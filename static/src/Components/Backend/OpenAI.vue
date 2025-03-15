<script setup>

import { ref, inject } from 'vue';
import { Post, Notification } from "@/fetch.js"
import Cookies from "js-cookie";

const url = inject('url')
const settings = ref()

async function getSetting() {
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "openai"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting_data
  }
}

async function sendSetting() {
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "editSetting",
      "setting_part": "openai",
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
        OpenAI
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置OpenAI的后端API属性，可用于文字生成，总结等</p>
          <mdui-text-field variant="outlined" label="API Key" type="password" toggle-password
                           :value="settings?settings[0][0]:''"
                           @change="settings[0][0] = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="API地址"
                           :value="settings?settings[0][5]:'https://api.openai.com/v1/chat/completions'"
                           @change="settings[0][5] = $event.target.value"></mdui-text-field>
          <p>默认文字生成模型</p>
          <mdui-radio-group :value="settings?settings[0][1]:'gpt-4o'" style="margin-top: 0"
                            @change="settings[0][1]=$event.target.value">
            <mdui-radio value="gpt-4o">GPT-4o</mdui-radio>
            <mdui-radio value="gpt-4o-mini">GPT-4o-mini</mdui-radio>
          </mdui-radio-group>
          <p>默认总结模型</p>
          <mdui-radio-group :value="settings?settings[0][2]:'gpt-4o'" style="margin-top: 0"
                            @change="settings[0][2]=$event.target.value">
            <mdui-radio value="gpt-4o">GPT-4o</mdui-radio>
            <mdui-radio value="gpt-4o-mini">GPT-4o-mini</mdui-radio>
          </mdui-radio-group>
          <p>默认图片生成模型</p>
          <mdui-radio-group :value="settings?settings[0][3]:'dall-e-3'" style="margin-top: 0"
                            @change="settings[0][3]=$event.target.value">
            <mdui-radio value="dall-e-3">DALL·E 3</mdui-radio>
            <mdui-radio value="dall-e-2">DALL·E 2</mdui-radio>
          </mdui-radio-group>
          <p>默认图片生成大小</p>
          <mdui-radio-group :value="settings?settings[0][4]:'1024x1024'" style="margin-top: 0"
                            v-if="settings?settings[0][3]=='dall-e-3':false"
                            @change="settings[0][4]=$event.target.value">
            <mdui-radio value="1024x1024">正方形 (1024x1024)</mdui-radio>
            <mdui-radio value="1792x1024">横屏（1792x1024）</mdui-radio>
            <mdui-radio value="1024x1792">竖屏 (1024x1792)</mdui-radio>
          </mdui-radio-group>
          <mdui-radio-group :value="settings?settings[0][4]:'1024x1024'" style="margin-top: 0"
                            v-if="settings?settings[0][3]=='dall-e-2':false"
                            @change="settings[0][4]=$event.target.value">
            <mdui-radio value="256x256">小 (256x256)</mdui-radio>
            <mdui-radio value="512x512">中（512x512）</mdui-radio>
            <mdui-radio value="1024x1024">大 (1024x1024)</mdui-radio>
          </mdui-radio-group>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>