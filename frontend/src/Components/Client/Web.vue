<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="web" rounded>
        网页
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p style="margin-bottom: 0">网页总结开关</p>
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
          <p style="margin-bottom: 0">网页缩略图开关</p>
          <mdui-radio-group value="openai" style="margin-top: 0">
            <mdui-radio value="chinese">开启</mdui-radio>
            <mdui-radio value="english">关闭</mdui-radio>
          </mdui-radio-group>
          <p style="margin-bottom: 0">允许生成缩略图的网站</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox>Github（需要网络支持）</mdui-checkbox>
            <mdui-checkbox>YouTube（需要网络支持）</mdui-checkbox>
            <mdui-checkbox>知乎</mdui-checkbox>
            <mdui-checkbox>CSDN</mdui-checkbox>
            <mdui-checkbox>Gitee</mdui-checkbox>
            <mdui-checkbox>Gitlab</mdui-checkbox>
            <mdui-checkbox>网易云</mdui-checkbox>
            <mdui-checkbox>Bilibili</mdui-checkbox>
            <mdui-button full-width>确认</mdui-button>
          </div>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">不能使用“总结”功能的网站（黑名单）</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable>
            <mdui-button-icon slot="end-icon" icon="add"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item>*</mdui-list-item>
            </mdui-list>
          </div>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>
.list {
  width: 80%;
  max-height: 20rem; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
</style>