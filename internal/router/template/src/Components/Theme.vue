<script setup>
  import {ref, onMounted} from 'vue'
  import {setColorScheme} from "mdui";
  import Cookies from "js-cookie";

  const theme = ref(0)
  const color = ref(0)

  function getTheme() {
    if (theme.value % 3 === 0) { // auto
      return "hdr_auto"
    } else if (theme.value % 3 === 1) {
      return "light_mode"
    } else if (theme.value % 3 === 2) {
      return "dark_mode"
    }
  }
  function setTheme() {
    theme.value ++
    if (theme.value % 3 === 0) { // dark -> auto
      document.documentElement.classList.remove("mdui-theme-dark");
      document.documentElement.classList.add("mdui-theme-auto");
      Cookies.set("theme", "mdui-theme-auto", {expires: 365});
    } else if (theme.value % 3 === 1) { // auto -> light
      document.documentElement.classList.remove("mdui-theme-auto");
      document.documentElement.classList.add("mdui-theme-light");
      Cookies.set("theme", "mdui-theme-light", {expires: 365});
    } else if (theme.value % 3 === 2) { // light -> dark
      document.documentElement.classList.remove("mdui-theme-light");
      document.documentElement.classList.add("mdui-theme-dark");
      Cookies.set("theme", "mdui-theme-dark", {expires: 365});
    }
  }
  function changeColor() {
    color.value ++;
    if (color.value % 4 === 0) setColorScheme('#00ff77');
    else if (color.value % 4 === 1) setColorScheme('#00aaff');
    else if (color.value % 4 === 2) setColorScheme('#eeff00');
    else if (color.value % 4 === 3) setColorScheme('#ff0000');
  }

  onMounted(() => {
    if (Cookies.get("theme")) document.documentElement.classList.add(Cookies.get("theme"));
    else document.documentElement.classList.add("mdui-theme-auto");
    setColorScheme('#00ff77');
  })
</script>

<template>
  <mdui-button-icon @click="changeColor()"
  icon="color_lens"></mdui-button-icon>
  <mdui-button-icon
  :icon="getTheme()" @click="setTheme()">
  </mdui-button-icon>
</template>

<style scoped>

</style>