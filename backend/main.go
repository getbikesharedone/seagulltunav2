// +build go1.8

package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// in lou of database for now
var bsn bikeShareNetwork
var netmap = make(map[string]Network)

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
	}

	srv := iris.New()
	srv.StaticWeb("/", "../frontend/www")
	srv.Get("/api/network/{id:string}", getDetail)
	srv.Get("/api/network", getNetworkList)
	srv.Run(iris.Addr(":8080"))
}

func getDetail(ctx context.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.NotFound()
	}
	net, ok := netmap[id]
	if !ok {
		log.Printf("network %v does not exist", id)
		ctx.NotFound()
	}
	ctx.Gzip(true)
	ctx.JSON(net)
}

type Where struct {
	Lat, Lng float64
	Rng      float64
}

func getNetworkList(ctx context.Context) {
	at := Where{}
	err := ctx.ReadForm(&at)
	if err != nil {
		log.Println(err)
	}
	type Shortnet struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Location `json:"location"`
	}

	type Response struct {
		Networks []Shortnet `json:"networks"`
	}
	ctx.Gzip(true)

	if at.Lat == 0 || at.Lng == 0 {
		var short Response
		for _, v := range bsn.Networks {
			s := Shortnet{
				ID:       v.ID,
				Name:     v.Name,
				Location: v.Location,
			}
			short.Networks = append(short.Networks, s)
		}
		ctx.JSON(short)
	}
	if at.Rng == 0 {
		at.Rng = 160000 // 160 km, 100 Miles
	}
	// compute diatnce
	var localised Response
	for _, v := range bsn.Networks {
		distance := v.Location.howfar(at)
		if distance < at.Rng {
			s := Shortnet{
				ID:       v.ID,
				Name:     v.Name,
				Location: v.Location,
			}
			s.Location.Distance = int(distance)
			localised.Networks = append(localised.Networks, s)
		}
	}
	ctx.JSON(localised)
}

func (l Location) howfar(at Where) float64 {
	R := 6371e3 // radius

	φ1 := (math.Pi * at.Lat) / 180
	φ2 := (math.Pi * at.Lng) / 180

	Δφ := (l.Latitude - at.Lat) * math.Pi / 180
	Δλ := (l.Longitude - at.Lng) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
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
	Distance  int     `json:"distance,omitempty"`
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
						log.Println("failed conversion to string")
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
