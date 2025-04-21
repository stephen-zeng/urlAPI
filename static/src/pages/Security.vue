<script setup>
import {onMounted, inject} from "vue";
import DashSecurity from "@/Components/Security/DashSecurity.vue";
import TaskBehavior from "@/Components/Security/TaskBehavior.vue";

const title = inject("title");

onMounted(()=>{
  title.value = "安全设置"
})
</script>

<template>
  <mdui-layout-main>
    <DashSecurity></DashSecurity>
    <TaskBehavior></TaskBehavior>
  </mdui-layout-main>
</template>

<style scoped>

</style>