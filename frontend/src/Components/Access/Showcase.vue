<script setup>
  import {ref, provide, inject} from 'vue';
  import Detail from "@/Components/Access/Detail.vue";

  const dialogStatus = ref(false);
  const target = ref(0)
  const catalog = inject('catalog');

  provide('dialogStatus',dialogStatus);
  provide('target',target);

  function showDetail(newTarget) {
    target.value = newTarget;
    dialogStatus.value = true;
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
      <tr @click="showDetail(1)">
        <td>1</td>
        <td>Mark</td>
        <td>Otto</td>
        <td>@mdo</td>
      </tr>
      <tr @click="showDetail(2)">
        <td>2</td>
        <td>Jacob</td>
        <td>Thornton</td>
        <td>@fat</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>Otto</td>
        <td>@mdo</td>
      </tr>
      <tr>
        <td>2</td>
        <td>Jacob</td>
        <td>Thornton</td>
        <td>@fat</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>Otto</td>
        <td>@mdo</td>
      </tr>
      <tr>
        <td>2</td>
        <td>Jacob</td>
        <td>Thornton</td>
        <td>@fat</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr><tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
      </tr>
      <tr>
        <td>1</td>
        <td>Mark</td>
        <td>asdf</td>
        <td>哈哈</td>
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