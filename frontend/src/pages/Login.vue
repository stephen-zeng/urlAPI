<script setup>
import {inject, ref} from 'vue';
import {snackbar} from "mdui";

const dialogStatus = ref(true);

function login() {
  snackbar({
    message: "Login",
    placement: "top-end",
  })
}
</script>

<template>
  <mdui-layout-main>
    <mdui-card variant="outlined">
      <mdui-icon name='login'></mdui-icon>
      <h1>登录后台</h1>
      <mdui-text-field type="password"
                       toggle-password label="密码"></mdui-text-field>
      <mdui-button @click="login()">登录</mdui-button>
    </mdui-card>
  </mdui-layout-main>
</template>

<style scoped>
mdui-layout-main {
  display: flex;
  justify-content: center;
  align-items: center;
}

mdui-layout-main mdui-card {
  height: 50%;
  width: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

mdui-text-field {
  width: 80%;
}
</style>
