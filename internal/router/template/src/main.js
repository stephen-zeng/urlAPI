import './assets/main.css'
import { createApp, ref } from 'vue'
import router from './router'
import App from './App.vue'
import 'mdui'
import 'mdui/mdui.css'

const app = createApp(App)
const url = "http://localhost:2233/session"
// const url = "/session"
const title = ref('');
const login= ref(false);
const emitter = ref(0);

app.provide("url", url)
app.provide("title", title)
app.provide("login", login)
app.provide("emitter", emitter)
app.use(router)
app.mount('#app')