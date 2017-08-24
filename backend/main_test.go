package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestMain(m *testing.M) {
	fmt.Println("start init")
	var err error
	bsn, err = getSeedData()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range bsn.Networks {
		netmap[v.ID] = v
		for _, vv := range v.Stations {
			stationmap[vv.ID] = vv
		}
	}
	exitStatus := m.Run()
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
	// check that length is the same as origional
	if len(networks) != len(bsn.Networks) {
		t.Errorf("expected array length of response to be %d but got %d", len(bsn.Networks), len(networks))
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

	for _, v := range bsn.Networks {
		// seems to be a nil in the data set
		if v.ID != "" {
			tests = append(tests, Test{v.ID, httptest.StatusOK, "application/json"})
		}
	}

	for k, test := range tests {
		// check response status

		rawResponse := e.GET("/api/network/" + test.id).Expect().Status(test.status).ContentType(test.content).Body().Raw()
		// check if response is valid json
		if test.status == httptest.StatusOK {
			var network Network
			if err := json.Unmarshal([]byte(rawResponse), &network); err != nil {
				fmt.Println(k, test)
				t.Errorf("Failed on %d: %s expected to be decode into Network: %v\n ", k, test.id, err)
			}

			if network.ID != test.id {
				t.Errorf("expected array length of response to be %s but got %s", network.ID, test.id)
			}
		}

	}

}
