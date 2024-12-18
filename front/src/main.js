import './assets/main.css'
import 'vue3-toastify/dist/index.css';

import Vue3Toastify from 'vue3-toastify';

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Vue3Toastify, {
    theme: "dark",
    autoClose: 1000
})

app.mount('#app')
