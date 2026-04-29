import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import router from './router'
import i18n from './i18n'
import App from './App.vue'
import './assets/styles/variables.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)
app.use(ElementPlus, { size: 'small' })
app.mount('#app')
