<template>
  <b-card
    :img-src="meme.Url"
    :img-alt="`Meme of ${meme.Author.Username}`"
    img-top
    tag="article"
    class="border-0 mb-4"
    no-body
  >
    <b-card-body class="px-3 py-2 position-relative">
      <b-card-title class="mb-0 text-break text-center">
        <NuxtLink
          :to="{ query: { ...$route.query, userName: meme.Author.Username } }"
          class="d-inline-flex align-items-center link-title text-left"
          ><b-avatar :src="avatar" alt="Avatar" class="mr-2" />@{{
            meme.Author.Username
          }}</NuxtLink
        >
      </b-card-title>
      <share :img-url="meme.Url" :user-name="meme.Author.Username" class="position-absolute card-share" />
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
import Share from '~/components/Share.vue'

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
  components: { Share },
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
      const userID = this.meme.Author.ID
      const avatar = this.meme.Author.Avatar
      return avatar
        ? `https://cdn.discordapp.com/avatars/${userID}/${avatar}.png?size=64`
        : ''
    },
  },
})
</script>

<style lang="scss" scoped>
.card-img-top {
  max-height: 500px;
  object-fit: contain;
}

.card-share {
  right: 0.5rem;
  top: -1.125rem;
}

.link-title {
  color: #000;
  transition: color 0.2s;

  &:hover {
    color: var(--color-covalent-pink);
    text-decoration: none;
  }
}
</style>
