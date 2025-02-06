<script setup>
import {inject, onMounted, ref} from 'vue';
import {snackbar} from "mdui";
import { sha256 } from "js-sha256";
import { Post, Notification } from "@/fetch.js";
import Cookies from 'js-cookie';

const pwd = ref("")
const term = ref(false)
const loginStatus = inject("login");
const url = inject("url");

async function login() {
  const session = await Post(url + "session", {
    "Token": sha256(pwd.value),
    "Send": {
      "operation": "login",
      "login_term": term.value,
    }
  })
  if (session.error) {
    Notification(session.error)
  } else {
    Cookies.set("token", session.session_token);
    loginStatus.value = true;
    Notification("Login successful!");
  }
}
</script>

<template>
  <mdui-layout-main>
    <mdui-card variant="outlined">
      <mdui-icon name='login'></mdui-icon>
      <h1>登录后台</h1>
      <mdui-text-field type="password" @keydown.enter="login()"
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
  height: 25rem;
  width: 35rem;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

mdui-text-field {
  width: 80%;
}
</style>
