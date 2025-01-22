<script setup>
  import {defineProps, inject} from 'vue'
  import Theme from "@/Components/Theme.vue";

  const props = defineProps(['title'])
  const sidebarStatus = inject('sidebarStatus')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
</script>

<template>
  <mdui-top-app-bar>
    <mdui-button-icon icon="menu"
    @click="SidebarStatusChanged()"></mdui-button-icon>
    <mdui-top-app-bar-title>{{ props.title }}</mdui-top-app-bar-title>
    <div style="flex-grow: 1"></div>
    <Theme></Theme>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>