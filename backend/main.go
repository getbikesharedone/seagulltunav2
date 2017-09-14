// +build go1.8

package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris"
	irisctx "github.com/kataras/iris/context"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func main() {
	log.Println("starting seagull")
	var err error

	// if err := buildDatabase(); err != nil {
	// 	log.Fatalf("build database error: ", err)
	// }
	// if err := buildNetworkExtents(); err != nil {
	// 	log.Fatalf("build database extents: ", err)
	// }

	if db, err = sqlx.Open("sqlite3", "bsn.db"); err != nil {
		log.Fatalf("database error: ", err)
	}

	var NetworkList []Network
	db.Select(&NetworkList, "SELECT Company, ID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks")

	time.Sleep(time.Millisecond * 50)
	log.Println(len(NetworkList))

	srv := iris.New()
	srv.Use()
	srv.StaticWeb("/", "../frontend/www")
	srv.Get("/api/network/{id:string}", getDetail)
	srv.Get("/api/network", getNetworkList)
	srv.Get("/api/station/{id:string}", getStation)
	// srv.Post("/api/station/{id:string}", updateStation)
	// srv.Post("/api/station/{id:string}/tag", tagStation)
	// srv.Post("/api/station/{id:string}/review", reviewStation)

	if err := srv.Run(iris.Addr(":9090"), iris.WithoutVersionChecker); err != nil {
		log.Fatalf("failed to start http server: %v\n", err)
	}
}

func getStation(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getStation")
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.NotFound()
	}
	stations := []Station{}
	db.Select(&stations, "SELECT ID, Name, Latitude, Longitude FROM stations where ID=$1", id)
	if len(stations) == 0 {
		ctx.NotFound()
	}
	ctx.Gzip(true)
	ctx.JSON(stations)

}

// func updateStation(ctx irisctx.Context) {
// 	id := ctx.Params().Get("id")
// 	if id == "" {
// 		ctx.Err()
// 	}
// }

// func tagStation(ctx irisctx.Context) {
// 	id := ctx.Params().Get("id")
// 	if id == "" {
// 		ctx.Err()
// 	}
// }

// func reviewStation(ctx irisctx.Context) {
// 	id := ctx.Params().Get("id")
// 	if id == "" {
// 		ctx.Err()
// 	}
// }

func getDetail(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getDetail")
	id := ctx.Params().Get("id")
	if id == "" {
		log.Println("\n\n id is nil")
		ctx.NotFound()
		return
	}
	var net Network
	db.Get(&net, "SELECT UID, Company, ID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks WHERE ID=$1", id)
	if net.ID == "" {
		log.Printf("network %v does not exist", id)
		ctx.NotFound()
		return
	}
	var stations = []Station{}
	db.Select(&stations, "SELECT ID, Name, Latitude, Longitude FROM stations where NetworkUID=$1", net.UID)
	if len(stations) == 0 {
		ctx.NotFound()
		return
	}
	net.Stations = append(net.Stations, stations...)

	ctx.Gzip(true)
	ctx.JSON(net)
}

func getNetworkList(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getNetworkList")

	var NetworkList []Network
	db.Select(&NetworkList, "SELECT Company, ID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks")
	if len(NetworkList) == 0 {
		ctx.NotFound()
		return
	}
	time.Sleep(time.Millisecond * 50)
	log.Println(len(NetworkList))

	ctx.Gzip(true)
	ctx.JSON(NetworkList)
	return

}

func timeLog(start time.Time, name string) {
	taken := time.Since(start)
	log.Printf("%s took %s", name, taken)
}
