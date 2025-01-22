<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="settings_applications" rounded>
        DeepSeek
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置DeepSeek的后端API key，可用于文字生成，总结等</p>
          <mdui-text-field variant="outlined" label="API Key"
                           clearable></mdui-text-field>
          <p style="margin-bottom: 0">默认文字生成模型</p>
          <mdui-radio-group value="chinese" style="margin-top: 0">
            <mdui-radio value="chinese">Chinese</mdui-radio>
            <mdui-radio value="english">English</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">默认总结模型</p>
          <mdui-radio-group value="chinese" style="margin-top: 0">
            <mdui-radio value="chinese">Chinese</mdui-radio>
            <mdui-radio value="english">English</mdui-radio>
          </mdui-radio-group>
          <mdui-button full-width>确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>