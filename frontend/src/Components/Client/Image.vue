<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="image" rounded>
        图像
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">生成</p>
          <mdui-radio-group value="on" style="margin-top: 0">
            <mdui-radio value="on">开启</mdui-radio>
            <mdui-radio value="off">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">图像生成使用的API</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">超分</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="chinese">开启</mdui-radio>
            <mdui-radio value="english">关闭</mdui-radio>
          </mdui-radio-group>
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