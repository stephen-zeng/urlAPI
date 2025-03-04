<script setup>
import {ref, provide, inject, onUnmounted, onMounted} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Access from "@/pages/Access.vue";
  import Backend from "@/pages/Backend.vue";
  import Client from "@/pages/Client.vue";
  import Workshop from "@/pages/Workshop.vue";
  import Login from "@/pages/Login.vue";
  import Cookies from "js-cookie";
  import {Notification, Post} from "@/fetch.js";

  const sidebarStatus = ref(false);
  const pages = ref([
      '所有记录',
      '接口设置',
      '功能设置',
      '工作台',
  ])
  const tab = ref(0);
  const tabAddition = ref("");
  const login = ref(false);
  const url = inject("url");
  const emitter = ref(0);
  // 1 for header & access correspond, reset by access

  provide('tab', tab);
  provide('tabAddition', tabAddition);
  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);
  provide('login', login);
  provide('emitter', emitter);

  onMounted(async() => {
    if (Cookies.get("tab")) {
      tab.value = Cookies.get("tab");
    }
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
    <Header :title="login ? pages[tab] + tabAddition : '登录后台' "></Header>
    <Sidebar v-if="login"></Sidebar>
    <Access v-if="tab==0 && login"></Access>
    <Backend v-if="tab==1 && login"></Backend>
    <Client v-if="tab==2 && login"></Client>
    <Workshop v-if="tab==3 && login"></Workshop>
    <Login v-if="!login"></Login>
  </mdui-layout>

</template>

<style scoped>

</style>