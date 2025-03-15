<script setup>
  import {inject} from 'vue'
  import Theme from "@/frameworks/Theme.vue";
  import Cookies from "js-cookie";
  import {Notification, Post} from "@/fetch.js";
  import {useRoute} from "vue-router";

  const title = inject("title");
  const sidebarStatus = inject('sidebarStatus')
  const login = inject('login')
  const url = inject('url')
  const emitter = inject('emitter')
  const router = useRoute()
  const page = inject('page')
  const maxPage = inject('maxPage')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
  async function logout() {
    const session = await Post(url, {
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
    <mdui-top-app-bar-title>{{ title }}</mdui-top-app-bar-title>
    <div style="flex-grow: 1"></div>
<!--    1 for refresh, 2 for backwards, 3 for forwards-->
    <mdui-segmented-button-group v-if="login && router.path == '/dash/task'">
      <mdui-segmented-button @click="(emitter=2)">←</mdui-segmented-button>
      <mdui-segmented-button>{{ page }} / {{ maxPage }}</mdui-segmented-button>
      <mdui-segmented-button @click="(emitter=3)">→</mdui-segmented-button>
    </mdui-segmented-button-group>
    <mdui-button-icon @click="(emitter=1)" v-if="login && router.path == '/dash/task'" icon="refresh"></mdui-button-icon>

    <Theme v-if="router.path != '/dash/task'"></Theme>
    <mdui-button-icon @click="logout()" v-if="login"
                      icon="exit_to_app"></mdui-button-icon>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>