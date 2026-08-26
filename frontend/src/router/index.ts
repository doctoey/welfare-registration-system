import { createRouter, createWebHistory } from "vue-router";
import Register from "../views/Register.vue";
import Status from "../views/Status.vue";
import Officer from "../views/Officer.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", redirect: "/register" },
    { path: "/register", component: Register },
    { path: "/status", component: Status },
    { path: "/officer", component: Officer },
  ],
});

export default router;
