<script setup>
  import Region from "@/Components/Access/Region.vue";
  import Type from "@/Components/Access/Type.vue"
  import Showcase from "@/Components/Access/Showcase.vue";
  import API from "@/Components/Access/API.vue";
  import Model from "@/Components/Access/Model.vue";
  import {ref, provide, onMounted, inject, watch} from 'vue';
  import Cookies from "js-cookie";
  import { Post, Notification } from "@/fetch.js"
  import Status from "@/Components/Access/Status.vue";
  import Referer from "@/Components/Access/Referer.vue";
  import Time from "@/Components/Access/Time.vue";
  import Device from "@/Components/Access/Device.vue";

  const catagory = ref("");
  const by = ref("")
  const url = inject("url");
  const task = ref([]);
  const emitter = inject("emitter");
  const title = inject("title");

  provide("catagory", catagory);
  provide("by", by);

  async function fetchTask(filter = 1) {
    title.value = "任务查看";
    if (filter === 1) {
      title.value += " → " + (by.value == "" ? "N/A" : by.value);
    }
    const session = await Post(url, {
      "Token": Cookies.get("token"),
      "Send": {
        "operation": "fetchTask",
        "task_catagory": catagory.value,
        "task_by": by.value,
      }
    })
    if (session.error) {
      Notification(session.error)
    } else {
      task.value = session.task_data.sort((a, b) => new Date(b.time) - new Date(a.time));
    }
  }

  onMounted(async() => {
    await fetchTask(0);
  })

  watch(emitter, async(newVal, oldVal) => {
    title.value = "任务查看";
    if (newVal == 1) {
      by.value = "";
      catagory.value = "";
      await fetchTask(0);
      emitter.value = 0; // reset the emitter
    }
  })
</script>

<template>
  <mdui-layout-main style="display: block">
    <Showcase :tasks="task" style="width: 100%" @refresh="fetchTask(0)"></Showcase>
    <div class="filter">
      <Region :tasks="task" @refresh="fetchTask"></Region>
      <Type :tasks="task" @refresh="fetchTask"></Type>
      <Status :tasks="task" @refresh="fetchTask"></Status>
      <API :tasks="task" @refresh="fetchTask"></API>
      <Model :tasks="task" @refresh="fetchTask"></Model>
      <Referer :tasks="task" @refresh="fetchTask"></Referer>
      <Time :tasks="task" @refresh="fetchTask"></Time>
      <Device :tasks="task" @refresh="fetchTask"></Device>
    </div>
  </mdui-layout-main>
</template>

<style scoped>
</style>