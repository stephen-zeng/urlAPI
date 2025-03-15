<script setup>
import {ref, provide, inject, onMounted, watch} from 'vue';
import Detail from "@/Components/Task/Detail.vue";
import {Notification, Post} from "@/fetch.js";
import Cookies from "js-cookie";

const dialogStatus = ref(false);
const target = ref(0)
const task = ref([]);

const url = inject("url");
const emitter = inject("emitter");
const catagory = inject("catagory");
const by = inject("by");
const page = inject("page");
const title = inject("title");
const maxPage = inject("maxPage");

provide('dialogStatus',dialogStatus);
provide('target',target);


async function fetchTask() {
  if (by.value == "") {
    title.value = "任务查看"
  } else {
    title.value = "任务查看 → " + (by.value == "" ? "N/A" : by.value);
  }
  const session = await Post(url, {
    "Token": Cookies.get("token"),
    "Send": {
      "operation": "fetchTask",
      "task_catagory": catagory.value,
      "task_by": by.value,
      "task_page": page.value,
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    task.value = session.task_data
    maxPage.value = session.task_max_page;
  }
}

onMounted(async () => {
  await fetchTask();
})

function showDetail(newTarget) {
  target.value = newTarget;
  dialogStatus.value = true;
}

function getTime(tim) {
  var date = new Date(tim).toJSON();
  return new Date(+new Date(date)+8*3600*1000).toISOString().replace(/T/g,' ').replace(/\.[\d]{3}Z/,'');
}

watch(emitter, async(newVal, oldVal) => {
  switch (newVal) {
    case 1: //
      page.value = 1;
      title.value = "任务查看"
      by.value = "";
      catagory.value = "none";
      await fetchTask();
      break;
    case 2:
      if (page.value > 1) {
        page.value--;
        await fetchTask();
      }
      break;
    case 3:
      if (page.value < maxPage.value) {
        page.value++;
        await fetchTask();
      }
      break;
    case 4: // from filter
      page.value = 1;
      await fetchTask();
      break;
  }
  emitter.value = 0;
})
</script>

<template>
  <div class="showcase">
    <table class="mdui-table">
      <thead>
      <tr>
        <th>时间</th>
        <th>类型</th>
        <th>地区</th>
        <th>状态</th>
      </tr>
      </thead>
      <tbody>
        <tr @click="showDetail(item)" v-for="(item,index) in task">
          <td>{{ getTime(item.time) }}</td>
          <td>{{ item.type }}</td>
          <td>{{ item.region }}</td>
          <td>{{ item.status }}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <Detail></Detail>
</template>

<style scoped>
.showcase {
  width: 100%;
  max-height: 75%; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
td {
  vertical-align: middle;
  text-align: center;
}
</style>