<template>
<nav id="navbar" class="navbar" role="navigation" aria-label="main navigation">
  <div class="navbar-brand">
    <a class="navbar-item">
      <i class="fa-solid fa-hat-cowboy fa-lg icon-color"></i>
    </a>
    
    <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="navbarmenu">
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
    </a>
  </div>

  <div class="navbar-menu" id="navbarmenu">
    <div class="navbar-end">
      <div class="navbar-item has-dropdown" id="dropdownId">
        <a class="navbar-link">
          <span class="icon-text">
            <span class="icon icon-color">
              <i class="fa-solid fa-user-secret"></i>
            </span>
            <span class="profile">{{ username }}</span>
          </span>
        </a>
        <div class="navbar-dropdown is-right is-boxing">
          <router-link to="/settings" class="navbar-item fix-navbar-item" @click="setting">
            <span class="icon-text">
              <span class="icon icon-color">
                <i class="fa-solid fa-gear"></i>
              </span>
              <span>Settings</span>
            </span>
          </router-link>
          <hr class="navbar-divider">
          <a class="navbar-item fix-navbar-item" @click="logout">
            <span class="icon-text">
              <span class="icon icon-color">
                <i class="fa-solid fa-right-from-bracket"></i>
              </span>
              <span>Log out</span>
            </span>
          </a>
        </div>
      </div>
    </div>
  </div>
</nav>
</template>

<script>
const api = useApi()

export default {
    setup() {
      const token = useSession.getToken()
      const logout = () => {
        api.post('/api/logout').then(res => {
          useSession.clearToken()
          location.href = '/'
        })
      }

      const setting = () => {
        const $menulist = document.querySelectorAll('.m-li')
        $menulist.forEach(el => {
          if(el.classList.contains('active')){
            el.classList.remove('active')
          }
        });
      }

      return {logout, setting, username: ref(token.username)}
    },

    mounted(){
        document.addEventListener('click', (e) => {
        const $classList = e.target.classList;
        if(!$classList.contains('navbar-link') && !$classList.contains('fa-user-secret') && !$classList.contains('profile')){
          const $dropdown = document.getElementById('dropdownId').classList
          if($dropdown.contains('is-active')){
            $dropdown.remove('is-active')
          }
        }
      });

      // Get all "navbar-burger" elements
      const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);
      // Check if there are any navbar burgers
      if ($navbarBurgers.length > 0) {

        // Add a click event on each of them
        $navbarBurgers.forEach( el => {
          el.addEventListener('click', () => {

            // Get the target from the "data-target" attribute
            const target = el.dataset.target;
            const $target = document.getElementById(target);

            // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
            el.classList.toggle('is-active');
            $target.classList.toggle('is-active');

          });
        });
      }

      const $target = document.getElementById('dropdownId');
      $target.addEventListener('click', () => {
        $target.classList.toggle('is-active');
      });
    }
}
</script>


<style scoped>
.icon-color { color: #183153; }
.navbar-dropdown a.navbar-item {
    padding-right: 1.3rem;
}
.icon-text {
    flex-wrap: nowrap;
}
</style>