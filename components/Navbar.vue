<template>
  <div class="border-bottom p-2">
    <b-nav pills>
      <b-nav-item
        v-if="$route.query.userName"
        :to="$route.query.sort ? { query: { sort: $route.query.sort } } : '/'"
        class="nav-item mr-1"
        exact
        exact-active-class="nav-link--active"
        >Back to all Users</b-nav-item
      >
      <b-nav-item
        :to="
          $route.query.userName
            ? { query: { userName: $route.query.userName } }
            : '/'
        "
        class="nav-item mr-1"
        exact
        exact-active-class="nav-link--active"
        >New</b-nav-item
      >
      <b-nav-item
        :to="{ query: { ...$route.query, sort: 'top' } }"
        class="nav-item"
        exact
        exact-active-class="nav-link--active"
        >Top</b-nav-item
      >
    </b-nav>
  </div>
</template>

<script>
export default {
  name: 'Navbar',
}
</script>

<style lang="scss" scoped>
.nav-pills {
  .nav-link {
    border-radius: 20px;
    color: #fff;
    line-height: 15px;
    padding: 0.5rem 1.5rem;
    transition: all 0.2s;

    &:hover {
      background-color: var(--color-covalent-gray);
      color: #000;
    }

    &--active {
      background-color: var(--color-covalent-pink);
      font-weight: 700;
    }
  }
}
</style>
