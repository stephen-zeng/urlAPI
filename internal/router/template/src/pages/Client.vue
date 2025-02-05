<script setup>
  import Text from "@/Components/Client/Text.vue";
  import Image from "@/Components/Client/Image.vue";
  import Web from "@/Components/Client/Web.vue";
  import Security from "@/Components/Client/Security.vue";
  import Rand from "@/Components/Client/Rand.vue"
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