import './assets/main.css'
import { createApp } from 'vue'
import App from './App.vue'
import 'mdui'
import 'mdui/mdui.css'

// const url = location.href
const url = "http://localhost:2233/dash"
const app = createApp(App)
app.provide("url", url)
app.mount('#app')