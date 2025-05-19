<script setup>
import {onMounted, inject} from "vue";
import DashSecurity from "@/Components/Security/DashSecurity.vue";
import TaskBehavior from "@/Components/Security/TaskBehavior.vue";
import ImgSecurity from "@/Components/Security/ImgSecurity.vue";
import TxtSecurity from "@/Components/Security/TxtSecurity.vue";

const title = inject("title");

onMounted(()=>{
  title.value = "安全设置"
})
</script>

<template>
  <mdui-layout-main>
    <DashSecurity></DashSecurity>
    <TaskBehavior></TaskBehavior>
    <ImgSecurity></ImgSecurity>
    <TxtSecurity></TxtSecurity>
  </mdui-layout-main>
</template>

<style scoped>

</style>