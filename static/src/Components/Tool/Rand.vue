<script setup>

import { ref } from 'vue';
import {Repo, Setting} from "@/js/util.js";

const settings = ref()
const repoAPI = ref("github")
const repoInfo = ref("")
const repos = ref([])

async function getRepos() {
  repos.value = await Repo("fetchRepo")
}
async function editRepo(operation, id) {
  await Repo(operation, id)
  await getRepos()
}
async function newRepo() {
  if (repoInfo.value === "") return
  await Repo("newRepo", "", repoAPI.value, repoInfo.value)
  repoInfo.value = ""
  await getRepos()
}

async function getSetting() {
  settings.value = await Setting("fetchSetting", "rand")
  await getRepos()
}
async function sendSetting() {
  await Setting("editSetting", "rand", settings.value)
}
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item>
      <mdui-list-item slot="header" icon="all_inclusive" rounded @click="getSetting">
        随机图片
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
          <p>默认API</p>
          <mdui-radio-group :value="settings?settings[0][3]:'github'"
                            @change="settings[0][3]=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="github">Github</mdui-radio>
            <mdui-radio value="gitee">Gitee</mdui-radio>
          </mdui-radio-group>
          <p>添加可用仓库</p>
          <mdui-radio-group :value="repoAPI"
                            @change="repoAPI=$event.target.value"
                            style="margin-top: 0; margin-bottom: 0">
            <mdui-radio value="github">Github</mdui-radio>
            <mdui-radio value="gitee">Gitee</mdui-radio>
          </mdui-radio-group>
          <mdui-text-field variant="outlined" label="用户名/仓库名" :value="repoInfo"
                           clearable @input="repoInfo = $event.target.value">
            <mdui-button-icon slot="end-icon" icon="add" @click="newRepo"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="item in repos" nonclickable>
                <mdui-button-icon slot="icon" icon="refresh" @click="editRepo('refreshRepo', item.uuid)"></mdui-button-icon>
                {{ item.api }} - {{ item.info }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="editRepo('delRepo', item.uuid)"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-divider></mdui-divider>
          <p>Github随机图片中RAW网址</p>
          <mdui-text-field variant="outlined" label="反代地址"
                           :value="settings?settings[0][1]:''"
                           @change="settings[0][1] = $event.target.value"></mdui-text-field>
          <p>生成失败时返回的图片</p>
          <mdui-text-field variant="outlined" label="URL"
                           :value="settings?settings[0][2]:''"
                           @change="settings[0][2] = $event.target.value"></mdui-text-field>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>