<template>
  <div class="container">
    <div>
      <Logo />
      <h1 class="title">
        memes
      </h1>
      <b-container class="mt-4">
        <b-row>
          <Meme v-for="meme in getMemes" :key="meme._id" :meme="meme._source"></Meme>
        </b-row>
        <b-row>
        </b-row>
      </b-container>
      <b-container id="more">
        <div class="links">
            <a
              href="#more"
              class="btn button--green"
            >More</a>
      </div>
      </b-container>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { RootState } from '~/store'
import Logo from '~/components/Logo.vue'
import Meme from "~/components/Meme.vue";

export default Vue.extend({
  components: {Meme, Logo},
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

<style>
.container {
  margin: 0 auto;
  min-height: 100vh;
  /*display: flex;*/
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family:
    'Quicksand',
    'Source Sans Pro',
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    'Helvetica Neue',
    Arial,
    sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.links {
  padding-top: 15px;
}
</style>
