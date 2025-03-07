<script setup>
  import { inject } from 'vue';
  import {useRouter} from "vue-router";

  const openStatus = inject('sidebarStatus');
  const tab = inject('tab');
  const tabAddition = inject('tabAddition');
  const pages = inject('pages');
  const router = useRouter();

  function sidebarClose() {
    openStatus.value = false;
  }
  function changeTab(newTab) {
    router.push(newTab);
    sidebarClose();
  }
</script>

<template>
  <mdui-navigation-drawer @close="sidebarClose()"
  :open="openStatus" close-on-esc close-on-overlay-click>
    <mdui-list>
      <mdui-list-item icon="task"
      @click="changeTab('/task')">{{ pages[0] }}</mdui-list-item>
      <mdui-list-item icon="miscellaneous_services"
      @click="changeTab('/backend')">{{ pages[1] }}</mdui-list-item>
      <mdui-list-item icon="auto_fix_high"
      @click="changeTab('/tool')">{{ pages[2] }}</mdui-list-item>
      <mdui-list-item icon="dashboard"
      @click="changeTab('/workshop')">{{ pages[3] }}</mdui-list-item>
    </mdui-list>
  </mdui-navigation-drawer>
</template>

<style scoped>

</style>