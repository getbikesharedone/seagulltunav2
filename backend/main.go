// +build go1.8

package main

import (
	"flag"
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

	var reBuildDB = flag.Bool("rebuild", false, "rebuild database")
	flag.Parse()

	var err error

	if *reBuildDB {
		if err := buildDatabase(); err != nil {
			log.Fatalf("build database error: ", err)
		}
		if err := buildNetworkExtents(); err != nil {
			log.Fatalf("build database extents: ", err)
		}
	}

	if db, err = sqlx.Open("sqlite3", "bsn.db"); err != nil {
		log.Fatalf("database error: ", err)
	}

	var NetworkList []Network
	db.Select(&NetworkList, "SELECT Company, ID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks")

	time.Sleep(time.Millisecond * 50)
	log.Println(len(NetworkList))

	srv := newSrv()

	if err := srv.Run(iris.Addr(":9090"), iris.WithoutVersionChecker); err != nil {
		log.Fatalf("failed to start http server: %v\n", err)
	}
}

func newSrv() *iris.Application {
	srv := iris.New()
	srv.Use()
	srv.StaticWeb("/", "www/")
	srv.Get("/api/network/{id:string}", getDetail)
	srv.Get("/api/network", getNetworkList)
	srv.Get("/api/station/{id:string}", getStation)

	// srv.Post("/api/station/{id:string}", updateStation)
	srv.Post("/api/station/{id:string}/tag", tagStation)
	srv.Post("/api/station/{id:string}/review", reviewStation)
	return srv
}

func getStation(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getStation")
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.NotFound()
		return
	}
	var stations = []Station{}
	db.Select(&stations, "SELECT UID, ID, Name, Latitude, Longitude, EmptySlots, FreeBikes, Safe FROM stations where ID=$1", id)
	if len(stations) == 0 {
		log.Printf("no stations in network %v", id)
		ctx.NotFound()
		return
	}

	log.Println(stations[0])
	var reviews = []Review{}
	db.Select(&reviews, "SELECT Body, Rating, TimeStamp FROM reviews where StationUID=$1", stations[0].UID)
	log.Println(reviews)
	if len(reviews) != 0 {
		log.Printf("no reviews for station %v", id)
		stations[0].Reviews = append(stations[0].Reviews, reviews...)
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

func tagStation(ctx irisctx.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
	}
}

func reviewStation(ctx irisctx.Context) {
	defer timeLog(time.Now(), "reviewStation")
	log.Printf("%+#v\n", ctx)
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
		return
	}
	var review Review
	err := ctx.ReadForm(&review)
	if err != nil {
		log.Println(err)
		ctx.Err()
		return
	}
	review.TimeStamp = time.Now()
	var station Station
	db.Get(&station, "SELECT UID FROM stations WHERE ID=$1", id)
	review.StationUID = station.UID
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO reviews (StationUID,TimeStamp,Body,Rating) VALUES ($1, $2, $3, $4)",
		review.StationUID,
		review.TimeStamp,
		review.Body,
		review.Rating)
	tx.Commit()
	return
}

func getDetail(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getDetail")
	id := ctx.Params().Get("id")
	if id == "" {
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
		log.Printf("no stations in network %v", id)
	}
	net.Stations = append(net.Stations, stations...)

	ctx.Gzip(true)
	ctx.JSON(net)
	return
}

func getNetworkList(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getNetworkList")

	var NetworkList []Network
	db.Select(&NetworkList, "SELECT Company, ID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks")
	if len(NetworkList) == 0 {
		ctx.NotFound()
		return
	}
	log.Println(len(NetworkList))

	ctx.Gzip(true)
	ctx.JSON(NetworkList)
	return

}

func timeLog(start time.Time, name string) {
	taken := time.Since(start)
	log.Printf("%s took %s", name, taken)
}
