<script setup>
import {ref, provide, inject, onUnmounted, onMounted, watch} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Cookies from "js-cookie";
  import {useRouter} from "vue-router";
import {Login} from "@/js/util.js";

const sidebarStatus = ref(false);
  const pages = ref([
      '所有记录',
      '接口设置',
      '功能设置',
      '安全设置',
      '工作台',
  ])
  const login = inject("login");
  const router = useRouter();
  // 1 for header & access correspond, reset by access

  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);

  onMounted(async() => {
    if (Cookies.get("token")) {
      login.value = await Login(Cookies.get("token"), false)
    }
    if (!login.value) {
      router.push("/dash/login");
    }
  })

  onUnmounted(async() => {
    if (Cookies.get("token")) {
      await Logout(Cookies.get("token"))
    }
  })

  watch(login, (newValue, oldValue) => {
    if (!newValue) {
      router.push("/dash/login");
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