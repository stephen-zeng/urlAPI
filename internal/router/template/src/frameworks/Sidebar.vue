<script setup>
  import { inject } from 'vue';
  import Cookies from "js-cookie";

  const openStatus = inject('sidebarStatus');
  const tab = inject('tab');
  const pages = inject('pages');

  function sidebarClose() {
    openStatus.value = false;
  }
  function changeTab(newTab) {
    tab.value = newTab;
    Cookies.set('tab', newTab);
    sidebarClose();
  }
</script>

<template>
  <mdui-navigation-drawer @close="sidebarClose()"
  :open="openStatus" close-on-esc close-on-overlay-click>
    <mdui-list>
      <mdui-list-item icon="task"
      @click="changeTab(0)">{{ pages[0] }}</mdui-list-item>
      <mdui-list-item icon="miscellaneous_services"
      @click="changeTab(1)">{{ pages[1] }}</mdui-list-item>
      <mdui-list-item icon="auto_fix_high"
      @click="changeTab(2)">{{ pages[2] }}</mdui-list-item>
      <mdui-list-item icon="dashboard"
      @click="changeTab(3)">{{ pages[3] }}</mdui-list-item>
    </mdui-list>
  </mdui-navigation-drawer>
</template>

<style scoped>

</style>