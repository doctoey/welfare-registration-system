import { createRouter, createWebHistory } from "vue-router";
import Register from "../views/Register.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", redirect: "/register" },
    { path: "/register", component: Register },
  ],
});

export default router;
