import {createRouter, createWebHistory} from "vue-router";
import Task from "@/pages/Task.vue";
import Tool from "@/pages/Tool.vue";
import Backend from "@/pages/Backend.vue";
import Workshop from "@/pages/Workshop.vue";
import Login from "@/pages/Login.vue";

const routes = [
    {path:"/dash",component:Task},
    {path:"/dash/task",component:Task},
    {path:"/dash/tool",component:Tool},
    {path:"/dash/backend",component:Backend},
    {path:"/dash/workshop",component:Workshop},
    {path:"/dash/login",component:Login},
]

const router=  createRouter({
    history:createWebHistory(),
    routes,
})

export default router;