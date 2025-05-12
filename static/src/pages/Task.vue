<script setup>
  import Showcase from "@/Components/Task/Showcase.vue";
  import {ref, provide} from 'vue';
  import Filter from "@/Components/Task/Filter.vue";

  const catagory = ref("none");
  const by = ref("");

  provide("catagory", catagory);
  provide("by", by);

</script>

<template>
  <mdui-layout-main style="display: block">
    <Showcase></Showcase>
    <Filter></Filter>
  </mdui-layout-main>
</template>

<style scoped>
</style>