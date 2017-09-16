window.Event = new Vue();

function addNetworkMarkers(map, networks) {
  for (let i = 0; i < networks.length; i++) {
    const network = networks[i]
    const marker = new google.maps.Marker({
      position: {
        lat: network.lat,
        lng: network.lng,
      },
      map,
      title: network.id,
      icon: {
        url: '/bike.png',
        size: new google.maps.Size(32, 32),
      },
      center: {
        lat: network.clat,
        lng: network.clng,
      },
      network,
    });
  }
}

const appVue = new Vue({
  el: "#app",
  data: {
    networks: []
  },
  created() {
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
  },
  methods: {
    initMap: function () {
      var myLatLng = { lat: 0, lng: 0 };

      var map = new google.maps.Map(document.getElementById('map'), {
        zoom: 3,
        center: myLatLng
      });

      addNetworkMarkers(map, this.networks);
    }
  },
  mounted () {
    Event.$on("networksLoaded", networks => {
      this.initMap()
    });
  }
}); 
