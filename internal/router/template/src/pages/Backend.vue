<script setup>

import OpenAI from "@/Components/Backend/OpenAI.vue";
import Alibaba from "@/Components/Backend/Alibaba.vue";
import Other from "@/Components/Backend/Other.vue";
import Deepseek from "@/Components/Backend/Deepseek.vue";
import Context from "@/Components/Backend/Context.vue";
</script>

<template>
  <mdui-layout-main>
    <OpenAI></OpenAI>
    <Alibaba></Alibaba>
    <Deepseek></Deepseek>
    <Other></Other>
    <Context></Context>
  </mdui-layout-main>
</template>

<style scoped>

</style>