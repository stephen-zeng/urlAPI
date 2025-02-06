<script setup>
import {inject, ref} from "vue";

const props = defineProps(["tasks"])
const emits = defineEmits(["refresh"])
const map = ref({})
const catagory = inject("catagory");
const by = inject("by");

function getValue(tasks) {
  if (Object.keys(map.value).length) {
    map.value = {}
  }
  for (let task of tasks) {
    if (task.api in map.value) {
      map.value[task.api]++;
    } else {
      map.value[task.api] = 1;
    }
  }
}
function setFilter(filter) {
  catagory.value = "api";
  by.value = filter;
  emits("refresh")
}
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="settings_applications" rounded @click="getValue(props.tasks)">
        API
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <div style="margin-left: 2.5rem">
        <mdui-list-item rounded slot="custom"
                        v-for="(value, key) in map"
                        @click="setFilter(key)">
          {{ key === "" ? "No API" : key }} - {{ value }}次
        </mdui-list-item>
      </div>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>