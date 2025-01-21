<script setup>
  import {ref, provide} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Access from "@/pages/Access.vue";
  import Backend from "@/pages/Backend.vue";
  import Client from "@/pages/Client.vue"

  const sidebarStatus = ref(false);
  const pages = ref([
      '访问情况',
      '后端API设置',
      '前端接口设置',
  ])
  const tab = ref(0);

  provide('tab', tab);
  provide('sidebarStatus', sidebarStatus);
  provide('pages', pages);

</script>

<template>
  <mdui-layout full-height>
    <Header :title="pages[tab]"></Header>
    <Sidebar></Sidebar>
    <Access v-if="tab==0"></Access>
    <Backend v-if="tab==1"></Backend>
    <Client v-if="tab==2"></Client>
  </mdui-layout>
</template>

<style scoped>

</style>