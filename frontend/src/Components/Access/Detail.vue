<script setup>
  import {ref, inject} from 'vue'

  const status = inject('dialogStatus');
  const target = inject('target');
</script>

<template>
  <mdui-dialog
  :open="status" @close="status = false"
  close-on-overlay-click close-on-esc
  headline="Detail" description="The detailed information of this access">
    <h1>{{ target }}</h1>
  </mdui-dialog>
</template>

<style scoped>

</style>