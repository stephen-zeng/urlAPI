<script setup>

import { ref, inject } from 'vue';
import { Post, Notification } from "@/fetch.js"
import Cookies from "js-cookie";
import {sha256} from "js-sha256";

const url = inject('url')
const settings = ref({})
const input1 = ref('')
const input2 = ref('')
const ip = ref('')

async function getSetting() {
  const session = await Post(url, {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "security"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting_data
    ip.value = session.session_ip
  }
}

async function sendSetting() {
  const session = await Post(url, {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "editSetting",
      "setting_part": "security",
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
    <mdui-collapse-item>
      <mdui-list-item slot="header" icon="security" rounded @click="getSetting">
        安全
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>后台登录密码</p>
          <mdui-text-field type="password" variant="outlined"
                           @change="settings[0][0] = sha256($event.target.value)"
                           toggle-password label="密码"></mdui-text-field>
          <p>允许登录后台的IP</p>
          <mdui-text-field variant="outlined" label="输入*为该子段都可以使用" clearable
                           @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!='') settings[1].push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[1]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[1].length>1) settings[1].splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="settings[1].push(ip)">添加本机IP</mdui-button>
          <mdui-divider></mdui-divider>
          <p>可以使用urlAPI的网站（防盗）</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable
                           @input="input2 = $event.target.value" :value="input2">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input2!='') settings[2].push(input2);input2=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[2]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[2].length>1) settings[2].splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>