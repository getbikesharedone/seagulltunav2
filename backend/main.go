package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var bsn bikeShareNetwork

func main() {
	log.Println("starting seagull")
	var err error
	bsn, err = getSeedData()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("../frontend/www")))
	http.HandleFunc("/api/networks", GzipFunc(listNetworksHandler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func listNetworksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("listNetworks", r.RemoteAddr, r.RequestURI)

	type Shortnet struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Location `json:"location"`
	}

	type Response struct {
		Networks []Shortnet `json:"networks"`
	}

	var networks Response
	for _, v := range bsn.Networks {
		short := Shortnet{
			ID:       v.ID,
			Name:     v.Name,
			Location: v.Location,
		}
		networks.Networks = append(networks.Networks, short)

	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(networks); err != nil {
		log.Println(err)
	}

}

func getSeedData() (bikeShareNetwork, error) {

	// check if we have it first
	if _, err := os.Stat("bsn.json"); os.IsNotExist(err) {
		const bikeShareAPI = "http://api.citybik.es/v2/networks"

		var bsn bikeShareNetwork
		resp, err := http.Get(bikeShareAPI)
		if err != nil {
			return bsn, err
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&bsn); err != nil {
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
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&detail); err != nil {
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

		enc := json.NewEncoder(out)
		err = enc.Encode(bsn)
		if err != nil {
			log.Println(err)
		}

	}

	var bsn bikeShareNetwork

	in, err := os.Open("bsn.json")
	if err != nil {
		return bsn, err
	}

	dec := json.NewDecoder(in)
	if err := dec.Decode(&bsn); err != nil {
		return bsn, err
	}
	// get data from file

	return bsn, nil

}

type bikeShareNetwork struct {
	Networks []Network `json:"networks"`
}

type Network struct {
	Company  []string `json:"company"`
	Href     string   `json:"href"`
	ID       string   `json:"id"`
	Location `json:"location"`
	Name     string    `json:"name"`
	Stations []Station `json:"stations,omitempty"`
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
}

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (n *Network) UnmarshalJSON(data []byte) error {
	// Need too handle the one case company strings vs []string
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
						log.Println("failed conversion to string")
					}
				}
			}
		}
	}
	return nil
}

type stations []Station

func (ss stations) within(lat, lng, dist float64) stations {

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].distance(lat, lng) < ss[j].distance(lat, lng)
	})

	var abriged []Station
	for _, v := range ss {
		if v.distance(lat, lng) < dist {
			abriged = append(abriged, v)
		} else {
			return abriged
		}
	}
	return abriged
}

func (s *Station) distance(lat, lng float64) float64 {
	// http://www.movable-type.co.uk/scripts/latlong.html
	R := 6371e3 // radius

	φ1 := (math.Pi * lat) / 180
	φ2 := (math.Pi * lng) / 180

	Δφ := (s.Latitude - lat) * math.Pi / 180
	Δλ := (s.Longitude - lng) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func GzipFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Content-Type", "application/json")
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
