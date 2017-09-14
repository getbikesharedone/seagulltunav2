package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type APINetList struct {
	Networks []struct {
		Company  interface{} `json:"company"`
		Href     string      `json:"href"`
		ID       string      `json:"id"`
		Location struct {
			City      string  `json:"city"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Name     string `json:"name"`
		GbfsHref string `json:"gbfs_href,omitempty"`
		License  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
	} `json:"networks"`
}

type APINetwork struct {
	Network struct {
		Company  interface{} `json:"company"`
		Href     string      `json:"href"`
		ID       string      `json:"id"`
		Location struct {
			City      string  `json:"city"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Name     string `json:"name"`
		GbfsHref string `json:"gbfs_href,omitempty"`
		License  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
		Stations []struct {
			EmptySlots int `json:"empty_slots"`
			// Extra      struct {
			// 	UID int `json:"uid"`
			// } `json:"extra"`
			FreeBikes int       `json:"free_bikes"`
			ID        string    `json:"id"`
			Latitude  float64   `json:"latitude"`
			Longitude float64   `json:"longitude"`
			Name      string    `json:"name"`
			TimeStamp time.Time `json:"timestamp"`
		} `json:"stations"`
	} `json:"network"`
}

var schema = `
DROP TABLE IF EXISTS networks;
CREATE TABLE networks (
	UID         INTEGER PRIMARY KEY,
	ID          VARCHAR(80),
	Company     VARCHAR(250),
	Href        VARCHAR(250),
	City        VARCHAR(250), 
	Country     VARCHAR(250), 
	Latitude    DOUBLE, 
	Longitude   DOUBLE,
	Name        VARCHAR(250),
	GbfsHref    VARCHAR(250),
	LicenseName VARCHAR(250),
	LicenseURL  VARCHAR(250),
	HSpan		INTEGER,
	VSpan		INTEGER,
	CenterLat	DOUBLE,
	CenterLng	DOUBLE
);
DROP TABLE IF EXISTS stations;
CREATE TABLE stations (
	UID         INTEGER PRIMARY KEY,
	ID          VARCHAR(80),
	NetworkUID  VARCHAR(80),
	Name        VARCHAR(250),
	EmptySlots  INTEGER,
	FreeBikes   INTEGER,
	Extra       INTEGER,
	Safe		INTEGER DEFAULT 1,
	Latitude    DOUBLE,
	Longitude   DOUBLE,
	TimeStamp   DATETIME
);
DROP TABLE IF EXISTS reviews;
CREATE TABLE reviews (
	UID         INTEGER PRIMARY KEY,
	StationUID	VARCHAR(80)  DEFAULT '',
	Body        VARCHAR(250) DEFAULT '',
	Rating      INTEGER,
	TimeStamp   DATETIME
);
`

// import pull all the data from http://api.citybik.es/v2/networks API
func downloadNetworks() error {

	const bikeShareAPI = "http://api.citybik.es/v2/networks"

	var bsnList APINetList
	resp, err := http.Get(bikeShareAPI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&bsnList); err != nil {
		return err
	}

	var networks []APINetwork

	for count, v := range bsnList.Networks {

		log.Printf("(%3d of %3d) detail of %s, \n", count+1, len(bsnList.Networks), v.ID)

		var network APINetwork
		resp, err := http.Get(bikeShareAPI + "/" + v.ID)
		if err != nil {
			log.Printf("error reading network detail for %s: %v", v.ID, err)
			break
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&network); err != nil {
			log.Printf("error decoding detail json for %s: %v", v.ID, err)
			break
		}
		networks = append(networks, network)

	}

	out, err := os.Create("raw.json")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := json.NewEncoder(out).Encode(networks); err != nil {
		log.Println(err)
	}

	return nil
}

func buildDatabase() error {

	if _, err := os.Stat("raw.json"); os.IsNotExist(err) {
		if err := downloadNetworks(); err != nil {
			return err
		}
	}

	data, err := os.Open("raw.json")
	if err != nil {
		log.Fatal(err)
	}

	var networks []APINetwork
	if err := json.NewDecoder(data).Decode(&networks); err != nil {
		return err
	}

	var db *sqlx.DB
	if db, err = sqlx.Open("sqlite3", "bsn.db"); err != nil {
		return nil
	}
	defer db.Close()
	db.MustExec(schema)

	tx := db.MustBegin()
	for _, net := range networks {
		network := net.Network
		var companySlice []string
		switch network.Company.(type) {
		case string:
			c, ok := network.Company.(string)
			if ok {
				companySlice = append(companySlice, c)
			}
		case []interface{}:
			cc, ok := network.Company.([]interface{})
			if ok {
				for _, vv := range cc {
					switch vv.(type) {
					case string:
						c, ok := vv.(string)
						if ok {
							companySlice = append(companySlice, c)
						} else {
							log.Printf("failed conversion to string for: %v\n", network.ID)
						}
					}
				}
			}
		}
		var company string
		for _, line := range companySlice {
			company += line
		}
		if network.ID == "" || network.Company == "" {
			continue
		}
		tx.MustExec("INSERT INTO networks (ID,Company,Href,City,Country,Latitude,Longitude,Name,GbfsHref,LicenseName,LicenseURL) VALUES ($1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
			network.ID,
			company,
			network.Href,
			network.Location.City,
			network.Location.Country,
			network.Location.Latitude,
			network.Location.Longitude,
			network.Name,
			network.GbfsHref,
			network.License.Name,
			network.License.URL)

		var NetworkUID string
		if err = tx.Get(&NetworkUID, "SELECT uid FROM networks WHERE ID=$1", network.ID); err != nil {
			log.Println(err)
			continue
		}

		for _, station := range network.Stations {
			tx.MustExec("INSERT INTO stations (ID,NetworkUID,Name,EmptySlots,FreeBikes,Latitude,Longitude,TimeStamp) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
				station.ID,
				NetworkUID,
				station.Name,
				station.EmptySlots,
				station.FreeBikes,
				station.Latitude,
				station.Longitude,
				station.TimeStamp)

		}
	}
	tx.Commit()

	return nil
}

func buildNetworkExtents() error {

	db, err := sqlx.Open("sqlite3", "bsn.db")
	if err != nil {
		return nil
	}
	defer db.Close()

	var Networks []Network
	db.Select(&Networks, "SELECT UID FROM networks")
	if len(Networks) == 0 {
		return errors.New("no networks in database")
	}
	for _, network := range Networks {
		var stations []Station
		db.Select(&stations, "SELECT UID, Latitude, Longitude FROM stations WHERE NetworkUID=$1", network.UID)
		if len(stations) == 0 {
			fmt.Println(network.UID)
			// return errors.New("no stations in network")
			continue
		}
		clat, clng, hspan, vspan := extents(stations)

		db.MustExec("UPDATE networks SET CenterLat=$1, CenterLng=$2, HSpan=$3, VSpan=$4 WHERE UID=$5",
			clat, clng, hspan, vspan, network.UID)

	}
	return nil
}

func distance(lat1, lng1, lat2, lng2 float64) float64 {
	R := 6371e3 // radius

	φ1 := (math.Pi * lat1) / 180
	φ2 := (math.Pi * lat2) / 180

	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	meters := R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return meters
}

func extents(ss []Station) (centerlat float64, centerlng float64, hspan int, vspan int) {
	// Meters per pixel = 156543.03392 * Math.cos(latLng.lat() * Math.PI / 180) / Math.pow(2, zoom)
	if len(ss) < 2 {
		return ss[0].Lat, ss[0].Lng, 500.0, 500.0
	}
	latMin := 90.0
	latMax := -90.0
	lngMin := 180.0
	lngMax := -180.0

	for _, v := range ss {
		if v.Lat < latMin {
			latMin = v.Lat
		}
		if v.Lat > latMax {
			latMax = v.Lat
		}
		if v.Lng < lngMin {
			lngMin = v.Lng
		}
		if v.Lng > lngMax {
			lngMax = v.Lng
		}
	}
	centerlat = ((latMax - latMin) / 2) + latMin
	centerlng = ((lngMax - lngMin) / 2) + lngMin

	vspan = int(distance(latMin, centerlng, latMax, centerlng))
	hspan = int(distance(centerlat, lngMin, centerlat, lngMax))
	return
}
