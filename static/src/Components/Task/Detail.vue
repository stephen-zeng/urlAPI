<script setup>
  import {ref, inject} from 'vue'
  import {Notification} from "@/js/util.js";

  const status = inject('dialogStatus');
  const target = inject('target');
  const copy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      Notification('Copied')
    } catch {
      Notification('Copy failed')
    }
  }
</script>

<template>
  <mdui-dialog
  :open="status" @close="status = false"
  close-on-overlay-click close-on-esc
  headline="Detail" description="The detailed information of this task">
    <mdui-list>
      <mdui-list-item v-for="(value, key) in target"
                      @click="copy(value)" class="copy" :data-clipboard-text="value">
        {{ key }} - {{ value }}
      </mdui-list-item>
    </mdui-list>
  </mdui-dialog>
</template>

<style scoped>

</style>