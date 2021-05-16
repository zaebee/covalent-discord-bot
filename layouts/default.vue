<template>
  <div id="app">
    <b-container class="text-center">
      <Logo />
      <h1 class="title mb-4">
        Memes
      </h1>
      <Navbar class="mb-4" />
    </b-container>
    <Nuxt />
  </div>
</template>

<script>
import Vue from 'vue'
import Logo from '~/components/Logo'
import Navbar from '~/components/Navbar'

export default Vue.extend({
  components: { Logo, Navbar },
})
</script>

<style lang="scss" scoped>
@font-face {
  font-family: 'What is this - some kind of joke';
  src: url('~/assets/fonts/What is this - some kind of joke.ttf')
    format('truetype');
}

.title {
  font-family: 'What is this - some kind of joke', sans-serif;
  font-size: 22vw;
  color: var(--color-covalent-gray);
  letter-spacing: 1px;
  text-transform: uppercase;

  @include media-breakpoint-up(sm) {
    font-size: 110px;
  }
}
</style>
