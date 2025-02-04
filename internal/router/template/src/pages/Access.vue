<script setup>
  import Region from "@/Components/Access/Region.vue";
  import Type from "@/Components/Access/Type.vue"
  import Showcase from "@/Components/Access/Showcase.vue";
  import {ref, provide, onMounted, inject} from 'vue';
  import Cookies from "js-cookie";
  import {snackbar} from "mdui";
  import { Post, Notification } from "@/fetch.js"
  import Status from "@/Components/Access/Status.vue";

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
        "operation": "task",
        "by": catagory.value,
        "task": by.value,
      }
    })
    if (session.error) {
      Notification(session.error)
    } else {
      task.value = session.task
    }
  }

  onMounted(async() => {
    await fetchTask()
  })
</script>

<template>
  <mdui-layout-main style="display: block">
    <Showcase :tasks="task" style="width: 100%" @refresh="fetchTask"></Showcase>
    <Region :tasks="task" style="width: 100%" @refresh="fetchTask"></Region>
    <Type :tasks="task" style="width: 100%" @refresh="fetchTask"></Type>
    <Status :tasks="task" style="width: 100%" @refresh="fetchTask"></Status>
  </mdui-layout-main>
</template>

<style scoped>

</style>