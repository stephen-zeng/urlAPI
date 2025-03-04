<script setup>
  import {inject} from 'vue'
  import Theme from "@/Components/Theme.vue";
  import Cookies from "js-cookie";
  import {snackbar} from "mdui";
  import {Notification, Post} from "@/fetch.js";
  import {sha256} from "js-sha256";

  const props = defineProps(['title'])
  const sidebarStatus = inject('sidebarStatus')
  const login = inject('login')
  const tab = inject('tab')
  const url = inject('url')
  const emitter = inject('emitter')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
  async function logout() {
    const session = await Post(url + "session", {
      "Token": Cookies.get("token"),
      "Send": {
        "operation": "logout",
        "login_term": false,
      }
    })
    if (session.error) {
      Notification(session.error)
    } else {
      Cookies.remove("token");
      login.value = false;
      Notification("Logout successful!");
    }
  }
</script>

<template>
  <mdui-top-app-bar>
    <mdui-button-icon icon="menu"
    @click="SidebarStatusChanged()"></mdui-button-icon>
    <mdui-top-app-bar-title>{{ props.title }}</mdui-top-app-bar-title>
    <div style="flex-grow: 1"></div>
    <mdui-button-icon @click="(emitter=1)" v-if="tab==0 && login" icon="refresh"></mdui-button-icon>
    <Theme></Theme>
    <mdui-button-icon @click="logout()" v-if="login"
                      icon="exit_to_app"></mdui-button-icon>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>