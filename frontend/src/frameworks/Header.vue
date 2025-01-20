<script setup>
  import {defineProps, inject} from 'vue'

  const props = defineProps(['title'])
  const sidebarStatus = inject('sidebarStatus')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
</script>

<template>
  <mdui-top-app-bar
  scroll-behavior="hide"
  scroll-target="content">
    <mdui-button-icon icon="menu"
    @click="SidebarStatusChanged()"></mdui-button-icon>
    <mdui-top-app-bar-title>{{ props.title }}</mdui-top-app-bar-title>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>