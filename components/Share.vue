<template>
  <div class="share">
    <button
      class="share__btn"
      :class="{ 'share__btn--active': isActive }"
      @click="share"
    >
      <b-icon-share-fill aria-hidden="true" />
    </button>
    <transition v-if="isActive" name="fade">
      <div class="share__list">
        <ShareNetwork
          v-for="network in networks"
          :network="network.network"
          :key="network.network"
          :style="{ backgroundColor: network.color }"
          :url="sharing.url"
          :title="sharing.title"
          :description="sharing.description"
          :quote="sharing.quote"
          :hashtags="sharing.hashtags"
          :twitterUser="sharing.twitterUser"
          class="share__item"
          :class="`share__item--${network.network}`"
        >
          <fa :icon="['fab', network.icon]" />
        </ShareNetwork>
      </div>
    </transition>
  </div>
</template>

<script>
import { BIconShareFill } from 'bootstrap-vue'

export default {
  name: 'Share',
  components: { BIconShareFill },
  props: {
    imgUrl: {
      type: String,
      required: true,
    },
    userName: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      isActive: false,
      sharing: {
        url: this.imgUrl,
        title: `Covalent meme of ${this.userName}`,
        description: '',
        quote: '',
        hashtags: 'Covalent,meme,DashboardMeme',
        twitterUser: this.userName,
      },
      networks: [
        {
          network: 'facebook',
          name: 'Facebook',
          icon: 'facebook-f',
        },
        {
          network: 'telegram',
          name: 'Telegram',
          icon: 'telegram-plane',
        },
        {
          network: 'twitter',
          name: 'Twitter',
          icon: 'twitter',
        },
      ],
    }
  },
  methods: {
    share() {
      this.isActive = !this.isActive
      console.log('share')
    },
  },
}
</script>

<style lang="scss" scoped>
.share {
  align-items: center;
  display: flex;
  flex-direction: column;

  &__btn {
    align-items: center;
    background-image: linear-gradient(to top, #000 50%, #fff 50%);
    background-size: 100% 200%;
    border: 1px solid #e5e5e5;
    border-radius: 50%;
    color: #000;
    display: flex;
    height: 2.25rem;
    justify-content: center;
    right: 1rem;
    transition: background-position 0.2s, border-color 0.2s;
    width: 2.25rem;

    &:hover,
    &--active {
      background-position: 0 100%;
      border-color: #000;
      color: #fff;
    }
  }

  &__list {
    bottom: 100%;
    margin-bottom: 0.75rem;
    position: absolute;
  }

  &__item {
    $color-facebook: #1877f2;
    $color-telegram: #08c;
    $color-twitter: #1da1f2;

    align-items: center;
    background-color: #fff;
    background-size: 100% 200%;
    border-radius: 50%;
    display: flex;
    height: 2rem;
    justify-content: center;
    transition: background-position 0.2s, color 0.2s;
    width: 2rem;

    &:hover {
      background-position: 0 100%;
      color: #fff;
    }

    & + & {
      margin-top: 0.5rem;
    }

    &--facebook {
      background-image: linear-gradient(to top, $color-facebook 50%, #fff 50%);
      color: $color-facebook;
    }

    &--telegram {
      background-image: linear-gradient(to top, $color-telegram 50%, #fff 50%);
      color: $color-telegram;
    }

    &--twitter {
      background-image: linear-gradient(to top, $color-twitter 50%, #fff 50%);
      color: $color-twitter;
    }
  }
}
</style>
