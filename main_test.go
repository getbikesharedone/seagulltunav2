package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/httptest"
	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	fmt.Println("start init")

	if _, err := os.Stat("bsn.db"); os.IsNotExist(err) {
		log.Panicf("cannot run tests without database, please run with -rebuild")
	}

	srcDB, err := os.Open("bsn.db")
	if err != nil {
		log.Panic(err)
	}
	defer srcDB.Close()

	testDB, err := os.Create("test.db")
	if err != nil {
		log.Panic(err)
	}

	if _, err = io.Copy(testDB, srcDB); err != nil {
		log.Panic(err)
	}
	err = testDB.Sync()
	if err != nil {
		log.Panic(err)

	}

	if db, err = sqlx.Open("sqlite3", "test.db"); err != nil {
		log.Fatalf("database error: ", err)
	}

	exitStatus := m.Run()
	db.Close()
	os.Remove("test.db")
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
	tests := []Test{
		{id: 20000, status: httptest.StatusNotFound, content: ""},
		{id: 0, status: httptest.StatusNotFound, content: ""},
	}

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
	tests := []Test{
		{id: 50000, status: httptest.StatusNotFound, content: ""},
		{id: 0, status: httptest.StatusNotFound, content: ""},
	}

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

func TestGetStationWithBadInput(t *testing.T) {
	testsrv := newSrv()

	e := httptest.New(t, testsrv)
	type Test struct {
		request string
		status  int
		content string
	}
	tests := []Test{
		{request: "50000", status: httptest.StatusNotFound, content: ""},
		{request: "0", status: httptest.StatusNotFound, content: ""},
		{request: "", status: httptest.StatusNotFound, content: ""},
		{request: "notint", status: httptest.StatusNotFound, content: ""},
	}
	for k, tc := range tests {

		response := e.GET("/api/station/" + tc.request).Expect().Status(tc.status).ContentType(tc.content).Body().Raw()
		if tc.status == httptest.StatusOK {
			var station Station
			if err := json.Unmarshal([]byte(response), &station); err != nil {
				t.Errorf("Failed on %d: %s expected to be decode into station with error: %v\n ", k, tc.request, err)
			}
		}
	}

}

// func Test_updateStationHandler(t *testing.T) {
// 	type args struct {
// 		ctx irisctx.Context
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			updateStationHandler(tt.args.ctx)
// 		})
// 	}
// }

func TestUpdateStation(t *testing.T) {

	tests := []struct {
		name  string
		input Station
		want  Station
	}{
		{
			name:  "no change",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 11, Safe: false, Open: false},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 11, Safe: false, Open: false},
		},
		{
			name:  "change Empty Slots",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 2, FreeBikes: 11, Safe: false, Open: false},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 2, FreeBikes: 11, Safe: false, Open: false},
		},
		{
			name:  "change Free Bikes",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: false, Open: false},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: false, Open: false},
		},
		{
			name:  "change Safe",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: true, Open: false},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: true, Open: false},
		},
		{
			name:  "change Open",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: true, Open: true},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 1, FreeBikes: 12, Safe: true, Open: true},
		},
		{
			name:  "change all",
			input: Station{StationID: 1, NetworkID: 1, EmptySlots: 0, FreeBikes: 0, Safe: false, Open: false},
			want:  Station{StationID: 1, NetworkID: 1, EmptySlots: 0, FreeBikes: 0, Safe: false, Open: false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := updateStation(tc.input)
			if err != nil {
				t.Errorf("got and error that was unexpected :%v", err)
			}
			if got.StationID != tc.want.StationID {
				t.Errorf("updateStation(StationID) = %v, want %v", got.StationID, tc.want.StationID)
			}
			if got.EmptySlots != tc.want.EmptySlots {
				t.Errorf("updateStation(EmptySlots) = %v, want %v", got.EmptySlots, tc.want.EmptySlots)
			}
			if got.FreeBikes != tc.want.FreeBikes {
				t.Errorf("updateStation(FreeBikes) = %v, want %v", got.FreeBikes, tc.want.FreeBikes)
			}
			if got.Safe != tc.want.Safe {
				t.Errorf("updateStation(Safe) = %v, want %v", got.Safe, tc.want.Safe)
			}
			if got.Open != tc.want.Open {
				t.Errorf("updateStation(Open) = %v, want %v", got.Open, tc.want.Open)
			}
		})
	}
}
