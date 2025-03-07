<script setup>
  import Text from "@/Components/Client/Text.vue";
  import Image from "@/Components/Client/Image.vue";
  import Web from "@/Components/Client/Web.vue";
  import Security from "@/Components/Client/Security.vue";
  import Rand from "@/Components/Client/Rand.vue"
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