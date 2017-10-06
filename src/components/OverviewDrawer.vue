<template>
  <v-card>
    <v-card-media :src="thumbnail320wSrc">
      <v-layout column class="media">
        <v-card-title>
          <v-spacer></v-spacer>
          <v-dialog v-model="dialog" persistent width="50%">
            <v-btn dark icon class="mr-3" slot="activator">
              <v-icon>edit</v-icon>
            </v-btn>
            <v-card>
              <v-card-title>
                <span class="headline">Station Settings</span>
              </v-card-title>
              <v-card-text>
                <v-container grid-list-md>
                  <v-layout wrap>
                    <v-checkbox label="Open" v-model="newOpen"></v-checkbox>
                    <v-checkbox label="Safe" v-model="newSafe"></v-checkbox>
                    <v-flex>
                  <v-text-field label="Available Bikes" v-model.number="newFree" :rules="freeBikeRules" required></v-text-field>
                </v-flex>
                  </v-layout>
                </v-container>
                <small>*indicates required field</small>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn class="blue--text darken-1" flat @click.native="dialog = false">Close</v-btn>
                <v-btn class="blue--text darken-1" flat @click.native="saveSettings">Save</v-btn>
                <v-flex x </v-card-actions>
            </v-card>
          </v-dialog>
          <v-btn dark icon @click="switchToReviews()">
            <v-icon>more_vert</v-icon>
          </v-btn>
        </v-card-title>
        <v-spacer></v-spacer>
        <v-card-title class="white--text pl-5 pt-5">
          <div class="">{{selectedStation.title}}</div>
        </v-card-title>
      </v-layout>
    </v-card-media>
    <v-list two-line>
      <v-list-tile>
        <v-list-tile-action></v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>{{free}}</v-list-tile-title>
          <v-list-tile-sub-title>Free</v-list-tile-sub-title>
        </v-list-tile-content>
        <v-list-tile-action>
        </v-list-tile-action>
      </v-list-tile>
      <v-list-tile>
        <v-list-tile-action>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>{{open}}</v-list-tile-title>
          <v-list-tile-sub-title>Open</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
      <v-list-tile>
        <v-list-tile-action></v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>{{safe}}</v-list-tile-title>
          <v-list-tile-sub-title>Safe</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list>
  </v-card>
</template>

<script>
import EventBus from '@/event-bus';
import axios from 'axios';

export default {
  /* eslint-disable no-unused-expressions */
  data() {
    return {
      selectedStation: { position: { lat: 0, lng: 0 } },
      dialog: false,
      freeBikeRules: [
        v => !!v || 'Available Bikes is required',
        v => v > -1 || 'Number must be non-negative',
      ],
      free: 0,
      empty: 0,
      safe: true,
      open: true,
      total: 0,
      thumbnail320wSrc: '',
      thumbnail640wSrc: '',
      thumbnail1024wSrc: '',
      thumbnail2048wSrc: '',
      imageKey: '',
      safeCheckbox: false,
      openCheckbox: false,
      newFree: 0,
      newOpen: false,
      newSafe: false,
    };
  },
  methods: {
    switchToReviews() {
      EventBus.$emit('switchToReviews', this.selectedStation);
    },
    getImagesSrc() {
      this.thumbnail320wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-320.jpg`;
      this.thumbnail640wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-640.jpg`;
      this.thumbnail1024wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-1024.jpg`;
      this.thumbnail2048wSrc = `https://d1cuyjsrcm0gby.cloudfront.net/${this.imageKey}/thumb-2048.jpg`;
    },
    saveSettings() {
      axios
        .post(`/api/station/${this.selectedStation.id}`, {
          id: this.selectedStation.id,
          free: this.newFree,
          open: this.newOpen,
          safe: this.newSafe,
        })
        .then((res) => {
          this.newFree = res.data.free;
          this.newOpen = res.data.open;
          this.newSafe = res.data.safe;
          this.free = res.data.free;
          this.open = res.data.open;
          this.safe = res.data.safe;
          this.selectedStation.free = this.newFree;
          this.selectedStation.open = this.newOpen;
          this.selectedStation.safe = this.newSafe;
          EventBus.$emit('saveSettings', this.selectedStation);
        })
        .catch((error) => {
          console.log(error);
        });
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
      this.free = selectedStation.free;
      this.empty = selectedStation.empty;
      this.safe = selectedStation.safe;
      this.open = selectedStation.open;
      this.newEmpty = selectedStation.empty;
      this.newSafe = selectedStation.safe;
      this.newOpen = selectedStation.open;
    });
  },
};
</script>