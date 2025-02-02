<script setup>
import {inject, ref} from 'vue';
import {snackbar} from "mdui";
import { sha256 } from "js-sha256";
import Cookies from 'js-cookie';

const dialogStatus = ref(true);
const pwd = ref("")
const term = ref(false)
const loginStatus = inject("login");
const url = inject("url");

function login() {
  fetch(url + "session", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": sha256(pwd.value),
    },
    body: JSON.stringify({
      operation: "login",
      term: term.value,
    })
  }).then((response) => response.json())
      .then((data) => {
        if (data.error) {
          snackbar({
            message: data.error,
            placement: "top-end",
          })
        } else {
          Cookies.set("token", data.token)
          loginStatus.value = true;
          snackbar({
            message: "Login success",
            placement: "top-end",
          })
        }
      })
}
</script>

<template>
  <mdui-layout-main>
    <mdui-card variant="outlined">
      <mdui-icon name='login'></mdui-icon>
      <h1>登录后台</h1>
      <mdui-text-field type="password"
                       toggle-password label="密码"
                       :value="pwd" @input="pwd = $event.target.value"></mdui-text-field>
      <mdui-checkbox :checked="term" @input="term = !$event.target.checked">7天内保持登录</mdui-checkbox>
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
