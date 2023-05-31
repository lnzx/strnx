<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <div>
    <Nav />
    <div class="columns">
      <Menu />
      <!-- main section -->
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.bg-color {
  background-color: #eef0f4;
}
</style>
