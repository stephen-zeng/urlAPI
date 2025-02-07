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
        "setting_part": "alibaba"
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
        "setting_part": "alibaba",
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
        阿里巴巴
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置阿里巴巴的后端API，可用于文字生成，总结等</p>
          <mdui-text-field variant="outlined" label="API Key" type="password" toggle-password
                           :value="settings?settings[0][0]:''"
                           @change="settings[0][0] = $event.target.value"></mdui-text-field>
          <p>默认文字生成模型</p>
          <mdui-radio-group :value="settings?settings[0][1]:'qwen-plus'"
                            @change="settings[0][1]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="qwen-max">通义千问-Max</mdui-radio>
            <mdui-radio value="qwen-plus">通义千问-Plus</mdui-radio>
            <mdui-radio value="qwen-turbo">通义千问-Turbo</mdui-radio>
            <mdui-radio value="qwen-long">通义千问-Long</mdui-radio>
          </mdui-radio-group>
          <p>默认总结模型</p>
          <mdui-radio-group :value="settings?settings[0][2]:'qwen-plus'"
                            @change="settings[0][2]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="qwen-max">通义千问 Max</mdui-radio>
            <mdui-radio value="qwen-plus">通义千问 Plus</mdui-radio>
            <mdui-radio value="qwen-turbo">通义千问 Turbo</mdui-radio>
            <mdui-radio value="qwen-long">通义千问 Long</mdui-radio>
          </mdui-radio-group>
          <p>默认图片生成模型</p>
          <mdui-radio-group :value="settings?settings[0][3]:'wanx2.1-t2i-turbo'"
                            @change="settings[0][3]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="wanx2.1-t2i-plus">通义万相文生图 2.1 Plus</mdui-radio>
            <mdui-radio value="wanx2.1-t2i-turbo">通义万相文生图 2.1 Turbo</mdui-radio>
            <mdui-radio value="wanx2.0-t2i-turbo">通义万相文生图 2.0</mdui-radio>
            <mdui-radio value="wanx-v1">通义万相文生图 1.0</mdui-radio>
          </mdui-radio-group>
          <p>默认图片生成大小</p>
          <mdui-radio-group :value="settings?settings[0][4]:'1024*1024'"
                            @change="settings[0][4]=$event.target.value"
                            v-if="settings?settings[0][3]=='wanx-v1':false"
                            style="margin-top: 0">
            <mdui-radio value="1024*1024">正方形 (1024x1024)</mdui-radio>
            <mdui-radio value="1280*720">横屏 (1280x720)</mdui-radio>
            <mdui-radio value="720*1280">竖屏 (720x1280)</mdui-radio>
            <mdui-radio value="768*1152">竖屏 (768x1152)</mdui-radio>
          </mdui-radio-group>
          <mdui-radio-group :value="settings?settings[0][4]:'1024*1024'"
                            @change="settings[0][4]=$event.target.value"
                            v-if="settings?settings[0][3]!='wanx-v1':false"
                            style="margin-top: 0">
            <mdui-radio value="1024*1024">大 (1024x1024)</mdui-radio>
            <mdui-radio value="512*512">中 (512x512)</mdui-radio>
            <mdui-radio value="1024*768">横屏 (1024x768)</mdui-radio>
            <mdui-radio value="768*1024">竖屏 (768x1024)</mdui-radio>
          </mdui-radio-group>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>