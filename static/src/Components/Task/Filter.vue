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

const task = ref([]);
const catagory = inject("catagory");
const by = inject("by");
const emitter = inject("emitter");

onMounted(async ()=>{
  const session = await Task("fetchTask", catagory.value, by.value);
  if (session) {
    task.value = session.task_data;
  }
})
</script>

<template>
  <div class="filter">
    <Region :tasks="task" @refresh="(emitter=4)"></Region>
    <Type :tasks="task" @refresh="(emitter=4)"></Type>
    <Status :tasks="task" @refresh="(emitter=4)"></Status>
    <API :tasks="task" @refresh="(emitter=4)"></API>
    <Model :tasks="task" @refresh="(emitter=4)"></Model>
    <Referer :tasks="task" @refresh="(emitter=4)"></Referer>
    <Time :tasks="task" @refresh="(emitter=4)"></Time>
    <Device :tasks="task" @refresh="(emitter=4)"></Device>
    <Info :tasks="task" @refresh="(emitter=4)"></Info>
  </div>
</template>

<style scoped>

</style>