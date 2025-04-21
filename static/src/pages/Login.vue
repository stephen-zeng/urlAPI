<script setup>
import {inject, onMounted, ref} from 'vue';
import { sha256 } from "js-sha256";
import { Post, Notification } from "@/js/fetch.js";
import Cookies from 'js-cookie';
import {useRouter} from "vue-router";
import {Login} from "@/js/util.js";

const pwd = ref("")
const term = ref(false)
const loginStatus = inject("login");
const title = inject("title");
const router = useRouter();

async function login() {
  if (await Login(sha256(pwd.value), term.value)) {
    loginStatus.value = true;
    router.push("/dash/task");
  }
}

onMounted(async() => {
  title.value = "登录";
  if (Cookies.get("token") && await (Cookies.get("token"), false)) {
      loginStatus.value = true;
      router.push("/dash/task");
  }
})
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
mdui-button {
  width: 80%;
}

mdui-text-field {
  width: 80%;
}
</style>
