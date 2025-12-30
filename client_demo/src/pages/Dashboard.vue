<script setup lang="ts">
import { computed, ref } from "vue";
import { auth } from "../config/firebase";
import { signOut } from "firebase/auth";
import { useRouter } from "vue-router";
import { FirebaseError } from "firebase/app";

const router = useRouter();

const user = computed(() => auth.currentUser);
const userEmail = computed(() => {
  if (!user.value?.email) {
    console.error("User email not found. Authentication state is invalid.");
    throw new Error("User email not found");
  }
  return user.value.email;
});
const logoutError = ref("");

const handleLogout = async () => {
  logoutError.value = "";
  try {
    await signOut(auth);
    router.push("/");
  } catch (err) {
    if (err instanceof FirebaseError) {
      logoutError.value = "ログアウトに失敗しました。もう一度お試しください。";
      if (import.meta.env.DEV) {
        console.error("Logout error:", err);
      }
    } else {
      logoutError.value = "予期しないエラーが発生しました。";
      if (import.meta.env.DEV) {
        console.error("Unexpected logout error:", err);
      }
    }
  }
};
</script>

<template>
  <div class="min-h-screen bg-base-200">
    <div class="navbar bg-base-100 shadow-lg">
      <div class="flex-1">
        <a class="btn btn-ghost text-xl">Dashboard</a>
      </div>
      <div class="flex-none gap-2">
        <button class="btn" @click="handleLogout">ログアウト</button>
      </div>
    </div>

    <div class="container mx-auto p-4 md:p-8">
      <div v-if="logoutError" class="alert alert-error mb-4">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6 shrink-0 stroke-current"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <span>{{ logoutError }}</span>
      </div>

      <div class="card bg-base-100 shadow-xl mb-8">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div>
              <h2 class="card-title text-3xl">ようこそ！</h2>
              <p class="text-base-content/70">{{ userEmail }}</p>
            </div>
          </div>
          <p class="mt-4">
            ログインに成功しました。このダッシュボードは、認証されたユーザーのみがアクセスできる保護されたページです。
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
