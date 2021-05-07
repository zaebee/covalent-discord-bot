<template>
  <b-col cols="4">
    <b-card
      tag="article"
      class="mb-4"
    >
      <b-img
        fluid
        :src="meme.Url"
        :alt="meme.Author.Username"
        style="overflow: hidden; height: 18rem;"
        class="mb-2 rounded-0"/>
      <b-card-title>@{{ meme.Author.Username }}</b-card-title>
      <b-card-text>2 hours ago</b-card-text>
    </b-card>
  </b-col>
</template>

<script lang="ts">
import Vue, { PropOptions } from 'vue'

interface User {
  ID: string
  Username: string
  Avatar: string
}

interface Meme {
  Url: string
  Author: User
}

export default Vue.extend({
  name: "Meme",
  props: {
    meme: {
      type: Object,
      required: true
    } as PropOptions<Meme>
  }
})
</script>
