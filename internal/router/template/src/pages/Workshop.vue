<script setup>
import Control from "@/Components/Workshop/Control.vue";
import Result from "@/Components/Workshop/Result.vue";
import { provide, ref, onMounted, inject} from "vue";

const title = inject("title");

onMounted(()=>{
  title.value = "工作台"
})
const url = ref("");
const type = ref("");
const show = ref(false);
provide("url", url)
provide("type", type)
provide("show", show)
</script>

<template>
  <mdui-layout-main>
    <Control></Control>
    <Result v-if="show"></Result>
  </mdui-layout-main>
</template>

<style scoped>
mdui-card {
  margin: 1rem;
  height: auto;
}
@media (min-width: 1500px) {
  .workshop-control-card {
    width: calc(100% / 4 - 2rem)
  }
  .workshop-result-card {
    width: calc(100% / 4 * 3 - 2rem)
  }
}
@media (min-width: 1000px) and (max-width: 1500px) {
  .workshop-control-card {
    width: calc(100% / 3 - 2rem)
  }
  .workshop-result-card {
    width: calc(100% / 3 * 2 - 2rem)
  }
}
@media (min-width: 500px) and (max-width: 1000px) {
  .workshop-control-card, .workshop-result-card {
    width: calc(100% / 2 - 2rem)
  }
}
@media (max-width: 500px) {
  .workshop-control-card, .workshop-result-card {
    width: calc(100% - 2rem)
  }
}
</style>