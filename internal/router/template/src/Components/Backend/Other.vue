<script setup>

import { ref, inject } from 'vue';
import { Post, Notification } from "@/fetch.js"
import Cookies from "js-cookie";

const url = inject('url')
const settings = ref()

async function getSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "otherapi"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting_data
  }
}

async function sendSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "editSetting",
      "setting_part": "otherapi",
      "setting_edit": settings.value,
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    Notification("Saved")
  }
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="settings_applications" rounded @click="getSetting">
        其他兼容API
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置其他兼容的后端API key，可用于文字生成，总结等，为OpenAI的格式</p>
          <mdui-text-field variant="outlined" label="API Key" type="password" toggle-password
                           :value="settings?settings[0][0]:''"
                           @change="settings[0][0] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="API地址"
                           :value="settings?settings[0][3]:''"
                           @change="settings[0][3] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="文字生成模型"
                           :value="settings?settings[0][2]:''"
                           @change="settings[0][2] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="文字总结模型"
                           :value="settings?settings[0][1]:''"
                           @change="settings[0][1] = $event.target.value"></mdui-text-field>
          <mdui-button full-width @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>