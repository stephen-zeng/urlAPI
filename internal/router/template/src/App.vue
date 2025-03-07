<script setup>
import {ref, provide, inject, onUnmounted, onMounted, watch} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Task from "@/pages/Task.vue";
  import Backend from "@/pages/Backend.vue";
  import Tool from "@/pages/Tool.vue";
  import Workshop from "@/pages/Workshop.vue";
  import Login from "@/pages/Login.vue";
  import Cookies from "js-cookie";
  import {Notification, Post} from "@/fetch.js";
  import {useRouter} from "vue-router";

const sidebarStatus = ref(false);
  const pages = ref([
      '所有记录',
      '接口设置',
      '功能设置',
      '工作台',
  ])
  const login = inject("login");
  const url = inject("url");
  const router = useRouter();
  // 1 for header & access correspond, reset by access

  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);

  onMounted(async() => {
    if (Cookies.get("token")) {
      const session = await Post(url, {
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
    if (!login.value) {
      router.push("/login");
    }
  })

  onUnmounted(async() => {
    if (Cookies.get("token")) {
      const session = await Post(url, {
        "Token": Cookies.get("token"),
        "Send": {
          "operation": "exit",
          "login_term": false,
        }
      })
      login.value = false
    }
  })

  watch(login, (newValue, oldValue) => {
    if (!newValue) {
      router.push("/login");
    } else {
      router.push('/task')
    }
  })

</script>

<template>
  <mdui-layout full-height>
    <Header></Header>
    <Sidebar v-if="login"></Sidebar>
<!--    <Task v-if="tab==0 && login"></Task>-->
<!--    <Backend v-if="tab==1 && login"></Backend>-->
<!--    <Tool v-if="tab==2 && login"></Tool>-->
<!--    <Workshop v-if="tab==3 && login"></Workshop>-->
    <router-view></router-view>
<!--    <Login v-if="!login"></Login>-->
  </mdui-layout>

</template>

<style scoped>

</style>