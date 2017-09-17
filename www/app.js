window.Event = new Vue();

// Stores references so clearMarkers() can be called
let markerCluster

// Reference to map so markers can be re-added on zoom_out
let map

Vue.component('modal', {
  template: '#modal-template'
})

Vue.component("stations-table", {
  template: `
  <table class="ui table unstackable inverted orange striped">
    <thead>
    <tr>
    <th>Name</th>
    <th>Free Bikes</th>
    <th>Open</th>
    <th>Safe</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="station in stations">
    <td>{{station.name}}</td>
    <td>{{station.free}}</td>
    <td>{{station.open}}</td>
    <td>{{station.safe}}</td>
    </tr>
    </tbody>
    </table>
    `,
  data() {
    return {
      stations: []
    };
  },
  created() {
    Event.$on("stationsLoaded", stations => {
      this.stations = stations;
    });
  }
});

const appVue = new Vue({
  el: "#app",
  data: {
    networks: [],
    stations: [],
    stationMarkers: [],
    activeNetwork: {},
    showModal: false
  },
  created() {
    this.getNetworks()
  },
  
  methods: {
    initMap: function () {
      let myLatLng = { lat: 0, lng: 0 };

      map = new google.maps.Map(document.getElementById('map'), {
        zoom: 3,
        center: myLatLng
      });

      networkMarkers = this.addNetworkMarkers(map, this.networks);
      markerCluster = new MarkerClusterer(map, networkMarkers,
        { imagePath: '/m' });

      const vm = this;
      map.addListener('zoom_changed', function () {
        zoomLevel = map.getZoom();
        if (markerCluster.getTotalMarkers() === 0 && zoomLevel < 10) {
          networkMarkers = vm.addNetworkMarkers(map, vm.networks)
          markerCluster = new MarkerClusterer(map, networkMarkers,
            { imagePath: '/m' })
          vm.deleteStationMarkers()
        }
      });
    },
    setMapOnAll: function (map) {
      for (var i = 0; i < this.stationMarkers.length; i++) {
        this.stationMarkers[i].setMap(map);
      }
    },
    deleteStationMarkers: function () {
      this.clearStationMarkers();
      this.stationMarkers = [];
    },
     clearStationMarkers:function() {
      this.setMapOnAll(null);
    },
    addNetworkMarkers: function (map, networks) {
      let networkMarkers = []
      const vm = this
      for (let i = 0; i < networks.length; i++) {
        const network = networks[i]

        let marker = new google.maps.Marker({
          position: {
            lat: network.lat,
            lng: network.lng,
          },
          map,
          title: network.name,
          icon: '/bike.png',
          network,
        })

        marker.addListener('click', function () {
          vm.getStations(this.network)
          Event.$on("stationsLoaded", function () {
            vm.addStationMarkers(map, vm.stations)
            var bounds = new google.maps.LatLngBounds();
            for (var i = 0; i < vm.stations.length; i++) {
              bounds.extend(new google.maps.LatLng(vm.stations[i].lat, vm.stations[i].lng));
            }
            map.fitBounds(bounds);
            markerCluster.clearMarkers()
          });
        });

        networkMarkers.push(marker)
      }
      return networkMarkers;
    },
    addStationMarkers: function (map, stations) {
      for (let i = 0; i < stations.length; i++) {
        const station = stations[i]
        let marker;
        if (station.lat !== undefined && station.lng !== undefined) {
          const lat = station.lat
          const lng = station.lng
          marker = new google.maps.Marker({
            position: {
              lat,
              lng,
            },
            map,
            title: station.name,
            icon: {
              url: '/helmet.png',
              size: new google.maps.Size(32, 32),
              origin: new google.maps.Point(0, 0),
              anchor: new google.maps.Point(16, 16),
            },
            station,
          })
          this.stationMarkers.push(marker)
        }
      }
      return this.stationMarkers
    },
    getStations: function (network) {
      axios
        .get("/api/network/" + network.id)
        .then(res => {
          if (res.status == 200) {
            if (res.data != null) {
              this.activeNetwork = res.data;
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
      console.log("Active network", this.activeNetwork)
    });
    Event.$on("clickStation", station => {
      // debug
      console.log(station)
      console.log("Active network", this.activeNetwork)
      console.log(this.showModal)
      // display modal
      this.showModal=true;
    });
  }
}); 