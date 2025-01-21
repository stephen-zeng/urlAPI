<script setup>
  import Region from "@/Components/Access/Region.vue";
  import Type from "@/Components/Access/Type.vue"
  import Showcase from "@/Components/Access/Showcase.vue";
  import {ref, provide} from 'vue';

  const catalog = ref("all");

  provide("catalog", catalog);
</script>

<template>
  <mdui-layout-main>
    <Showcase></Showcase>
    <Region></Region>
    <Type></Type>
  </mdui-layout-main>
</template>

<style scoped>

</style>