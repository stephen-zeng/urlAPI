<script setup>

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item>
      <mdui-list-item slot="header" icon="security" rounded>
        安全
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <mdui-text-field variant="outlined" label="访问秘钥"
                           type="password" toggle-password></mdui-text-field>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">外链白名单，单击对应项目以删除</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable>
            <mdui-button-icon slot="end-icon" icon="add"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item>*</mdui-list-item>
            </mdui-list>
          </div>
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">外链黑名单，单击对应项目以删除</p>
          <mdui-text-field variant="outlined" label="输入*为该子域都可以使用" clearable>
            <mdui-button-icon slot="end-icon" icon="add"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item>*</mdui-list-item>
            </mdui-list>
          </div>
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
.list {
  width: 100%;
  max-height: 20rem; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
</style>