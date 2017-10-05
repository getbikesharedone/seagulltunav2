<template>
  <v-app id="e3" style="max-width: 400px; margin: auto;" standalone>
    <v-toolbar class="pink">
      <v-toolbar-title class="white--text">Reviews</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon>add</v-icon>
      </v-btn>
      <v-btn dark icon @click="switchToOverview()">
                <v-icon>more_vert</v-icon>
              </v-btn>
    </v-toolbar>
    <main>
      <v-container
        fluid
        style="min-height: 0;"
        grid-list-lg
      >
        <v-layout row wrap>
          <v-flex xs12 v-for="review in reviews" :key="review.id">
            <v-card class="blue-grey darken-2 white--text">
              <v-card-title primary-title>
                <div class="headline">Unlimited music now</div>
                <div>Listen to your favorite artists and albums whenenver and wherever, online and offline.</div>
              </v-card-title>
              <v-card-actions>
                <v-btn flat dark>Listen now</v-btn>
              </v-card-actions>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
    </main>
  </v-app>
</template>

<script>
import axios from 'axios';
import EventBus from '@/event-bus';

export default {
  /* eslint-disable no-unused-expressions */
  data() {
    return {
      thumbnail320wSrc: '',
      thumbnail640wSrc: '',
      thumbnail1024wSrc: '',
      thumbnail2048wSrc: '',
      imageKey: '',
      reviews: [],
    };
  },
  props: ['selectedStation'],
  methods: {
    getImagesSrc() {
      this.thumbnail320wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail640wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail1024wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail2048wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
    },
    switchToOverview() {
      EventBus.$emit('switchToOverview');
    },
  },
  computed: {
    getSrcset() {
      const srcSet = `${this.thumbnail640wSrc} 640w, ${this.thumbnail1024wSrc} 1024w, ${this.thumbnail2048wSrc} 2048w`;
      return srcSet;
    },
    getStationImageFeature() {
      axios
        .get(this.imageUrl)
        .then((res) => {
          this.imageKey = res.data.features[0].properties.key;
          this.getImagesSrc();
        })
        .catch((error) => {
          console.log(error);
        });
    },
    imageUrl() {
      const lat = this.selectedStation.position.lat;
      const lng = this.selectedStation.position.lng;
      console.log(`lat: ${lat}, lng: ${lng}`);
      const baseUrl = 'https://a.mapillary.com/v3/images?';
      const clientId = 'client_id=SHNGU2JaY3Z4M3hEMWd5eW1CNElhQTowM2FhZjZhZWIyYmVkOTY0';
      const lookAt = `&lookat=${lat},${lng}`;
      const finalUrl = baseUrl + clientId + lookAt;
      return finalUrl;
    },
    photoPageUrl() {
      return `https://www.mapillary.com/app/?focus=photo&pKey=${this.key}`;
    },
  },
};
</script>