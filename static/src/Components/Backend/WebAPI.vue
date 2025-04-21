<script setup>
import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()

async function getSetting() {
  settings.value = await Setting("fetchSetting", "web")
}

async function sendSetting() {
  await Setting("editSetting", "web", settings.value)
}
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="web" rounded @click="getSetting()">
        其他
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>Github App Token（选填）</p>
          <mdui-text-field variant="outlined" label="App Token" type="password" toggle-password
                           :value="settings?settings[0][5]:''"
                           @change="settings[0][5] = $event.target.value"></mdui-text-field>
          <p>YouTube API Token（必填）</p>
          <mdui-text-field variant="outlined" label="API Token" type="password" toggle-password
                           :value="settings?settings[0][6]:''"
                           @change="settings[0][6] = $event.target.value"></mdui-text-field>
          <mdui-divider></mdui-divider>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>