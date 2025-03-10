<script setup>
  import {ref, inject} from 'vue'
  import {snackbar} from "mdui";
  import Clipboard from "clipboard"

  const status = inject('dialogStatus');
  const target = inject('target');
  const catagory = inject('catagory');
  const by = inject('by');

  function notification(message) {
    snackbar({
      message: message,
      placement: "top-end",
    })
  }
  const copy = () => {
    let clipboard = new Clipboard('.copy')
    clipboard.on('success', (e) => {
      notification('Copied')
      clipboard.destroy()
    })
    clipboard.on('error', (e) => {
      notification('Error copying')
      clipboard.destroy()
    })
  }
</script>

<template>
  <mdui-dialog
  :open="status" @close="status = false"
  close-on-overlay-click close-on-esc
  headline="Detail" description="The detailed information of this task">
    <mdui-list>
      <mdui-list-item v-for="(value, key) in target"
                      @click="copy" class="copy" :data-clipboard-text="value">
        {{ key }} - {{ value }}
      </mdui-list-item>
    </mdui-list>
  </mdui-dialog>
</template>

<style scoped>

</style>