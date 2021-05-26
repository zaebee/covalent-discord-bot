<template>
  <b-card tag="article" class="border-0 mb-5 mb-md-4" no-body>
    <b-card-img-lazy
      v-if="isImgMeme"
      :src="meme.Url"
      :alt="`Meme of ${meme.Author.Username}`"
      top
    />
    <video
      v-else-if="isVideoMeme"
      :src="meme.Url"
      autoplay
      controls
      loop
      muted
    />
    <p v-else class="m-0 text-center p-2">
      {{ `This file type (.${getExt(meme.Url)}) is not supported!` }}
    </p>
    <b-card-body class="border-top px-3 py-2 position-relative">
      <b-card-title class="mb-0 text-break text-center">
        <NuxtLink
          :to="{ query: { ...$route.query, userName: meme.Author.Username } }"
          class="d-inline-flex align-items-center link-title text-left"
          ><b-avatar :src="avatar" alt="Avatar" class="mr-2" />@{{
            meme.Author.Username
          }}</NuxtLink
        >
      </b-card-title>
      <share
        :img-url="meme.Url"
        :user-name="meme.Author.Username"
        class="position-absolute card-share"
      />
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

const IMG_EXT = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp']

const VIDEO_EXT = [
  'mp4',
  'webm',
  'mpg',
  'mp2',
  'mpeg',
  'mpe',
  'mpv',
  'ogg',
  'mp4',
  'm4p',
  'm4v',
  'avi',
  'wmv',
  'mov',
  'qt',
  'flv',
  'swf',
]

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
    isImgMeme() {
      const ext = this.getExt(this.meme.Url)
      return IMG_EXT.includes(ext)
    },
    isVideoMeme() {
      const ext = this.getExt(this.meme.Url)
      return VIDEO_EXT.includes(ext)
    },
  },
  methods: {
    getExt(str) {
      return str.split('.').pop().toLowerCase()
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
