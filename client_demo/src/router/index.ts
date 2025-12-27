import { createRouter, createWebHistory } from "vue-router";
import { auth } from "../config/firebase";
import { onAuthStateChanged } from "firebase/auth";
import LoginPage from "../pages/Login.vue";
import DashboardPage from "../pages/Dashboard.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "login",
      component: LoginPage,
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardPage,
      meta: { requiresAuth: true },
    },
  ],
});

// Get current user
const getCurrentUser = () => {
  return new Promise((resolve, reject) => {
    const removeListener = onAuthStateChanged(
      auth,
      (user) => {
        removeListener();
        resolve(user);
      },
      (error) => {
        removeListener();
        reject(error);
      },
    );
  });
};

// Navigation guard
router.beforeEach(async (to, _from, next) => {
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const user = await getCurrentUser();

  if (requiresAuth && !user) {
    // If the user is not logged in, redirect them to the login screen.
    next("/");
  } else if (to.path === "/" && user) {
    // If the user is already logged in,
    // attempting to access the login screen will redirect you to the dashboard screen.
    next("/dashboard");
  } else {
    next();
  }
});

export default router;
