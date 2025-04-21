<script setup>

import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()

async function getSetting() {
  settings.value = await Setting("fetchSetting", "contxt")
}

async function sendSetting() {
  await Setting("editSetting", "contxt", settings.value)
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="texture" rounded @click="getSetting">
        提示词相关
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置提示词语境及具体的提示词</p>
          <mdui-text-field variant="outlined" label="生成的语境"
                           :value="settings?settings[0][0]:''"
                           @change="settings[0][0] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="总结的语境"
                           :value="settings?settings[0][1]:''"
                           @change="settings[0][1] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成笑话的提示词 (laugh)"
                           :value="settings?settings[1][0]:''"
                           @change="settings[1][0] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成诗句的提示词 (poem)"
                           :value="settings?settings[1][1]:''"
                           @change="settings[1][1] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成鸡汤的提示词 (sentence)"
                           :value="settings?settings[1][2]:''"
                           @change="settings[1][2] = $event.target.value"></mdui-text-field>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>