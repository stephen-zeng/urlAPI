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
      "part": "security"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting
    ip.value = session.ip
  }
}

async function sendSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "edit",
      "part": "security",
      "edit": settings.value,
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    Notification("Saved")
  }
}
function del(list, index) {
  if (list.length > 1) {
    list.splice(index, 1)
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
          <p style="margin-bottom: 0">后台登录密码</p>
          <mdui-text-field type="password" variant="outlined"
                           @change="settings[0][0] = sha256($event.target.value)"
                           toggle-password label="密码"></mdui-text-field>
          <p style="margin-bottom: 0">允许登录后台的IP</p>
          <mdui-text-field variant="outlined" label="输入*为该子段都可以使用" clearable @input="input1 = $event.target.value">
            <mdui-button-icon slot="end-icon" icon="add" @click="settings[1].push(input1)"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[1]:[]" @click="del(settings[1], index)">
                {{ item }}
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="settings[1].push(ip)">添加本机IP</mdui-button>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">可以使用urlAPI的网站（防盗）</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable @input="input2 = $event.target.value">
            <mdui-button-icon slot="end-icon" icon="add" @click="settings[2].push(input2)"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[2]:[]" @click="del(settings[2], index)">
                {{ item }}
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button full-width @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

.list {
  width: 80%;
  max-height: 20rem; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
</style>