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
      "part": "web"
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    settings.value = session.setting
    console.log(settings.value)
  }
}

async function sendSetting() {
  const session = await Post(url+"session", {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "edit",
      "part": "web",
      "edit": settings.value,
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    Notification("Saved")
  }
}

function find(list, status, value, operation) {
  let index
  if (operation == "edit" && status == true) {
    list.push(value)
  }
  for (let i = 0; i < list.length; i++) {
    if (list[i] == value) {
      index = i
      if (operation == "find") {
        return true
      }
    }
  }
  if (operation == "find") {
    return false
  }
  if (operation == "edit" && status == false) {
    list.splice(index, 1)
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
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="web" rounded @click="getSetting">
        网页
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">网页总结开关</p>
          <mdui-radio-group :value="settings?settings[0][0]:'false'"
                            @change="settings[0][0]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">网页总结使用的API</p>
          <mdui-radio-group :value="settings?settings[0][2]:'openai'"
                            @change="settings[0][2]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="otherapi">其他API</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">网页缩略图开关</p>
          <mdui-radio-group :value="settings?settings[0][1]:'false'"
                            @change="settings[0][1]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">允许生成缩略图的网站</p>
          <div class="mdui-checkbox-group">
<!--            <mdui-checkbox>Github（需要网络支持）</mdui-checkbox>-->
<!--            <mdui-checkbox>YouTube（需要网络支持）</mdui-checkbox>-->
<!--            <mdui-checkbox>Gitee</mdui-checkbox>-->
<!--            <mdui-checkbox>网易云</mdui-checkbox>-->
<!--            <mdui-checkbox>Bilibili</mdui-checkbox>-->
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'github.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'github.com', 'edit')">
              Github（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'www.youtube.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'www.youtube.com', 'edit')">
              YouTube（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'gitee.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'gitee.com', 'edit')">
              Gitee</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'music.163.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'music.163.com', 'edit')">
              网易云</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'www.bilibili.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'www.bilibili.com', 'edit')">
              B站</mdui-checkbox>
          </div>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">不能使用“总结”功能的网站（黑名单）</p>
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
          <mdui-button full-width @click="sendSetting">确认</mdui-button>
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