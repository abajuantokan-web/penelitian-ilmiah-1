import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { getImageUrl } from './utils/imageUtils'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.config.globalProperties.$getImageUrl = getImageUrl

app.mount('#app')
