<script setup>
  import {inject} from 'vue'
  import Theme from "@/Components/Theme.vue";
  import Cookies from "js-cookie";
  import {snackbar} from "mdui";

  const props = defineProps(['title'])
  const sidebarStatus = inject('sidebarStatus')
  const login = inject('login')
  const url = inject('url')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
  function logout() {
    fetch(url+"session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": Cookies.get("token"),
      },
      body: JSON.stringify({
        operation: "logout",
      })
    }).then(res => res.json()).then((data) => {
      if (data.error) {
        snackbar({
          message: data.error,
          placement: "top-end",
        })
      } else {
        login.value = false;
        Cookies.remove("token");
        snackbar({
          message: "Logged out",
          placement: "top-end",
        })
      }
    })
  }
</script>

<template>
  <mdui-top-app-bar>
    <mdui-button-icon icon="menu"
    @click="SidebarStatusChanged()"></mdui-button-icon>
    <mdui-top-app-bar-title>{{ props.title }}</mdui-top-app-bar-title>
    <div style="flex-grow: 1"></div>
    <Theme></Theme>
    <mdui-button-icon @click="logout()" v-if="login"
                      icon="exit_to_app"></mdui-button-icon>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>