<script setup>
  import {inject} from "vue";

  const catalog = inject("catalog");
</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="location_on" rounded>
        地区
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <div style="margin-left: 2.5rem">
        <mdui-list-item rounded slot="custom">
          China - 123456次
        </mdui-list-item>
        <mdui-list-item rounded>US - 4567次</mdui-list-item>
      </div>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>