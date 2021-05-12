<template>
  <b-container>
    <b-row>
      <b-col
        v-for="meme in getMemes"
        :key="meme._id"
        sm="6"
        lg="4"
      >
        <Meme
          :meme="meme._source"/>
      </b-col>
    </b-row>
    <div class="mt-4 text-center">
      <button
        class="covalent-button-pink"
      >More</button>
    </div>
  </b-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { RootState } from '~/store'
import Meme from "~/components/Meme.vue"

export default Vue.extend({
  name: 'New',
  components: { Meme },
  async fetch ({ store }) {
    await store.dispatch('fetchMemes')
  },
  computed: {
    getMemes () {
      return (this.$store.state as RootState).memes
    }
  }
})
</script>

<style lang="scss">
.covalent-button-pink {
  background-image: linear-gradient(to top, var(--color-covalent-turquoise) 50%, var(--color-covalent-pink) 50%);
  background-size: 100% 200%;
  border-radius: $border-radius;
  color: white;
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.5rem;
  margin-bottom: 2rem;
  padding: 0.75rem 4rem;
  transition: background-position 0.2s;

  &:hover {
     background-position: 0 100%;
     color: #000426;
   }
}
</style>
