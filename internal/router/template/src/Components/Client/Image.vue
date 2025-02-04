<script setup>

import { ref, inject } from 'vue';
import { Post, Notification } from "@/fetch.js"
import Cookies from "js-cookie";
import {sha256} from "js-sha256";

const url = inject('url')
const settings = ref()
const input1 = ref('')
const input2 = ref('')
const ip = ref('')

async function getSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetch",
      "part": "img"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting
  }
}

async function sendSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "edit",
      "part": "img",
      "edit": settings.value,
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
      <mdui-list-item slot="header" icon="image" rounded @click="getSetting">
        图像
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">总开关</p>
          <mdui-radio-group :value="settings?settings[0][0]:'false'"
                            @change="settings[0][0]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">图像生成使用的API</p>
          <mdui-radio-group :value="settings?settings[0][1]:'openai'"
                            @change="settings[0][1]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">Github随机图片中RAW网址</p>
          <mdui-text-field variant="outlined" label="反代地址"
                           :value="settings?settings[0][2]:''"
                           @change="settings[0][2] = $event.target.value"></mdui-text-field>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">过期时间</p>
          <mdui-text-field variant="outlined" label="分钟"
                           :value="settings?settings[0][3]:'60'"
                           @change="settings[0][3] = $event.target.value"></mdui-text-field>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">生成失败时返回的图片</p>
          <mdui-text-field variant="outlined" label="URL"
                           :value="settings?settings[0][4]:'60'"
                           @change="settings[0][4] = $event.target.value"></mdui-text-field>
          <mdui-button full-width @click="sendSetting">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>