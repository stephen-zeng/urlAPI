<script setup>
import {inject, ref} from "vue";

const props = defineProps(["tasks", "fetched"])
const emits = defineEmits(["refresh"])
const map = ref({})
const catagory = inject("catagory");
const by = inject("by");

function getValue(tasks) {
  if (Object.keys(map.value).length) {
    map.value = {}
  }
  for (let task of tasks) {
    if (task.temp in map.value) {
      map.value[task.temp]++;
    } else {
      map.value[task.temp] = 1;
    }
  }
  const sortedEntries = Object.entries(map.value).sort((a, b) => b[1] - a[1]);
  map.value = Object.fromEntries(sortedEntries);
}
function setFilter(filter) {
  catagory.value = "temp";
  by.value = (filter === "" ? "N/A" : filter);
  emits("refresh")
}
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="file_download_done" rounded @click="getValue(props.tasks)" :disabled="!props.fetched">
        缓存
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <div style="margin-left: 2.5rem">
        <mdui-list-item rounded slot="custom"
                        v-for="(value, key) in map"
                        @click="setFilter(key)">
          {{ key === "" ? "N/A" : key }} - {{ value }}次
        </mdui-list-item>
      </div>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>