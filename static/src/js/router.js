import {createRouter, createWebHistory} from "vue-router";
import Task from "@/pages/Task.vue";
import Tool from "@/pages/Tool.vue";
import Backend from "@/pages/Backend.vue";
import Workshop from "@/pages/Workshop.vue";
import Login from "@/pages/Login.vue";
import Security from "@/pages/Security.vue";

const routes = [
    { path: "/dash", component: () => import('@/pages/Task.vue') },
    { path: "/dash/task", component: () => import('@/pages/Task.vue') },
    { path: "/dash/tool", component: () => import('@/pages/Tool.vue') },
    { path: "/dash/security", component: () => import('@/pages/Security.vue') },
    { path: "/dash/backend", component: () => import('@/pages/Backend.vue') },
    { path: "/dash/workshop", component: () => import('@/pages/Workshop.vue') },
    {path:"/dash/login",component:Login},
]


const router=  createRouter({
    history:createWebHistory(),
    routes,
})

export default router;