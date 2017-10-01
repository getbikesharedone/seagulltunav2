<template>
  <v-app id="example-1" toolbar footer>
    <v-navigation-drawer persistent v-model="drawer" light enable-resize-watcher absolute>
      <v-layout row>
    <v-flex>
      <drawer></drawer>
    </v-flex>
  </v-layout>
    </v-navigation-drawer>
    <v-toolbar class="indigo" dark fixed>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
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
import Drawer from '@/components/Drawer';
import StarRating from 'vue-star-rating';

Vue.component('star-rating', StarRating);

Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyDDy5IUrvL4bVAdeQ_MBvcqsy1rgs5X3V4',
  },
});

export default {
  data() {
    return {
      drawer: true,
    };
  },
  components: { drawer: Drawer },
};
</script>
