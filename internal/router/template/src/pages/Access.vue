<script setup>
  import Region from "@/Components/Access/Region.vue";
  import Type from "@/Components/Access/Type.vue"
  import Showcase from "@/Components/Access/Showcase.vue";
  import API from "@/Components/Access/API.vue";
  import Model from "@/Components/Access/Model.vue";
  import {ref, provide, onMounted, inject} from 'vue';
  import Cookies from "js-cookie";
  import {snackbar} from "mdui";
  import { Post, Notification } from "@/fetch.js"
  import Status from "@/Components/Access/Status.vue";
  import Referer from "@/Components/Access/Referer.vue";
  import Time from "@/Components/Access/Time.vue";

  const catagory = ref("");
  const by = ref("")
  const url = inject("url");
  const task = ref([]);

  provide("catagory", catagory);
  provide("by", by);

  async function fetchTask() {
    const session = await Post(url + "session", {
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
    await fetchTask()
  })
</script>

<template>
  <mdui-layout-main style="display: block">
    <Showcase :tasks="task" style="width: 100%" @refresh="fetchTask"></Showcase>
    <div class="filter">
      <Region :tasks="task" @refresh="fetchTask"></Region>
      <Type :tasks="task" @refresh="fetchTask"></Type>
      <Status :tasks="task" @refresh="fetchTask"></Status>
      <API :tasks="task" @refresh="fetchTask"></API>
      <Model :tasks="task" @refresh="fetchTask"></Model>
      <Referer :tasks="task" @refresh="fetchTask"></Referer>
      <Time :tasks="task" @refresh="fetchTask"></Time>
    </div>
  </mdui-layout-main>
</template>

<style scoped>
</style>