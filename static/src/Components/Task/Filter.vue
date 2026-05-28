<script setup>

import Type from "@/Components/Task/Type.vue";
import Model from "@/Components/Task/Model.vue";
import Time from "@/Components/Task/Time.vue";
import Status from "@/Components/Task/Status.vue";
import API from "@/Components/Task/API.vue";
import Region from "@/Components/Task/Region.vue";
import Device from "@/Components/Task/Device.vue";
import Referer from "@/Components/Task/Referer.vue";
import {inject, onMounted, ref} from "vue";
import Info from "@/Components/Task/Info.vue";
import {Task} from "@/js/util.js";
import Temp from "@/Components/Task/Temp.vue";

const task = ref([]);
const catagory = inject("catagory");
const by = inject("by");
const emitter = inject("emitter");
const fetched = ref(false);

onMounted(async ()=>{
  fetched.value = false;
  const session = await Task("fetchTaskStats", catagory.value, by.value);
  if (session) {
    task.value = session.task_stats;
    fetched.value = true;
  }
})
</script>

<template>
  <div class="filter">
    <Region :tasks="task.region" @refresh="(emitter=4)" :fetched="fetched"></Region>
    <Type :tasks="task.type" @refresh="(emitter=4)" :fetched="fetched"></Type>
    <Status :tasks="task.status" @refresh="(emitter=4)" :fetched="fetched"></Status>
    <API :tasks="task.api" @refresh="(emitter=4)" :fetched="fetched"></API>
    <Model :tasks="task.model" @refresh="(emitter=4)" :fetched="fetched"></Model>
    <Referer :tasks="task.referer" @refresh="(emitter=4)" :fetched="fetched"></Referer>
    <Time :tasks="task.time" @refresh="(emitter=4)" :fetched="fetched"></Time>
    <Device :tasks="task.device" @refresh="(emitter=4)" :fetched="fetched"></Device>
    <Info :tasks="task.more_info" @refresh="(emitter=4)" :fetched="fetched"></Info>
    <Temp :tasks="task.temp" @refresh="(emitter=4)" :fetched="fetched"></Temp>
  </div>
</template>

<style scoped>

</style>
