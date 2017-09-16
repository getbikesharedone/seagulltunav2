window.Event = new Vue();

// Stores cluster reference so clearMarkers() can be called
let markerCluster
let map
const appVue = new Vue({
  el: "#app",
  data: {
    networks: [],
    stations: []
  },
  created() {
    this.getNetworks()
  },
  
  methods: {
    initMap: function () {
      var myLatLng = { lat: 0, lng: 0 };

      map = new google.maps.Map(document.getElementById('map'), {
        zoom: 3,
        center: myLatLng
      });

      networkMarkers = this.addNetworkMarkers(map, this.networks);
      // stationMarkers = this.addStationsMarkers(map, this.stations);
      markerCluster = new MarkerClusterer(map, networkMarkers,
        { imagePath: '/m' });


    },
    addNetworkMarkers: function (map, networks) {
      let networkMarkers = []
      for (let i = 0; i < networks.length; i++) {
        const network = networks[i]
        let marker = new google.maps.Marker({
          position: {
            lat: network.lat,
            lng: network.lng,
          },
          map,
          title: network.name + "\n" + network.city,
          icon: {
            url: '/bike.png',
            size: new google.maps.Size(32, 32),
          },
          center: {
            lat: network.clat,
            lng: network.clng,
          },
          network,
        })
        const vm = this
        marker.addListener('click', function () {
          markerCluster.clearMarkers()
          vm.getStations(this.network)
          map.panTo(marker.center);
          // zoom
        });

        networkMarkers.push(marker)
      }
      return networkMarkers;
    },
    addStationsMarkers: function (map, stations) {
      let stationsMarkers = []
      for (let i = 0; i < stations.length; i++) {
        const station = stations[i]
        let marker = new google.maps.Marker({
          position: {
            lat: station.lat,
            lng: station.lng,
          },
          map,
          title: station.name,
          icon: {
            url: '/helmet.png',
            size: new google.maps.Size(32, 32),
          },
        })
        stationsMarkers.push(marker)
      }
      return stationsMarkers;
    },
    getStations: function (network) {
      axios
        .get("/api/network/" + network.id)
        .then(res => {
          if (res.status == 200) {
            if (res.data != null) {
              console.log(res.data)
              this.stations = res.data.stations;
              Event.$emit("stationsLoaded", this.stations);
            }
          }
        })
        .catch(error => {
          this.advice = "There was an error: " + error.message;
        });
    },
    getNetworks: function () {
      axios
        .get("/api/network")
        .then(res => {
          if (res.status == 200) {
            if (res.data != null) {
              this.networks = res.data;
              Event.$emit("networksLoaded", this.networks);
            }
          }
        })
        .catch(error => {
          this.advice = "There was an error: " + error.message;
        });
    }
  },
  mounted() {
    Event.$on("networksLoaded", networks => {
      this.initMap()
    });
    Event.$on("stationsLoaded", stations => {
      this.addStationsMarkers(map,stations);
    });
  }
}); 
