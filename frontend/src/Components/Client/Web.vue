<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="web" rounded>
        网页缩略图
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">允许的网站</p>
          <mdui-checkbox>Github（需要网络支持）</mdui-checkbox>
          <mdui-checkbox>YouTube（需要网络支持）</mdui-checkbox>
          <mdui-checkbox>知乎</mdui-checkbox>
          <mdui-checkbox>CSDN</mdui-checkbox>
          <mdui-checkbox>Gitee</mdui-checkbox>
          <mdui-checkbox>Gitlab</mdui-checkbox>
          <mdui-checkbox>网易云</mdui-checkbox>
          <mdui-checkbox>Bilibili</mdui-checkbox>
          <mdui-button full-width>确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>
mdui-list-item mdui-card {
  width: 100%;
  height: 50%;
}
mdui-text-field {
  width: 96%;
  margin-left: 2%;
}
</style>