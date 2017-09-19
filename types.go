package main

import "time"

type Review struct {
	ReviewID  int       `db:"ReviewID" json:"id"`
	StationID int       `db:"StationID" json:"stationuid"`
	UserID    int       `db:"UserID" json:"-"`
	User      string    `db:"User" json:"user"`
	TimeStamp time.Time `db:"TimeStamp" json:"time"`
	Body      string    `db:"Body" json:"body"`
	Rating    int       `db:"Rating" json:"rating"`
}

type Station struct {
	StationID  int       `db:"StationID" json:"id"`
	NetworkID  int       `db:"NetworkID" json:"-"`
	Name       string    `db:"Name" json:"name"`
	EmptySlots int       `db:"EmptySlots" json:"empty"`
	FreeBikes  int       `db:"FreeBikes" json:"free"`
	Safe       bool      `db:"Safe" json:"safe"`
	Open       bool      `db:"Open" json:"open"`
	TimeStamp  time.Time `db:"TimeStamp" json:"time"`
	Lat        float64   `db:"Latitude" json:"lat"`
	Lng        float64   `db:"Longitude" json:"lng"`
	Reviews    []Review  `db:"-" json:"reviews,omitempty"`
}

type Network struct {
	NetworkID int       `db:"NetworkID" json:"id,omitempty"`
	Company   string    `db:"Company" json:"company,omitempty"`
	Name      string    `db:"Name" json:"name,omitempty"`
	City      string    `db:"City" json:"city,omitempty"`
	Country   string    `db:"Country" json:"country,omitempty"`
	Lat       float64   `db:"Latitude" json:"lat"`
	Lng       float64   `db:"Longitude" json:"lng"`
	HSpan     int       `db:"HSpan" json:"hspan,omitempty"`
	VSpan     int       `db:"VSpan" json:"vspan,omitempty"`
	CenterLat float64   `db:"CenterLat" json:"clat,omitempty"`
	CenterLng float64   `db:"CenterLng" json:"clng,omitempty"`
	Stations  []Station `db:"-" json:"stations,omitempty"`
}

var schema = `
DROP TABLE IF EXISTS networks;
CREATE TABLE networks (  
	NetworkID   INTEGER PRIMARY KEY,
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
	HSpan		INTEGER DEFAULT 0,
	VSpan		INTEGER DEFAULT 0,
	CenterLat	DOUBLE DEFAULT 0,
	CenterLng	DOUBLE DEFAULT 0
);
DROP TABLE IF EXISTS stations;
CREATE TABLE stations (
	StationID   INTEGER PRIMARY KEY,
	NetworkID   INTEGER,
	Name        VARCHAR(250),
	EmptySlots  INTEGER,
	FreeBikes   INTEGER,
	Extra       INTEGER,
	Safe		INTEGER DEFAULT 1,
	Open		INTEGER DEFAULT 1,
	Latitude    DOUBLE,
	Longitude   DOUBLE,
	TimeStamp   DATETIME
);
DROP TABLE IF EXISTS reviews;
CREATE TABLE reviews (
	ReviewID    INTEGER PRIMARY KEY,
	StationID	INTEGER,
	UserID		INTEGER,
	User		VARCHAR(80)  DEFAULT '',
	Body        VARCHAR(250) DEFAULT '',
	Rating      INTEGER,
	TimeStamp   DATETIME
);
`
