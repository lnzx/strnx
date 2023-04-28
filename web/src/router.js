import { createRouter, createWebHistory } from 'vue-router/auto'

export const router = createRouter({
    history: createWebHistory(),
    // You don't need to pass the routes anymore, the plugin writes it for you 🤖
})

router.beforeEach((to, from) => {
    if(to.meta.requiresAuth && !useSession.isLogin()){
        // 将用户重定向到login页面
        return {
            path: '/login',
            // save the location we were at to come back later
            query: { redirect: to.fullPath },
        }
    }
    return true
})


