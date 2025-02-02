<script setup>
  import Region from "@/Components/Access/Region.vue";
  import Type from "@/Components/Access/Type.vue"
  import Showcase from "@/Components/Access/Showcase.vue";
  import {ref, provide, onMounted, inject} from 'vue';
  import Cookies from "js-cookie";
  import {snackbar} from "mdui";

  const catalog = ref("all");
  const url = inject("url");
  const task = ref([]);

  provide("catalog", catalog);

  onMounted(() => {
    fetch(url+"session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": Cookies.get("token"),
      },
      body: JSON.stringify({
        operation: "task",
      })
    }).then(res => res.json()).then((data) => {
      if (data.error) {
        snackbar({
          message: data.error,
          placement: "top-end",
        })
      } else {
        task.value = data.task
      }
    })
  })
</script>

<template>
  <mdui-layout-main>
    <Showcase :tasks="task"></Showcase>
    <Region :tasks="task"></Region>
    <Type :tasks="task"></Type>
  </mdui-layout-main>
</template>

<style scoped>

</style>