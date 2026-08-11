import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import axios from 'axios'
import { getImageUrl } from './utils/imageUtils'

const app = createApp(App)

app.use(createPinia())
app.use(router)


app.config.globalProperties.$getImageUrl = getImageUrl


const token = localStorage.getItem('openpeo_token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

app.mount('#app')
