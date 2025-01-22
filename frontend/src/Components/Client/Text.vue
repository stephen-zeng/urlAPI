<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="translate" rounded>
        文字
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">随机生成的启用情况</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox>随机笑话</mdui-checkbox>
            <mdui-checkbox>随机名言</mdui-checkbox>
            <mdui-checkbox>随机段子</mdui-checkbox>
            <mdui-checkbox>自定义提示词</mdui-checkbox>
          </div>
          <p style="margin-bottom: 0">随机生成使用的API</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">网页总结</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="chinese">开启</mdui-radio>
            <mdui-radio value="english">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">网页总结使用的API</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">文字总结的启用情况</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox>Word</mdui-checkbox>
            <mdui-checkbox>PDF</mdui-checkbox>
          </div>
          <p style="margin-bottom: 0">文字总结使用的API</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
            <mdui-radio value="deepseek">DeepSeek</mdui-radio>
          </mdui-radio-group>
          <mdui-button full-width>确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>