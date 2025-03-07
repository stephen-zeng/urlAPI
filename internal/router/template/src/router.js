import {createRouter, createWebHistory} from "vue-router";
import Task from "@/pages/Task.vue";
import Tool from "@/pages/Tool.vue";
import Backend from "@/pages/Backend.vue";
import Workshop from "@/pages/Workshop.vue";
import Login from "@/pages/Login.vue";

const routes = [
    {path:"/:pathMatch(.*)*",component:Task},
    {path:"/task",component:Task},
    {path:"/tool",component:Tool},
    {path:"/backend",component:Backend},
    {path:"/workshop",component:Workshop},
    {path:"/login",component:Login},
]

const router=  createRouter({
    history:createWebHistory(),
    routes,
})

export default router;