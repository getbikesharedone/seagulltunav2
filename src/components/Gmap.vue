<template>
  <!-- The inline styles for the below div and gmap-map. They bypass a height bug, thus allowing
      Google Maps to be flex and 100% height. -->
  <div style="display: flex;
        min-height: 100%;
        flex-direction: column;">
    <gmap-map @zoom_changed="zoom = $event" @bounds_changed="bounds = $event" ref="map" style="" :center="center" :zoom="zoom" style="position: absolute;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;">
      <google-cluster ref="networkCluster">
        <gmap-marker ref="networkMarkers" :key="index" v-for="(m, index) in networkMarkers" :position="m.position" :clickable="true" :draggable="true" @click="createStationMarkers(m)"></gmap-marker>
      </google-cluster>
      <gmap-marker ref="stationMarkers" :key="index" v-for="(m, index) in stationMarkers" :position="m.position" :clickable="true" :draggable="true" @click="selectStation(m)"></gmap-marker>
    </gmap-map>

  </div>
</template>

<script>
import * as VueGoogleMaps from 'vue2-google-maps';
import Vue from 'vue';
import Axios from 'axios';

Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyDDy5IUrvL4bVAdeQ_MBvcqsy1rgs5X3V4',
  },
});
Vue.component('google-cluster', VueGoogleMaps.Cluster);
/* global google */

export default {
  data() {
    return {
      bounds: {},
      center: { lat: 10.0, lng: 10.0 },
      created: false,
      networks: [],
      networkMarkers: [],
      selectedNetwork: {},
      selectedStation: {},
      stationMarkers: [],
      stations: [],
      zoom: 3,
    };
  },
  methods: {
    /* eslint-disable no-unused-expressions */
    createStationMarkers(selectedNetworkMarker) {
      this.selectedNetwork = selectedNetworkMarker;
      // this.hideNetworkMarkers(); // causes slowness when re-adding markers
      this.getStations.then(() => {
        this.createSMarkers;
        this.fitStationBounds;
      });
    },
    selectStation(station) {
      this.selectedStation = station;
    },
  },
  watch: {
    /* eslint-disable object-shorthand, no-unused-vars */
    // This prevents grey areas on the map.
    '$route'(to, from) {
      // Call resizePreserveCenter() on all maps
      Vue.$gmapDefaultResizeBus.$emit('resize');
    },
    zoom(newZoom) {
      if (this.created === true && newZoom <= 10) {
        if (this.stationMarkers.length !== 0) {
          this.hideStationMarkers;
        }
        // this.createNetworkMarkers(); // used in conjunction with other same named method
        this.created = false;
      }
    },
  },
  created() {
    this.getNetworks;
  },
  computed: {
    createNetworkMarkers() {
      this.networks.forEach((network) => {
        const marker = {
          title: network.name,
          position: {
            lat: network.lat, lng: network.lng,
          },
          id: network.id,
        };
        this.networkMarkers.push(marker);
        this.created = true;
      });
    },
    createSMarkers() {
      this.stations.forEach((station) => {
        const marker = {
          empty: station.empty,
          free: station.free,
          open: station.open,
          safe: station.safe,
          time: station.time,
          title: station.name,
          position: {
            lat: station.lat, lng: station.lng,
          },
          id: station.id,
        };
        this.stationMarkers.push(marker);
      });
    },
    fitStationBounds() {
      const bounds = new google.maps.LatLngBounds();
      this.stationMarkers.forEach((marker) => {
        bounds.extend(marker.position);
      });
      this.$refs.map.fitBounds(bounds);
      this.zoom -= 1; // Remove one zoom level to ensure no marker is on the edge.
      /* Set a minimum zoom.
      If you got only 1 marker, or all markers are on the same address,
      the map will be zoomed too much. */
      if (this.zoom > 15) {
        this.zoom = 15;
      }
    },
    getNetworks() {
      Axios
        .get('/api/network')
        .then((res) => {
          this.networks = res.data;
          this.createNetworkMarkers;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getStations() {
      return Axios
        .get(`/api/network/${this.selectedNetwork.id}`)
        .then((res) => {
          this.stations = res.data.stations;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    hideNetworkMarkers() {
      this.$refs.networkCluster.$clusterObject.clearMarkers();
    },
    hideStationMarkers() {
      /* eslint-disable no-param-reassign */
      this.stationMarkers = [];
    },
  },
};
</script>