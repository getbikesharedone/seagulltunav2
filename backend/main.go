// +build go1.8

package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// in lou of database for now
var bsn bikeShareNetwork

var netmapMu *sync.RWMutex
var netmap = make(map[string]Network)
var stationmap = make(map[string]Station)

func main() {
	log.Println("starting seagull")
	var err error
	bsn, err = getSeedData()
	if err != nil {
		log.Fatal(err)
	}

	// read only map from here
	// replace with database
	for _, v := range bsn.Networks {
		netmap[v.ID] = v
		for _, vv := range v.Stations {
			stationmap[vv.ID] = vv
		}
	}

	srv := newSrv()

	if err := srv.Run(iris.Addr(":8080"), iris.WithoutVersionChecker); err != nil {
		log.Fatalf("failed to start http server: %v\n", err)
	}
}

func newSrv() *iris.Application {
	s := iris.New()
	s.StaticWeb("/", "../frontend/www")
	s.Get("/api/network/{id:string}", getDetail)
	s.Get("/api/network", getNetworkList)
	s.Get("/api/station/{id:string}", getStation)
	s.Post("/api/station/{id:string}", updateStation)
	s.Post("/api/station/{id:string}/tag", tagStation)
	s.Post("/api/station/{id:string}/review", reviewStation)
	return s
}

func getStation(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.NotFound()
	}
	station, ok := stationmap[id]
	if !ok {
		ctx.NotFound()
	}
	ctx.Gzip(true)
	ctx.JSON(station)

}

func updateStation(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
	}
}

func tagStation(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
	}
}

func reviewStation(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
	}
}

func getDetail(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		log.Println("\n\n id is nil")
		ctx.NotFound()
		return
	}
	net, ok := netmap[id]
	if !ok {
		log.Printf("network %v does not exist", id)
		ctx.NotFound()
		return
	}
	ctx.Gzip(true)
	ctx.JSON(net)
}

func getNetworkList(ctx context.Context) {

	type Shortnet struct {
		Company   []string `json:"company,omitempty"`
		ID        string   `json:"id,omitempty"`
		Name      string   `json:"name,omitempty"`
		Location  `json:"location,omitempty"`
		MapWindow MapView `json:"mapwindow,omitempty"`
	}

	ctx.Gzip(true)

	var short []Shortnet
	for _, v := range bsn.Networks {
		s := Shortnet{
			Company:   v.Company,
			ID:        v.ID,
			Name:      v.Name,
			Location:  v.Location,
			MapWindow: v.MapWindow,
		}
		short = append(short, s)
	}
	ctx.JSON(short)
	return

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

func (m *MapView) extents(ss []Station) {
	// Meters per pixel = 156543.03392 * Math.cos(latLng.lat() * Math.PI / 180) / Math.pow(2, zoom)
	if len(ss) < 2 {
		return
	}
	latMin := 90.0
	latMax := -90.0
	lngMin := 180.0
	lngMax := -180.0
	for _, v := range ss {
		if v.Latitude < latMin {
			latMin = v.Latitude
		}
		if v.Latitude > latMax {
			latMax = v.Latitude
		}
		if v.Longitude < lngMin {
			lngMin = v.Longitude
		}
		if v.Longitude > lngMax {
			lngMax = v.Longitude
		}
	}
	latCenter := ((latMax - latMin) / 2) + latMin
	lngCenter := ((lngMax - lngMin) / 2) + lngMin

	m.Vspan = distance(latMin, lngCenter, latMax, lngCenter)
	m.Hspan = distance(latCenter, lngMin, latCenter, lngMax)

	m.Center.Lat = latCenter
	m.Center.Lng = lngCenter

}

func getSeedData() (bikeShareNetwork, error) {
	defer timeLog(time.Now(), "getSeedData")
	// check if we have it first
	if _, err := os.Stat("bsn.json"); os.IsNotExist(err) {
		const bikeShareAPI = "http://api.citybik.es/v2/networks"

		var bsn bikeShareNetwork
		resp, err := http.Get(bikeShareAPI)
		if err != nil {
			return bsn, err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&bsn); err != nil {
			return bsn, err
		}
		networks := len(bsn.Networks)
		for k, v := range bsn.Networks {
			log.Printf("(%3d of %3d) detail of %s, \n", k+1, networks, v.ID)
			var detail networkDetail
			resp, err := http.Get(bikeShareAPI + "/" + v.ID)
			if err != nil {
				log.Printf("error reading network detail for %s: %v", v.ID, err)
				break
			}
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
				log.Printf("error decoding detail json for %s: %v", v.ID, err)
				break
			}
			bsn.Networks[k] = detail.Detail
		}

		out, err := os.Create("bsn.json")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		if err := json.NewEncoder(out).Encode(bsn); err != nil {
			log.Println(err)
		}

	}

	var bsn bikeShareNetwork
	// get data from file
	in, err := os.Open("bsn.json")
	if err != nil {
		return bsn, err
	}

	if err := json.NewDecoder(in).Decode(&bsn); err != nil {
		return bsn, err
	}

	for k, v := range bsn.Networks {
		m := v.MapWindow
		m.extents(v.Stations)
		bsn.Networks[k].MapWindow = m
	}

	return bsn, nil

}

type bikeShareNetwork struct {
	Networks []Network `json:"networks"`
}

type Network struct {
	Company   []string `json:"company"`
	ID        string   `json:"id"`
	Location  `json:"location"`
	Name      string    `json:"name"`
	Stations  []Station `json:"stations,omitempty"`
	MapWindow MapView   `json:"mapwindow,omitempty"`
}

type MapView struct {
	Vspan  float64    `json:"vspan,omitempty"`
	Hspan  float64    `json:"hspan,omitempty"`
	Center Coordinate `json:"center,omitempty"`
}

type Coordinate struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

type networkDetail struct {
	Detail Network `json:"network"`
}

type Station struct {
	EmptySlots int       `json:"empty_slots"`
	FreeBikes  int       `json:"free_bikes"`
	ID         string    `json:"id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Name       string    ` json:"name"`
	Timestamp  time.Time `json:"timestamp"` // look up local offset at location
	Reviews    []Review  `json:"reviews,omitempty"`
	Closed     bool      `json:"closed,omitempty"`
}

type Review struct {
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Rating    int       `json:"rating"`
}

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance,omitempty"`
	CenterLat float64 `json:"centerlat,omitempty"`
	CenterLng float64 `json:"centerlng,omitempty"`
	Scale     int     `json:"scale,omitempty"`
}

func (n *Network) UnmarshalJSON(data []byte) error {
	// Need too handle the one case where company is string vs []string
	type ServerNetworks Network
	aux := &struct {
		Company interface{} `json:"company"`
		*ServerNetworks
	}{
		ServerNetworks: (*ServerNetworks)(n),
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	switch aux.Company.(type) {
	case string:
		c, ok := aux.Company.(string)
		if ok {
			n.Company = append(n.Company, c)
		}
	case []interface{}:
		cc, ok := aux.Company.([]interface{})
		if ok {
			for _, vv := range cc {
				switch vv.(type) {
				case string:
					c, ok := vv.(string)
					if ok {
						n.Company = append(n.Company, c)
					} else {
						log.Println("failed conversion to string for: ", aux.ServerNetworks.ID)
					}
				}
			}
		}
	}
	return nil
}

func timeLog(start time.Time, name string) {
	taken := time.Since(start)
	log.Printf("%s took %s", name, taken)
}
