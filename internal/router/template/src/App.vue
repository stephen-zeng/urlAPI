<script setup>
import {ref, provide, inject, onUnmounted, onMounted} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Access from "@/pages/Access.vue";
  import Backend from "@/pages/Backend.vue";
  import Client from "@/pages/Client.vue";
  import Login from "@/pages/Login.vue";
  import Cookies from "js-cookie";
  import {Notification, Post} from "@/fetch.js";
x
  const sidebarStatus = ref(false);
  const pages = ref([
      '访问情况',
      '接口设置',
      '功能设置',
  ])
  const tab = ref(0);
  const login = ref(false);
  const url = inject("url");

  provide('tab', tab);
  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);
  provide('login', login);

  onMounted(async() => {
    if (Cookies.get("token")) {
      const session = await Post(url + "session", {
        "Token": Cookies.get("token"),
        "Send": {
          "operation": "login",
          "login_term": false,
        }
      })
      if (session.error) {
        Notification(session.error)
      } else {
        login.value = true
      }
    }
  })

  onUnmounted(async() => {
    if (Cookies.get("token")) {
      const session = await Post(url + "session", {
        "Token": Cookies.get("token"),
        "Send": {
          "operation": "exit",
          "login_term": false,
        }
      })
      login.value = false
    }
  })

</script>

<template>
  <mdui-layout full-height>
    <Header :title="login ? pages[tab] : '登录后台' "></Header>
    <Sidebar v-if="login"></Sidebar>
    <Access v-if="tab==0 && login"></Access>
    <Backend v-if="tab==1 && login"></Backend>
    <Client v-if="tab==2 && login"></Client>
    <Login v-if="!login"></Login>
  </mdui-layout>

</template>

<style scoped>

</style>