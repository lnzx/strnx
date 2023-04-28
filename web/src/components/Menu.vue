<template>
  <aside id="aside" class="column is-narrow menu pr-0">
    <ul class="menu-list is-flex-touch" id="menulist">
      <li class="m-li active"><router-link to="/dash" class="m-a icon-color" title="Dashboard"><i class="fa-solid fa-gauge-simple"></i></router-link></li>
      <li class="m-li"><router-link to="/nodes" class="m-a icon-color" title="Nodes"><i class="fas fa-server"></i></router-link></li>
      <li class="m-li"><router-link to="/wallets" class="m-a icon-color" title="Wallets"><i class="fa-solid fa-sack-dollar"></i></router-link></li>
    </ul>
  </aside>
</template>

<script>
export default {
    setup() {
        
    },
    mounted (){
      const $menulist = document.querySelectorAll('.m-li')
      $menulist.forEach(el => {
        el.addEventListener('click', () => {
          if(!el.classList.contains('active')){
            $menulist.forEach(e => {
              e.classList.remove('active')
            })
            el.classList.add('active')
          }
        })
      });
  }
}
</script>


<style scoped>
.m-a {text-align:center; align-items:center; height: 3rem;}
.m-li {text-align:center; align-items:center; border-left: 2px solid transparent;}
.icon-color { color: #183153; }
.active {
    border-left: 2px solid #2075E8;
}
</style>