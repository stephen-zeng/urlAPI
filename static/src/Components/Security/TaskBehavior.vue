<script setup>

import {ref} from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref({})
const input1 = ref('')
const input2 = ref('')

async function getSetting() {
  settings.value = await Setting("fetchSetting", "taskBehavior")
}

async function sendSetting() {
  await Setting("editSetting", "taskBehavior", settings.value)
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
          <p>不记录任务的Referer</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable
                           @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!=='') settings[0].push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[0]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[0].length>1) settings[0].splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <p>不记录任务的Info</p>
          <mdui-text-field variant="outlined" label="输入*表示这个部分遵从任何信息" clearable
                           @input="input2 = $event.target.value" :value="input2">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input2!=='') settings[1].push(input2);input2=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?settings[1]:[]" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings[1].length>1) settings[1].splice(index, 1)}"></mdui-button-icon>
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