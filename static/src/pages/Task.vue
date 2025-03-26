<script setup>
  import Region from "@/Components/Task/Region.vue";
  import Type from "@/Components/Task/Type.vue"
  import Showcase from "@/Components/Task/Showcase.vue";
  import API from "@/Components/Task/API.vue";
  import Model from "@/Components/Task/Model.vue";
  import {ref, provide, onMounted, inject, watch} from 'vue';
  import Cookies from "js-cookie";
  import { Post, Notification } from "@/fetch.js"
  import Status from "@/Components/Task/Status.vue";
  import Referer from "@/Components/Task/Referer.vue";
  import Time from "@/Components/Task/Time.vue";
  import Device from "@/Components/Task/Device.vue";
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