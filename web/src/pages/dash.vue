<route lang="yaml">
meta:
    requiresAuth: true
</route>

<template>
    <div style="height:95.8%">
        <Nav />
        <div class="columns" style="height:95.8%">
          <Menu />
          <!-- main section -->
          <router-view /> 
        </div>
    </div>
</template>


<style scoped>

.bg-color {background-color:#eef0f4;}

</style>
