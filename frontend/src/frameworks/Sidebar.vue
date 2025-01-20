<script setup>
  import { inject } from 'vue';

  const openStatus = inject('sidebarStatus');
  const tab = inject('tab');
  const pages = inject('pages');

  function sidebarClose() {
    openStatus.value = false;
  }
</script>

<template>
  <mdui-navigation-drawer @close="sidebarClose()"
  :open="openStatus" close-on-esc close-on-overlay-click>
    <mdui-list>
      <mdui-list-item icon="task"
      @click="tab=0">{{ pages[0] }}</mdui-list-item>
      <mdui-list-item icon="miscellaneous_services"
      @click="tab=1">{{ pages[1] }}</mdui-list-item>
      <mdui-list-item icon="auto_fix_high"
      @click="tab=2">{{ pages[2] }}</mdui-list-item>
    </mdui-list>
  </mdui-navigation-drawer>
</template>

<style scoped>

</style>