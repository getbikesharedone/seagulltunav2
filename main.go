// +build go1.8

package main

import (
	"flag"
	"log"
	"strconv"
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
	var reBuildDB = flag.Bool("rebuild", false, "rebuild database")
	flag.Parse()

	// var err error

	if *reBuildDB {
		if err := buildDatabase(); err != nil {
			log.Fatalf("build database error: ", err)
		}
		if err := buildNetworkExtents(); err != nil {
			log.Fatalf("build database extents: %v", err)
		}
	}

	if db, err = sqlx.Open("sqlite3", "bsn.db"); err != nil {
		log.Fatalf("database error: ", err)
	}

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
	srv.Post("/api/station/{id:string}", updateStationHandler)
	srv.Post("/api/station/{id:string}/review", reviewStation)
	return srv
}

func getStation(ctx irisctx.Context) {
	defer timeLog(time.Now(), "getStation")
	idStr := ctx.Params().Get("id")

	if idStr == "" {
		ctx.NotFound()
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.NotFound()
		return
	}
	var stations Station
	db.Get(&stations, "SELECT StationID, Name, Latitude, Longitude, EmptySlots, FreeBikes, Safe, Open, TimeStamp FROM stations where StationID=$1", id)

	if stations.StationID == 0 {
		ctx.NotFound()
		return
	}
	var reviews = []Review{}
	db.Select(&reviews, "SELECT ReviewID, Body, Rating, TimeStamp FROM reviews where StationID=$1", stations.StationID)
	if len(reviews) != 0 {
		stations.Reviews = append(stations.Reviews, reviews...)
	}

	ctx.Gzip(true)
	ctx.JSON(stations)

}

func updateStationHandler(ctx irisctx.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.Err()
	}
	var s Station
	err := ctx.ReadForm(&s)
	if err != nil {
		log.Println(err)
		ctx.Err()
		return
	}
	u, err := updateStation(s)
	if err != nil {
		log.Printf("error updating staion: %v\n", err)
		ctx.Err()
		return
	}
	ctx.Gzip(true)
	ctx.JSON(u)

}

func updateStation(update Station) (Station, error) {
	var existing Station
	err := db.Get(&existing, "SELECT StationID, EmptySlots, FreeBikes, Safe FROM stations WHERE StationID=$1", update.StationID)
	if err != nil {
		return Station{}, err
	}
	tx := db.MustBegin()
	tx.MustExec("UPDATE stations SET FreeBikes=$1 WHERE StationID=$2", update.FreeBikes, existing.StationID)
	tx.MustExec("UPDATE stations SET EmptySlots=$1 WHERE StationID=$2", update.EmptySlots, existing.StationID)
	tx.MustExec("UPDATE stations SET Safe=$1 WHERE StationID=$2", update.Safe, existing.StationID)
	tx.MustExec("UPDATE stations SET Open=$1 WHERE StationID=$2", update.Open, existing.StationID)
	tx.MustExec("UPDATE stations SET TimeStamp=$1 WHERE StationID=$2", time.Now().UTC(), existing.StationID)
	if err := tx.Commit(); err != nil {
		return existing, err
	}

	var updated Station
	err = db.Get(&updated, "SELECT StationID, Name, Latitude, Longitude, EmptySlots, FreeBikes, Safe, Open, TimeStamp FROM stations where StationID=$1", existing.StationID)
	if err != nil {
		return existing, err
	}

	var reviews = []Review{}
	err = db.Select(&reviews, "SELECT Body, Rating, TimeStamp FROM reviews where StationID=$1", updated.StationID)
	if err != nil {
		return existing, err
	}
	if len(reviews) != 0 {
		updated.Reviews = append(updated.Reviews, reviews...)
	}

	return updated, nil

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
	review.TimeStamp = time.Now().UTC()
	var station Station
	db.Get(&station, "SELECT UID FROM stations WHERE StationID=$1", id)
	review.StationID = station.StationID
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO reviews (StationUID,TimeStamp,Body,Rating) VALUES ($1, $2, $3, $4)",
		review.StationID,
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
	db.Get(&net, "SELECT NetworkID, Company, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks WHERE NetworkID=$1", id)
	if net.NetworkID == 0 {
		log.Printf("network %v does not exist", id)
		ctx.NotFound()
		return
	}
	var stations = []Station{}
	db.Select(&stations, "SELECT StationID, Name, Latitude, Longitude, EmptySlots, FreeBikes, Safe, Open, TimeStamp FROM stations where NetworkID=$1", id)
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
	db.Select(&NetworkList, "SELECT Company, NetworkID, Name, City, Country, Latitude, Longitude, HSpan, VSpan, CenterLat, CenterLng FROM networks")
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
