<script setup>
import {inject, onMounted} from "vue";
const url = inject("url")
const type = inject("type")
</script>

<template>
 <mdui-card variant="outlined" class="workshop-result-card">
   <img v-if="type=='img'" :src="url">
   <iframe v-if="type=='txt'" :src="url"></iframe>
 </mdui-card>
</template>
<style scoped>
img, iframe {
  height: 100%;
  width: 100%;
}
</style>