<template>
  <b-img
    alt="Covalent"
    class="logo"
    fluid
    src="~/assets/img/logo.png"
    height="52"
    width="220"
  />
</template>

<style lang="scss" scoped>
.logo {
  animation: 1s appear;
  height: 52px;
  margin: 40px auto;
  width: 220px;
}

@keyframes appear {
  0% {
    opacity: 0;
  }
}
</style>
