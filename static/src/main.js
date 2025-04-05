import './assets/main.css'
import { createApp, ref } from 'vue'
import router from './router'
import App from './App.vue'
import 'mdui'
import 'mdui/mdui.css'

const app = createApp(App)
const title = ref('');
const login= ref(false);
const emitter = ref(0);
const page = ref(1);
const maxPage = ref(Infinity);

export const url = "http://localhost:2233/session"
// export const url = "/session"

app.provide("url", url)
app.provide("title", title)
app.provide("login", login)
app.provide("emitter", emitter)
app.provide("page", page)
app.provide("maxPage", maxPage)
app.use(router)
app.mount('#app')