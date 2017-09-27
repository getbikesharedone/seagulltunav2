<template>
  <!-- The inline styles for the below div and gmap-map. They bypass a height bug, thus allowing
  Google Maps to be flex and 100% height. -->
  <div style="display: flex;
    min-height: 100%;
    flex-direction: column;">
    <gmap-map style="" :center="center" :zoom="zoom" style="position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;">
      <gmap-marker :key="index" v-for="(m, index) in networkMarkers" :position="m.position" :clickable="true" :draggable="true" @click="center=m.position"></gmap-marker>
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

export default {
  data() {
    return {
      center: { lat: 10.0, lng: 10.0 },
      networks: [],
      networkMarkers: [],
      zoom: 3,
    };
  },
  methods: {
    getNetworks() {
      Axios
        .get('/api/network')
        .then((res) => {
          this.networks = res.data;
          this.createNetworkMarkers();
        })
        .catch((error) => {
          console.log(error);
        });
    },
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
      });
    },
  },
  watch: {
    /* eslint-disable object-shorthand, no-unused-vars */
    // This prevents grey areas on the map.
    '$route'(to, from) {
      // Call resizePreserveCenter() on all maps
      Vue.$gmapDefaultResizeBus.$emit('resize');
    },
  },
  created() {
    this.getNetworks();
  },
};
</script>