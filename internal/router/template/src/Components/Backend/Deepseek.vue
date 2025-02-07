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
      "setting_part": "deepseek"
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
      "setting_part": "deepseek",
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
      <mdui-list-item slot="header" icon="settings_applications" rounded @click="getSetting()">
        DeepSeek
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置DeepSeek的后端API，可用于文字生成，总结等</p>
          <mdui-text-field variant="outlined" label="API Key" type="password" toggle-password
                           :value="settings?settings[0][0]:''"
                           @change="settings[0][0] = $event.target.value"></mdui-text-field>
          <p>默认文字生成模型</p>
          <mdui-radio-group :value="settings?settings[0][1]:'deepseek-chat'" style="margin-top: 0"
                            @change="settings[0][1]=$event.target.value">
            <mdui-radio value="deepseek-chat">DeepSeek V3</mdui-radio>
            <mdui-radio value="deepseek-reasoner">DeepSeek R1</mdui-radio>
          </mdui-radio-group>
          <p>默认总结模型</p>
          <mdui-radio-group :value="settings?settings[0][2]:'deepseek-chat'" style="margin-top: 0"
                            @change="settings[0][2]=$event.target.value">
            <mdui-radio value="deepseek-chat">DeepSeek V3</mdui-radio>
            <mdui-radio value="deepseek-reasoner">DeepSeek R1</mdui-radio>
          </mdui-radio-group>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>