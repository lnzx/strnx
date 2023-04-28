<template>
    <div>
        <Nav />
        <div class="columns" style="height:95.8%">
          <Menu />
          <!-- main section -->
          <slot></slot>
        </div>
    </div>
</template>
    
<style scoped>
    .m-a {text-align:center; align-items:center; height: 3rem;}
    .m-li {text-align:center; align-items:center; border-left: 2px solid transparent;}
    .bg-color {background-color:#eef0f4;}
    .icon-color { color: #183153; }
    .active {
        border-left: 2px solid #2075E8;
    }
</style>
