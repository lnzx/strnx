<template>
  <section class="columns is-centered is-vcentered is-vfull">
    <div class="column is-5">
      <figure class="image">
        <img src="/station.svg">
      </figure>
    </div>
    <div class="column is-3">
      <fieldset>
        <div class="field">
          <h4 class="title is-4">Log in to StrnX</h4>
        </div>

        <div class="field">
          <label class="label c-label"></label>
          <div class="control has-icons-left">
            <input class="input" type="text" placeholder="Username" v-model="username" @keydown.enter="login">
            <span class="icon is-small is-left">
              <i class="fa-solid fa-user-secret"></i>
            </span>
          </div>
        </div>

        <div class="field">
          <label class="label c-label"></label>
          <div class="control has-icons-left">
            <input class="input" type="password" placeholder="Password" v-model="password" @keydown.enter="login">
            <span class="icon is-small is-left">
              <i class="fa-solid fa-key"></i>
            </span>
          </div>
        </div>

        <div class="field">
          <div class="control">
            <label class="checkbox">
              <input type="checkbox" v-model="remember">
              Remember me
              <span class="ml-3 has-text-danger" :class="{'is-hidden': errHide}">Credentials do not match</span>
            </label>
          </div>
        </div>

        <div class="field">
          <div class="control">
            <button class="button is-link is-fullwidth" style="background-color: #5769df;" @click="login" >Log in</button>
          </div>
        </div>

      </fieldset>
    </div>
  </section>
</template>

<script setup>
  const username = ref('')
  const password = ref('')
  const remember = ref(false)
  const errHide = ref(true)
  const router = useRouter()
  const route = useRoute()
  const redirect = route.query.redirect
  const api = useApi()

  const login = () => {
    if(username.value && password.value){
      api.post('/api/login', {username:username.value, password:password.value}).then(res => {
        const data = res.data
        if(data.token){
          useSession.setToken(data, remember.value)
          if(redirect){
            router.push(redirect)
          }else{
            router.push({
              path: '/dash',
            })
          }
        }else{
          errHide.value = false
        }
      })
    }
  }
</script>

<style scoped>
.input {box-shadow:none;}
.c-label {font-size: .875rem;}
</style>