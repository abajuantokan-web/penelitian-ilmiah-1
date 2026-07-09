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

// Provide image utility globally to all components
app.config.globalProperties.$getImageUrl = getImageUrl

// Initialize auth header if token exists
const token = localStorage.getItem('openpeo_token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

app.mount('#app')
