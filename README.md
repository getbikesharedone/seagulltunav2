<a href="https://aimeos.org/">
    <img src="https://avatars1.githubusercontent.com/u/31987199?v=4&s=200" alt="Getterdone logo" title="Getterdone" align="right" height="60" />
</a>

Bike Network Location Finder
======================

The repository contains the Bike Network Location Finder website. It's the v2 rework of https://github.com/getbikesharedone/seagulltuna. The front end utilizes Vue, Vuetify, and Vuex. A Go server for handling API calls and an SQLite server are provided. Over 500 network points and 400k+ station points are included in the database. Webpack is setup for minification, hot-reload, and includes webpack-dev-server for running the website locally.

[![Bike Network Location Finder demo](https://i.imgur.com/zvSk9PM.png)](https://seagulltunav2.neuralspaz.com/#//)

## Table of content

- [Installation](#installation)
    - [I already have my own server](#i-already-have-my-own-server)
    - [webpack-dev-server](#webpack-dev-server)
- [Other Useful Commands](#other-useful-commands)
    - [Extension](#extension)
    - [Database](#database)
- [Usage](#usage)
    - [Hamburger](#hamburger)
    - [Markers](#markers)
    - [Station Settings](#station-settings)
    - [Station Reviews](#station-reviews)

## Installation

### I already have my own server

Download the latest release at https://github.com/getbikesharedone/seagulltunav2/releases.

Place the index.html and unzipped static.7z in your server's public directory.


### webpack-dev-server (for those who don't have their own)

Clone the repo:
```bash
git clone https://github.com/getbikesharedone/seagulltuna.git
```

Install dependencies:
```bash
npm i
```

Port 9090 must be available. Start Go server:
```bash
./seagulltuna.exe -dev
```

Port 8080 must be available. Start webpack-dev-server server:
```bash
npm run dev
```

### Other Useful Commands

```bash
# build for production with minification
npm run build
```

```bash
# build for production and view the bundle analyzer report
npm run build --report
```

## Usage

The website isn't perfect as of the current release, but if you follow this guide, then all features should work.

### Hamburger

![hamburger](https://i.imgur.com/YQsbnTk.png)

Opens and closes the drawer.

### Markers

As of the current release, you can't tell the difference between network and station markers. The best way to tell is when the drawer loads station or network info for the marker that you clicked.

![network loaded](https://i.imgur.com/y5HFQhP.png)

Once you have a network loaded, all the markers become stations. Click on a station marker loads the station's info, including a street view image that replaces the black box.

![station loaded](https://i.imgur.com/xkaErjv.png)

### Station Settings

Each station comes with several options that indicate it's status. After you have a station loaded, you can click the pencil icon to permantly change an option.

![edit station settings](https://i.imgur.com/denxjVl.png)

### Station Reviews

Each station comes with its own set of reviews. You may access them by click the vertical dots next to the pencil. The drawer will change.

![reviews drawer](https://i.imgur.com/Cd3cVkp.png)

You can add a review by clicking the plus icon at the top-left.

![add review](https://i.imgur.com/VDVi1Yw.png)

You can edit the reviews by clicking the Edit button located inside each review's card.

![edit reviews](https://i.imgur.com/Ik0thJH.png)