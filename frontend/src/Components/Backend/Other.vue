<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="settings_applications" rounded>
        其他兼容API
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置其他兼容的后端API key，可用于文字生成，总结等，为OpenAI的格式</p>
          <mdui-text-field variant="outlined" label="API地址"
                           clearable></mdui-text-field>
          <mdui-text-field variant="outlined" label="API Key"
                           clearable></mdui-text-field>
          <mdui-text-field variant="outlined" label="默认文字生成模型"
                           clearable></mdui-text-field>
          <mdui-text-field variant="outlined" label="默认总结模型"
                           clearable></mdui-text-field>
          <mdui-button full-width>确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>