<template>
  <b-container>
    <template v-if="getMemes.length">
      <b-row>
        <b-col v-for="meme in getMemes" :key="meme._id" sm="6" lg="4">
          <Meme :meme="meme._source" />
        </b-col>
      </b-row>
      <div class="text-center" v-show="hasNext">
        <b-overlay
          :show="busy"
          opacity="0.6"
          spinner-small
          spinner-variant="secondary"
          class="d-inline-block my-4"
        >
          <b-button
            ref="button"
            class="covalent-button-pink"
            :disabled="busy"
            @click="showMore"
          >
            More
          </b-button>
        </b-overlay>
      </div>
    </template>
    <b-alert v-else variant="light" class="text-center" show
      >No memes here yet</b-alert
    >
  </b-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { RootState } from '~/store'
import Meme from '~/components/Meme.vue'

export default Vue.extend({
  components: { Meme },
  async fetch({ store, query }) {
    const { sort, userName } = query
    await store.dispatch(sort === 'top' ? 'fetchTopMemes' : 'fetchMemes', {
      page: 0,
      userName,
    })
  },
  data() {
    return {
      busy: false,
      nextPage: 1,
    }
  },
  computed: {
    hasNext() {
      return this.$store.getters.hasNext as boolean
    },
    getMemes() {
      return (this.$store.state as RootState).memes
    },
  },
  methods: {
    showMore() {
      const { sort, userName } = this.$route.query
      this.busy = true
      this.$store
        .dispatch(sort === 'top' ? 'fetchTopMemes' : 'fetchMemes', {
          page: this.nextPage,
          userName,
        })
        .then(() => {
          this.nextPage += 1
          this.busy = false
        })
    },
  },
  watchQuery: ['sort', 'userName'],
})
</script>

<style lang="scss">
.covalent-button-pink {
  background-image: linear-gradient(
    to top,
    var(--color-covalent-turquoise) 50%,
    var(--color-covalent-pink) 50%
  );
  background-size: 100% 200%;
  border: 0;
  border-radius: $border-radius;
  color: #fff;
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.5rem;
  padding: 0.75rem 4rem;
  transition: background-position 0.2s;

  &:hover {
    background-position: 0 100%;
    color: #000426;
  }
}
</style>
