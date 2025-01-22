<script setup>
  import {ref, provide} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Access from "@/pages/Access.vue";
  import Backend from "@/pages/Backend.vue";
  import Client from "@/pages/Client.vue";
  import Login from "@/pages/Login.vue";

  const sidebarStatus = ref(false);
  const pages = ref([
      '访问情况',
      '接口设置',
      '功能设置',
  ])
  const tab = ref(0);
  const login = ref(true);

  provide('tab', tab);
  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);

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