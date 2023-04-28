import { createApp } from 'vue'
import './style.scss'
import App from './App.vue'
import { router } from './router'
import { makeServer } from './server'
import VueApexCharts from "vue3-apexcharts";

if (process.env.NODE_ENV === "development") {
    makeServer()
}

createApp(App).use(router).use(VueApexCharts).mount('#app')
