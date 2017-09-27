<template>
  <!-- The inline styles for the below div and gmap-map. They bypass a height bug, thus allowing
  Google Maps to be flex and 100% height. -->
  <div style="display: flex;
    min-height: 100%;
    flex-direction: column;">
    <gmap-map style="" :center="center" :zoom="7" style="position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;">
      <gmap-marker :key="index" v-for="(m, index) in markers" :position="m.position" :clickable="true" :draggable="true" @click="center=m.position"></gmap-marker>
    </gmap-map>
  </div>
</template>

<script>
import * as VueGoogleMaps from 'vue2-google-maps';
import Vue from 'vue';

Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyDDy5IUrvL4bVAdeQ_MBvcqsy1rgs5X3V4',
  },
});

export default {
  data() {
    return {
      center: { lat: 10.0, lng: 10.0 },
      markers: [{
        position: { lat: 10.0, lng: 10.0 },
      }, {
        position: { lat: 11.0, lng: 11.0 },
      }],
    };
  },
};
</script>