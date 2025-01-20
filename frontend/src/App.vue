<script setup>
  import {ref, provide} from 'vue';
  import Header from "@/frameworks/Header.vue";
  import Sidebar from "@/frameworks/Sidebar.vue";
  import Welcome from "@/pages/Welcome.vue";

  const sidebarStatus = ref(false);
  const pages = ref([
      '任务队列',
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
    <Welcome v-if="tab==0"></Welcome>
  </mdui-layout>
</template>

<style scoped>

</style>