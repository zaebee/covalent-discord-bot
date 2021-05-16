<template>
  <b-card
    :img-src="meme.Url"
    :img-alt="`Meme of ${meme.Author.Username}`"
    img-top
    tag="article"
    class="border-0 mb-4"
    no-body
  >
    <b-card-body class="px-3 py-2">
      <b-card-title
        class="d-flex align-items-center justify-content-center mb-0 text-break"
      >
        <b-avatar :src="avatar" alt="Avatar" class="mr-2" />
        @{{ meme.Author.Username }}
      </b-card-title>
    </b-card-body>
    <b-card-footer footer-class="d-flex px-3 py-2 card-footer">
      <small class="text-muted"
        >Created {{ $moment(meme.Timestamp).fromNow() }}</small
      >
      <small class="ml-auto mr-0 text-muted text-right"
        >{{ meme.Reactions.length }}
        {{ meme.Reactions.length | pluralize('reaction') }}</small
      >
    </b-card-footer>
  </b-card>
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
  Timestamp: Date
  Reactions: string[]
}

export default Vue.extend({
  name: 'Meme',
  props: {
    meme: {
      type: Object,
      required: true,
    } as PropOptions<Meme>,
  },
  data() {
    return {}
  },
  computed: {
    avatar() {
      let userID = this.meme.Author.ID
      let avatar = this.meme.Author.Avatar
      return `https://cdn.discordapp.com/avatars/${userID}/${avatar}.png?size=64`
    },
  },
})
</script>
