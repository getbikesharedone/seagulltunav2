<template>
  <v-app id="e3" style="max-width: 400px; margin: auto;" standalone>
    <v-toolbar class="pink">
      <v-toolbar-title class="white--text">Reviews</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-dialog v-model="dialog" persistent width="50%">
      <v-btn icon slot="activator">
        <v-icon>add</v-icon>
      </v-btn>
      <v-card>
              <v-card-title>
                <span class="headline">Add Review</span>
              </v-card-title>
              <v-card-text>
                <v-container grid-list-md>
                  <v-layout>
                      <v-flex>
                  <v-text-field label="Author" v-model="user" :rules="authorRules" required></v-text-field>
                </v-flex>
                  </v-layout>
                  <v-layout>
                       <v-flex>
                  <v-text-field multi-line label="Content" v-model="body" :rules="contentRules" required></v-text-field>
                </v-flex>
                  </v-layout>
                  <v-layout>
                      <v-flex><star-rating v-model="rating"></star-rating></v-flex>
                  </v-layout>
                </v-container>
                <small>*indicates required field</small>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn class="blue--text darken-1" flat @click.native="dialog = false">Close</v-btn>
                <v-btn class="blue--text darken-1" flat @click.native="saveReview">Save</v-btn>
                <v-flex x </v-card-actions>
            </v-card>
      </v-dialog>
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
          <v-flex v-for="(review, index) in reviews" :key="review.id">
            <review :review="review" :index="index"></review>
          </v-flex>
        </v-layout>
      </v-container>
    </main>
  </v-app>
</template>

<script>
import axios from 'axios';
import EventBus from '@/event-bus';
import StarRating from 'vue-star-rating';
import Review from '@/components/Review';


export default {
  /* eslint-disable no-unused-expressions */
  data() {
    return {
      reviews: [],
      dialog: false,
      user: '',
      body: '',
      authorRules: [
        v => !!v || 'Author is required',
      ],
      contentRules: [
        v => !!v || 'Body is required',
      ],
      rating: 0,
    };
  },
  components: {
    StarRating,
    Review,
  },
  props: ['selectedStation'],
  methods: {
    switchToOverview() {
      EventBus.$emit('switchToOverview');
    },
    saveReview() {
      axios
        .post(`/api/station/${this.selectedStation.id}/review`, {
          user: this.user,
          body: this.body,
          rating: this.rating,
        })
        .then((res) => {
          const review = {
            user: res.data.user,
            body: res.data.body,
            rating: res.data.rating,
          };
          this.reviews.push(review);
          review.index = this.selectedStation.index;
          EventBus.$emit('addReview', review);
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
  created() {
    axios
      .get(`/api/station/${this.selectedStation.id}`)
      .then((res) => {
        if (res.data.reviews !== undefined) {
          this.reviews = res.data.reviews;
        }
      })
      .catch((error) => {
        console.log(error);
      });

    EventBus.$on('updatedReview', (review) => {
      this.$set(this.reviews, review.index, review);
    });
  },
};
</script>