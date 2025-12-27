<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { auth } from "../config/firebase";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
} from "firebase/auth";
import { FirebaseError } from "firebase/app";
import AuthForm from "../components/auth/AuthForm.vue";

const router = useRouter();

const isSignUp = ref(false);
const error = ref("");
const loading = ref(false);

const handleSubmit = async (email: string, password: string) => {
  error.value = "";
  loading.value = true;

  try {
    if (isSignUp.value) {
      await createUserWithEmailAndPassword(auth, email, password);
    } else {
      await signInWithEmailAndPassword(auth, email, password);
    }
    router.push("/dashboard");
  } catch (err) {
    if (err instanceof FirebaseError) {
      switch (err.code) {
        case "auth/email-already-in-use":
          error.value = "このメールアドレスは既に使用されています";
          break;
        case "auth/invalid-email":
          error.value = "メールアドレスの形式が正しくありません";
          break;
        case "auth/weak-password":
          error.value = "パスワードは6文字以上で入力してください";
          break;
        case "auth/user-not-found":
          error.value = "ユーザーが見つかりません";
          break;
        case "auth/wrong-password":
          error.value = "パスワードが正しくありません";
          break;
        case "auth/invalid-credential":
          error.value = "メールアドレスまたはパスワードが正しくありません";
          break;
        default:
          error.value = "エラーが発生しました。もう一度お試しください";
          if (import.meta.env.DEV) {
            console.error("Authentication error:", err);
          }
      }
    } else {
      error.value = "予期しないエラーが発生しました";
      if (import.meta.env.DEV) {
        console.error("Unexpected error:", err);
      }
    }
  } finally {
    loading.value = false;
  }
};

const toggleMode = () => {
  isSignUp.value = !isSignUp.value;
  error.value = "";
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-base-200">
    <AuthForm
      :is-sign-up="isSignUp"
      :loading="loading"
      :error="error"
      @submit="handleSubmit"
      @toggle-mode="toggleMode"
    />
  </div>
</template>
