<script setup>
  import {inject} from 'vue'
  import Theme from "@/frameworks/Theme.vue";
  import Cookies from "js-cookie";
  import {useRoute, useRouter} from "vue-router";
  import {Logout} from "@/js/util.js";

  const title = inject("title");
  const sidebarStatus = inject('sidebarStatus')
  const login = inject('login')
  const emitter = inject('emitter')
  const router = useRouter();
  const route = useRoute();
  const page = inject('page')
  const maxPage = inject('maxPage')

  function SidebarStatusChanged() {
    sidebarStatus.value = !sidebarStatus.value;
  }
  async function logout() {
    if (await Logout(Cookies.get("token"))) {
      Cookies.remove("token");
      router.push("/dash/login");
    }
  }
</script>

<template>
  <mdui-top-app-bar>
    <mdui-button-icon icon="menu"
    @click="SidebarStatusChanged()"></mdui-button-icon>
    <mdui-top-app-bar-title>{{ title }}</mdui-top-app-bar-title>
    <div style="flex-grow: 1"></div>
<!--    1 for refresh, 2 for backwards, 3 for forwards-->
    <mdui-segmented-button-group v-if="login && route.path === '/dash/task'">
      <mdui-segmented-button @click="(emitter=2)">←</mdui-segmented-button>
      <mdui-segmented-button>{{ page }} / {{ maxPage }}</mdui-segmented-button>
      <mdui-segmented-button @click="(emitter=3)">→</mdui-segmented-button>
    </mdui-segmented-button-group>
    <mdui-button-icon @click="(emitter=1)" v-if="login && route.path === '/dash/task'" icon="refresh"></mdui-button-icon>

    <Theme v-if="route.path !== '/dash/task'"></Theme>
    <mdui-button-icon @click="logout()" v-if="login"
                      icon="exit_to_app"></mdui-button-icon>
  </mdui-top-app-bar>
</template>

<style scoped>

</style>