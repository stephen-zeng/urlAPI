<script setup>
import {inject, ref} from "vue";

const catalog = inject("catalog");
const props = defineProps(["tasks"])
const map = ref({})

function getValue(tasks) {
  if (Object.keys(map.value).length) return
  for (let task of tasks) {
    if (task.type in map.value) {
      map.value[task.type]++;
    } else {
      map.value[task.type] = 1;
    }
  }
}
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded @click="getValue(props.tasks)">
      <mdui-list-item slot="header" icon="assessment" rounded>
        种类
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <div style="margin-left: 2.5rem">
        <mdui-list-item rounded slot="custom" v-for="(value, key) in map">
          {{ key }} - {{ value }}次
        </mdui-list-item>
      </div>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>