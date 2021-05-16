<template>
  <div class="border-bottom p-2">
    <b-nav pills>
      <b-nav-item
        to="/"
        class="nav-item mr-1"
        exact
        exact-active-class="nav-link--active"
        >New</b-nav-item
      >
      <b-nav-item
        to="top"
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
    color: #fff;
    border-radius: 20px;
    line-height: 15px;
    padding: 0.5rem 1.5rem;

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
