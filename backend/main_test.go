package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/httptest"
)

func TestMain(m *testing.M) {
	fmt.Println("start init")
	var err error
	if db, err = sqlx.Open("sqlite3", "bsn.db"); err != nil {
		log.Fatalf("database error: ", err)
	}

	exitStatus := m.Run()
	db.Close()
	os.Exit(exitStatus)
}

func TestIndex(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)

	// check page exists
	e.GET("/").Expect().Status(httptest.StatusOK).ContentType("text/html")

}

func TestGetNetworkList(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)

	// check response status
	rawResponse := e.GET("/api/network/").Expect().Status(httptest.StatusOK).ContentType("application/json").Body().Raw()
	// check if response is valid json
	var networks []Network
	if err := json.Unmarshal([]byte(rawResponse), &networks); err != nil {
		t.Errorf("expected to be decode into bsn network: %v", err)
	}
}

func TestGetNetworkDetail(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)
	type Test struct {
		id      string
		status  int
		content string
	}
	tests := []Test{{id: "garbage", status: httptest.StatusNotFound, content: ""}}

	var Networks []Network
	db.Select(&Networks, "SELECT ID FROM networks")

	for _, v := range Networks {
		tests = append(tests, Test{v.ID, httptest.StatusOK, "application/json"})
	}

	for k, test := range tests {
		// check response status

		rawResponse := e.GET("/api/network/" + test.id).Expect().Status(test.status).ContentType(test.content).Body().Raw()
		// check if response is valid json
		if test.status == httptest.StatusOK {
			var network Network
			if err := json.Unmarshal([]byte(rawResponse), &network); err != nil {
				fmt.Println(k, test)
				t.Errorf("Failed on %d: %s expected to be decode into Network with error: %v\n ", k, test.id, err)
			}

			if network.ID != test.id {
				t.Errorf("expected response to be %s but got %s", network.ID, test.id)
			}
		}

	}

}

func TestGetNetworkDetailConcurrent(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)
	type Test struct {
		id      string
		status  int
		content string
	}
	tests := []Test{{id: "garbage", status: httptest.StatusNotFound, content: ""}}

	var Networks []Network
	db.Select(&Networks, "SELECT ID FROM networks")

	for _, v := range Networks {
		tests = append(tests, Test{v.ID, httptest.StatusOK, "application/json"})
	}
	var wg sync.WaitGroup
	for k, testcase := range tests {
		// throw all the requests at once at server
		wg.Add(1)
		go func(tc Test, tn int) {

			rawResponse := e.GET("/api/network/" + tc.id).Expect().Status(tc.status).ContentType(tc.content).Body().Raw()
			if tc.status == httptest.StatusOK {
				var network Network
				if err := json.Unmarshal([]byte(rawResponse), &network); err != nil {
					fmt.Println(k, tc)
					t.Errorf("Failed on %d: %s expected to be decode into Network with error: %v\n ", tn, tc.id, err)
				}

				if network.ID != tc.id {
					t.Errorf("expected response to be %s but got %s", network.ID, tc.id)
				}
			}
			wg.Done()
		}(testcase, k)

	}
	wg.Wait()
}
