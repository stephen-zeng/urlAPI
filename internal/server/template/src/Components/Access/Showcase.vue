<script setup>
import {ref, provide, inject, onMounted} from 'vue';
  import Detail from "@/Components/Access/Detail.vue";

  const dialogStatus = ref(false);
  const target = ref(0)
  const props = defineProps(["tasks"])

  provide('dialogStatus',dialogStatus);
  provide('target',target);

  function showDetail(newTarget) {
    target.value = newTarget;
    dialogStatus.value = true;
  }

  function getTime(tim) {
    var date = new Date(tim).toJSON();
    return new Date(+new Date(date)+8*3600*1000).toISOString().replace(/T/g,' ').replace(/\.[\d]{3}Z/,'');
  }
</script>

<template>
  <div class="showcase">
    <table class="mdui-table">
      <thead>
      <tr>
        <th>时间</th>
        <th>类型</th>
        <th>地区</th>
        <th>状态</th>
      </tr>
      </thead>
      <tbody>
        <tr @click="showDetail(item)" v-for="(item,index) in props.tasks">
          <td>{{ getTime(item.time) }}</td>
          <td>{{ item.type }}</td>
          <td>{{ item.region }}</td>
          <td>{{ item.status }}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <Detail></Detail>
</template>

<style scoped>
.showcase {
  width: 100%;
  max-height: 75%; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
td {
  vertical-align: middle;
  text-align: center;
}
</style>