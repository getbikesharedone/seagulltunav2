<template>
  <v-card class="blue-grey darken-2 white--text">
              <v-card-title primary-title>
                <div class="headline">{{review.user}}</div>
                
              </v-card-title>
              <v-flex>{{review.body}}</v-flex>
              <v-flex><star-rating :rating="review.rating" :show-rating="false" :read-only="true"></star-rating></v-flex>
              <v-card-actions>
                  <v-dialog v-model="dialog" persistent width="50%">
                <v-btn flat dark slot="activator">Edit</v-btn>
                <v-card>
              <v-card-title>
                <span class="headline">Edit Review</span>
              </v-card-title>
              <v-card-text>
                <v-container grid-list-md>
                  <v-layout>
                      <v-flex>
                  <v-text-field label="Author" :value="review.user" :rules="authorRules" required readonly></v-text-field>
                </v-flex>
                  </v-layout>
                  <v-layout>
                      <v-flex>
                  <v-text-field multi-line label="Content" v-model="newBody" :value="review.body" :rules="contentRules" required></v-text-field>
                </v-flex>
                  </v-layout>
                  <v-layout>
                      <v-flex><star-rating :rating="review.rating" v-model="newRating"></star-rating></v-flex>
                  </v-layout>
                </v-container>
                <small>*indicates required field</small>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn class="blue--text darken-1" flat @click.native="dialog = false">Close</v-btn>
                <v-btn class="blue--text darken-1" flat @click.native="updateReview">Save</v-btn>
                <v-flex x </v-card-actions>
            </v-card>
                  </v-dialog>
              </v-card-actions>
            </v-card>
</template>

<script>
import StarRating from 'vue-star-rating';
import axios from 'axios';
import EventBus from '@/event-bus';

export default {
  data() {
    return {
      body: '',
      newBody: '',
      rating: 0,
      newRating: 0,
      user: '',
      authorRules: [
        v => !!v || 'Author is required',
      ],
      contentRules: [
        v => !!v || 'Body is required',
      ],
      dialog: false,
    };
  },
  props: ['review', 'index'],
  components: {
    StarRating,
  },
  methods: {
    updateReview() {
      axios
        .put(`/api/review/${this.review.id}`, {
          body: this.newBody,
          rating: this.newRating,
        })
        .then((res) => {
          const review = res.data;
          this.body = res.data.body;
          this.rating = res.data.rating;
          review.index = this.index;
          EventBus.$emit('updatedReview', review);
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
};
</script>
