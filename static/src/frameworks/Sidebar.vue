<script setup>
  import { inject } from 'vue';
  import {useRouter} from "vue-router";
  import Cookies from "js-cookie";

  const openStatus = inject('sidebarStatus');
  const pages = inject('pages');
  const router = useRouter();

  function sidebarClose() {
    openStatus.value = false;
  }
  function changeTab(newTab) {
    router.push(newTab);
    Cookies.set("tab", newTab);
    sidebarClose();
  }
</script>

<template>
  <mdui-navigation-drawer @close="sidebarClose()"
  :open="openStatus" close-on-esc close-on-overlay-click>
    <mdui-list>
      <mdui-list-item icon="task"
      @click="changeTab('/dash/task')">{{ pages[0] }}</mdui-list-item>
      <mdui-list-item icon="miscellaneous_services"
      @click="changeTab('/dash/backend')">{{ pages[1] }}</mdui-list-item>
      <mdui-list-item icon="auto_fix_high"
      @click="changeTab('/dash/tool')">{{ pages[2] }}</mdui-list-item>
      <mdui-list-item icon="security"
      @click="changeTab('/dash/security')">{{ pages[3] }}</mdui-list-item>
      <mdui-list-item icon="dashboard"
      @click="changeTab('/dash/workshop')">{{ pages[4] }}</mdui-list-item>
    </mdui-list>
  </mdui-navigation-drawer>
</template>

<style scoped>

</style>