package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
		id      int
		status  int
		content string
	}
	tests := []Test{{id: 20000, status: httptest.StatusNotFound, content: ""}}

	var Networks []Network
	db.Select(&Networks, "SELECT NetworkID FROM networks")

	for _, v := range Networks {
		tests = append(tests, Test{v.NetworkID, httptest.StatusOK, "application/json"})
	}

	for k, test := range tests {
		// check response status
		if k > 50 {
			break
		}
		rawResponse := e.GET("/api/network/" + strconv.Itoa(test.id)).Expect().Status(test.status).ContentType(test.content).Body().Raw()
		// check if response is valid json
		if test.status == httptest.StatusOK {
			var network Network
			if err := json.Unmarshal([]byte(rawResponse), &network); err != nil {
				t.Errorf("Failed on %d: %s expected to be decode into Network with error: %v\n ", k, test.id, err)
			}

			if network.NetworkID != test.id {
				t.Errorf("expected response to be %s but got %s", network.NetworkID, test.id)
			}
		}

	}

}

func TestGetNetworkDetailConcurrent(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)
	type Test struct {
		id      int
		status  int
		content string
	}
	tests := []Test{{id: 20000, status: httptest.StatusNotFound, content: ""}}

	var Networks []Network
	db.Select(&Networks, "SELECT ID FROM networks")

	for _, v := range Networks {
		tests = append(tests, Test{v.NetworkID, httptest.StatusOK, "application/json"})
	}
	var wg sync.WaitGroup
	for k, testcase := range tests {
		// throw all the requests at once at server
		if k > 50 {
			break
		}
		wg.Add(1)
		go func(tc Test, tn int) {

			rawResponse := e.GET("/api/network/" + strconv.Itoa(tc.id)).Expect().Status(tc.status).ContentType(tc.content).Body().Raw()
			if tc.status == httptest.StatusOK {
				var network Network
				if err := json.Unmarshal([]byte(rawResponse), &network); err != nil {
					t.Errorf("Failed on %d: %s expected to be decode into Network with error: %v\n ", tn, tc.id, err)
				}

				if network.NetworkID != tc.id {
					t.Errorf("expected response to be %s but got %s", network.NetworkID, tc.id)
				}
			}
			wg.Done()
		}(testcase, k)

	}
	wg.Wait()
}

func TestGetStation(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)
	type Test struct {
		id      int
		status  int
		content string
	}
	tests := []Test{{id: 50000, status: httptest.StatusNotFound, content: ""}}

	var Stations []Station
	db.Select(&Stations, "SELECT StationID FROM stations")

	for _, v := range Stations {
		tests = append(tests, Test{v.StationID, httptest.StatusOK, "application/json"})
	}
	var wg sync.WaitGroup
	for k, tc := range tests {
		if k > 1000 {
			break
		}
		rawResponse := e.GET("/api/station/" + strconv.Itoa(tc.id)).Expect().Status(tc.status).ContentType(tc.content).Body().Raw()
		if tc.status == httptest.StatusOK {
			var station Station
			if err := json.Unmarshal([]byte(rawResponse), &station); err != nil {
				t.Errorf("Failed on %d: %s expected to be decode into station with error: %v\n ", k, tc.id, err)
			}

			if station.StationID != tc.id {
				t.Errorf("expected response to be %d but got %d", station.StationID, tc.id)
			}
		}
	}
	wg.Wait()
}
