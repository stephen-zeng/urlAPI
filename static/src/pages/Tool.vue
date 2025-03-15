<script setup>
  import Text from "@/Components/Tool/Text.vue";
  import Image from "@/Components/Tool/Image.vue";
  import Web from "@/Components/Tool/Web.vue";
  import Security from "@/Components/Tool/Security.vue";
  import Rand from "@/Components/Tool/Rand.vue"
  import {onMounted, inject} from "vue";

  const title = inject("title");

  onMounted(()=>{
    title.value = "功能设置"
  })
</script>

<template>
  <mdui-layout-main>
    <Security></Security>
    <Text></Text>
    <Image></Image>
    <Rand></Rand>
    <Web></Web>
  </mdui-layout-main>
</template>

<style scoped>

</style>