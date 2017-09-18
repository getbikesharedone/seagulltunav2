window.Event = new Vue();

// Stores references so clearMarkers() can be called
let markerCluster

// Reference to map so markers can be re-added on zoom_out
let map

Vue.component('modal', {
  template: '#modal-template'
})

Vue.component('rating', {
  template: '<div class="ui rating" data-rating="3" data-max-rating="5"></div>',
  mounted() {
    $('.ui.rating')
      .rating()
      ;

  }
})

Vue.component("stations-table", {
  template: `
  <div style="width: 100%; height: 100%;overflow-x: scroll;">

  <div class="ui styled fluid accordion">



  <table class="ui table unstackable inverted orange ">
  <thead>
    <tr><th><div class="ui horizontal yellow statistic">
    
    <div class="value">
      {{stations.length}}
    </div>
    <div class="label">
      Stations
    </div>
    <div class="value">
    {{networksLength}}
  </div>

  <div class="label">
    Networks
  </div>
  </div></th>
  </tr></thead>
  <tbody>
    <tr style="width:95%;margin-left:auto; margin-right:auto; margin-top:15px; margin-bottom:15px" class="ui card" v-for="(station, index) in stations">
    <div v-on:click="loadReviews(station.id)" class="title">
    <td>
    
    <span class="header">{{station.name}}</span>
    <div>
    <span>
      Open
      <i v-bind:class="{ 'check square icon': station.open,  'square outline icon': !station.open }"></i>
    </span>
    <span>
    Safe
    <i v-bind:class="{ 'check square icon': station.safe, 'square outline icon': !station.safe }"></i>
    </span>
    <span>
    Available: <span>{{station.free}}</span>
    </span>
    
  </div></span>
    
    
      </td>
      </div>  
  <div class=" content">
  <button v-on:click="callMultiple(station, index)" class="ui button orange">
  Settings
</button>
  <p v-if="availableReviews" style="color:black">No reviews available</p>
    <p v-else style="color:black">{{reviews[0]}}</p>
  </div>
  <i class="settings icon"></i>
  <div :class="'idx' + index + ' ui modal'">
  <i class="close icon"></i>
  <div class="header">
    Settings
  </div>
  <div style="margin:15px" class="ui form">
  <div class="field">
  <div class="ui checkbox slider">
  <input type="checkbox" name="isOpen" >
  <label>Open</label>
</div></div
<div class="field">
<div class="ui checkbox slider">
<input type="checkbox" name="isSafe">
<label>Safe</label>
</div></div>
<div style="margin:15px" class="ui right action input">

<input type="text" :value="station.free">
<button class="ui orange labeled icon button">
<i class="bicycle icon"></i>
Update Available
</button>
</div>
  <div class="actions">
    <div class="approve ui button">Close</div>
  </div>
  </div>
    </tr>
    
  </tbody>
</table>


</div>
    </div>
    `,
  data() {
    return {
      stations: [],
      networksLength: 0,
      reviews: [],
      isSafe: true,
      isOpen: true,
      currentStation: {}

    };
  },
  created() {
    Event.$on("stationsLoaded", stations => {
      this.stations = stations;
    });
    Event.$on("networksLoaded", networks => {
      this.networksLength = networks.length;
    });
  },
  mounted() {
    $('.ui.accordion')
      .accordion()
      ;
    $('.ui.modal')
      .modal()
      ;
  },
  methods: {
    callMultiple: function(station, index){
      this.showModal(index)
      this.addCheckboxListener(station)
      this.currentStation = station;
    },
    addCheckboxListener: function(station){
      vm = this;
      $('.ui.checkbox').checkbox({
        onChange: function () {
          $('.ui.checkbox').hasClass('checked') ? this.isSafe = true : this.isSafe = false;
          axios
          .post("/api/station/" + station.id, {
            station
          })
          .then(res => {
            console.log(res)
            if (res.status == 200) {
              if (res.data != null) {
                vm.currentStation.safe = res.data.safe
              }
            }
          })
          .catch(error => {
            this.advice = "There was an error: " + error.message;
          });
         
      }});
    },
    loadReviews: function (stationId) {
      axios
        .get("/api/station/" + stationId)
        .then(res => {
          if (res.status == 200) {
            if (res.data != null) {
              this.reviews = res.data.reviews;
              Event.$emit("reviewsLoaded", this.reviews);
            }
          }
        })
        .catch(error => {
          this.advice = "There was an error: " + error.message;
        });
    },
    availableReviews: function () {
      reviews[0] === undefined ? false : true
    },
    showModal: function (index) {
      $('.ui.modal.idx' + index)
        .modal('show')
        ;
    }
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
    clearStationMarkers: function () {
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
      this.addStationMarkers(map, stations);
    });
    Event.$on("clickStation", station => {
      console.log(this.showModal)
      // display modal
      this.showModal = true;

    });
    

  }
}); 