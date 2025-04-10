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
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "txt"
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
      "setting_part": "txt",
      "setting_edit": settings.value,
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

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="translate" rounded @click="getSetting">
        文字
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>总开关</p>
          <mdui-radio-group :value="settings?settings[0][0]:'false'"
                            @change="settings[0][0]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p>随机生成的启用情况</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'laugh', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'laugh', 'edit')">
              随机笑话</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'poem', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'poem', 'edit')">
              随机诗句</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'sentence', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'sentence', 'edit')">
              随机鸡汤</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'other', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'other', 'edit')">
              其他提示词</mdui-checkbox>
          </div>
          <p>随机生成使用的API</p>
          <mdui-radio-group :value="settings?settings[0][1]:'openai'"
                            @change="settings[0][1]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="otherapi">其他API</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p>文字总结使用的API</p>
          <mdui-radio-group :value="settings?settings[0][2]:'openai'"
                            @change="settings[0][2]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="otherapi">其他API</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p>允许的Prompt</p>
          <mdui-text-field variant="outlined" label="适配通配符" clearable
                           @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!='') settings[2].push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[2]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[2].length>1) settings[2].splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-divider></mdui-divider>
          <p>过期时间</p>
          <mdui-text-field variant="outlined" label="分钟"
                           :value="settings?settings[0][3]:'60'"
                           @change="settings[0][3] = $event.target.value"></mdui-text-field>
          <p>生成失败时返回的图片</p>
          <mdui-text-field variant="outlined" label="URL"
                           :value="settings?settings[0][4]:''"
                           @change="settings[0][4] = $event.target.value"></mdui-text-field>
          <mdui-button @click="sendSetting">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>