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
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "taskBehavior",
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
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "editSetting",
      "setting_part": "taskBehavior",
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
      <mdui-list-item slot="header" icon="task" rounded @click="getSetting">
        任务记录行为
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>任务不记录的Referer</p>
          <mdui-text-field variant="outlined" label="输入*为该子段都可以使用" clearable
                           @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!='') settings[0].push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[0]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[0].length>1) settings[0].splice(index, 1)}"></mdui-button-icon>
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