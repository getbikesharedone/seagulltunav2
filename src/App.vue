<template>
  <v-app id="example-1" toolbar footer>
    <v-navigation-drawer persistent v-model="overviewDrawer" light enable-resize-watcher absolute>
      <div :is="currentDrawer" :selectedStation="selectedStation"></div>
    </v-navigation-drawer>
    <v-toolbar class="indigo" dark fixed>
      <v-toolbar-side-icon @click.stop="overviewDrawer = !overviewDrawer"></v-toolbar-side-icon>
      <v-toolbar-title>Bike Network Location Finder</v-toolbar-title>
    </v-toolbar>
    <main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </main>
    <v-footer class="indigo">
      <span class="white--text">© 2017</span>
    </v-footer>
  </v-app>
</template>

<script>
import * as VueGoogleMaps from 'vue2-google-maps';
import Vue from 'vue';
import OverviewDrawer from '@/components/OverviewDrawer';
import ReviewsDrawer from '@/components/ReviewsDrawer';
import StarRating from 'vue-star-rating';
import EventBus from '@/event-bus';

Vue.component('star-rating', StarRating);

Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyDDy5IUrvL4bVAdeQ_MBvcqsy1rgs5X3V4',
  },
});

export default {
  data() {
    return {
      currentDrawer: 'OverviewDrawer',
      reviewsDrawer: 'ReviewsDrawer',
      overviewDrawer: 'OverviewDrawer',
      selectedStation: {},
      selectedNetwork: {},
    };
  },
  components: {
    OverviewDrawer,
    ReviewsDrawer,
  },
  created() {
    EventBus.$on('switchToReviews', (selectedStation) => {
      this.selectedStation = selectedStation;
      this.currentDrawer = 'ReviewsDrawer';
    });
    EventBus.$on('switchToOverview', () => {
      this.currentDrawer = 'OverviewDrawer';
    });
    EventBus.$on('addReview', (review) => {
      if (this.selectedStation.reviews !== undefined) {
        this.selectedStation.reviews[review.index].push(review);
      } else {
        this.selectedStation.reviews = [review];
      }
    });
    EventBus.$on('selectedNetwork', (selectedNetwork) => {
      this.selectedNetwork = selectedNetwork;
    });
  },
};
</script>
