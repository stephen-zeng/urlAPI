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
          <p style="margin-bottom: 0">后台登录密码</p>
          <mdui-text-field type="password" variant="outlined"
                           toggle-password label="密码"></mdui-text-field>
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
          <mdui-divider></mdui-divider>
          <p style="margin-bottom: 0">能够使用的网站</p>
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

.list {
  width: 80%;
  max-height: 20rem; /* 设置最大高度 */
  overflow-y: auto;  /* 启用垂直滚动 */
}
</style>