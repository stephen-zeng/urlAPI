<script setup>

import { ref, inject } from 'vue';
import { Post, Notification } from "@/fetch.js"
import Cookies from "js-cookie";

const url = inject('url')
const settings = ref()
const input2 = ref('')

async function getSetting() {
  const session = await Post({
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchSetting",
      "setting_part": "web"
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
      "setting_part": "web",
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
          <p>网页缩略图开关</p>
          <mdui-radio-group :value="settings?settings[0][1]:'false'"
                            @change="settings[0][1]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p>允许生成缩略图的网站</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'github.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'github.com', 'edit')">
              Github（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'gitee.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'gitee.com', 'edit')">
              Gitee</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'www.youtube.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'www.youtube.com', 'edit')">
              Youtube（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'www.bilibili.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'www.bilibili.com', 'edit')">
              B站</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'arxiv.org', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'arxiv.org', 'edit')">
              Arxiv</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?settings[1]:[], false, 'www.ithome.com', 'find')"
                           @change="find(settings?settings[1]:[], $event.target.checked, 'www.ithome.com', 'edit')">
              IT之家</mdui-checkbox>
          </div>
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