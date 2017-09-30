<template>
  <v-card>
        <v-card-media :src="thumbnail320wSrc">
          <v-layout column class="media">
            <v-card-title>
              <v-btn dark icon>
                <v-icon>chevron_left</v-icon>
              </v-btn>
              <v-spacer></v-spacer>
              <v-btn dark icon class="mr-3">
                <v-icon>edit</v-icon>
              </v-btn>
              <v-btn dark icon>
                <v-icon>more_vert</v-icon>
              </v-btn>
            </v-card-title>
            <v-spacer></v-spacer>
            <v-card-title class="white--text pl-5 pt-5">
              <div class="display-1 pl-5 pt-5">Ali Conners</div>
            </v-card-title>
          </v-layout>
        </v-card-media>
        <v-list two-line>
          <v-list-tile @click="">
            <v-list-tile-action>
              <v-icon class="indigo--text">phone</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>(650) 555-1234</v-list-tile-title>
              <v-list-tile-sub-title>Mobile</v-list-tile-sub-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-icon dark>chat</v-icon>
            </v-list-tile-action>
          </v-list-tile>
          <v-list-tile @click="">
            <v-list-tile-action></v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>(323) 555-6789</v-list-tile-title>
              <v-list-tile-sub-title>Work</v-list-tile-sub-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-icon dark>chat</v-icon>
            </v-list-tile-action>
          </v-list-tile>
        <v-divider inset></v-divider>
          <v-list-tile @click="">
            <v-list-tile-action>
              <v-icon class="indigo--text">mail</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>aliconnors@example.com</v-list-tile-title>
              <v-list-tile-sub-title>Personal</v-list-tile-sub-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile @click="">
            <v-list-tile-action></v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>ali_connors@example.com</v-list-tile-title>
              <v-list-tile-sub-title>Work</v-list-tile-sub-title>
            </v-list-tile-content>
          </v-list-tile>
        <v-divider inset></v-divider>
          <v-list-tile @click="">
            <v-list-tile-action>
              <v-icon class="indigo--text">location_on</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>1400 Main Street</v-list-tile-title>
              <v-list-tile-sub-title>Orlando, FL 79938</v-list-tile-sub-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list>
      </v-card>
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
      selectedStation: { position: { lat: 0, lng: 0 } },
      imageKey: '',
    };
  },
  methods: {
    getImagesSrc() {
      this.thumbnail320wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail640wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail1024wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail2048wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
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
  created() {
    EventBus.$on('stationSelected', (selectedStation) => {
      this.selectedStation = selectedStation;
      this.getStationImageFeature;
    });
  },
};
</script>