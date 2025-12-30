<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { auth } from "../config/firebase";
import {
  createUserWithEmailAndPassword,
  signInWithCustomToken,
} from "firebase/auth";
import { FirebaseError } from "firebase/app";
import { useAuthApi } from "../composables/useAuthApi";
import SignupForm from "../components/auth/SignupForm.vue";
import OTPRequestForm from "../components/auth/OTPRequestForm.vue";
import OTPVerifyForm from "../components/auth/OTPVerifyForm.vue";

const router = useRouter();
const { requestOTP, verifyOTP } = useAuthApi();

type AuthMode = "signup" | "login";
type LoginStep = "email" | "otp";

const authMode = ref<AuthMode>("login");
const loginStep = ref<LoginStep>("email");
const email = ref("");
const error = ref("");
const loading = ref(false);

/**
 * Handle signup with email and password
 */
const handleSignup = async (userEmail: string, password: string) => {
  error.value = "";
  loading.value = true;

  try {
    await createUserWithEmailAndPassword(auth, userEmail, password);
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
        default:
          error.value = "エラーが発生しました。もう一度お試しください";
          if (import.meta.env.DEV) {
            console.error("Signup error:", err);
          }
      }
    } else {
      error.value = "予期しないエラーが発生しました";
      if (import.meta.env.DEV) {
        console.error("Unexpected signup error:", err);
      }
    }
  } finally {
    loading.value = false;
  }
};

/**
 * Handle OTP request
 */
const handleOTPRequest = async (userEmail: string) => {
  error.value = "";
  loading.value = true;
  email.value = userEmail;

  try {
    const response = await requestOTP(userEmail);
    if (response) {
      loginStep.value = "otp";
      // Show OTP in development mode
      if (import.meta.env.DEV && response.otp) {
        console.log("🔐 OTP Code (Development):", response.otp);
      }
    } else {
      error.value = "OTPの送信に失敗しました";
    }
  } catch (err) {
    error.value = "予期しないエラーが発生しました";
    if (import.meta.env.DEV) {
      console.error("OTP request error:", err);
    }
  } finally {
    loading.value = false;
  }
};

/**
 * Handle OTP verification and login with custom token
 */
const handleOTPVerify = async (otp: string) => {
  error.value = "";
  loading.value = true;

  try {
    const customToken = await verifyOTP(email.value, otp);
    if (customToken) {
      // Sign in with custom token
      await signInWithCustomToken(auth, customToken);
      router.push("/dashboard");
    } else {
      error.value = "OTPが正しくありません";
    }
  } catch (err) {
    if (err instanceof FirebaseError) {
      switch (err.code) {
        case "auth/invalid-custom-token":
          error.value = "認証トークンが無効です";
          break;
        default:
          error.value = "ログインに失敗しました";
          if (import.meta.env.DEV) {
            console.error("Custom token error:", err);
          }
      }
    } else {
      error.value = "予期しないエラーが発生しました";
      if (import.meta.env.DEV) {
        console.error("Unexpected verify error:", err);
      }
    }
  } finally {
    loading.value = false;
  }
};

/**
 * Switch to login mode
 */
const switchToLogin = () => {
  authMode.value = "login";
  loginStep.value = "email";
  error.value = "";
};

/**
 * Switch to signup mode
 */
const switchToSignup = () => {
  authMode.value = "signup";
  error.value = "";
};

/**
 * Go back to email input step
 */
const backToEmailStep = () => {
  loginStep.value = "email";
  error.value = "";
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-base-200">
    <!-- Signup Form -->
    <SignupForm
      v-if="authMode === 'signup'"
      :loading="loading"
      :error="error"
      @submit="handleSignup"
      @switch-to-login="switchToLogin"
    />

    <!-- Login Flow -->
    <template v-else>
      <!-- Step 1: Email Input -->
      <OTPRequestForm
        v-if="loginStep === 'email'"
        :loading="loading"
        :error="error"
        @submit="handleOTPRequest"
        @switch-to-signup="switchToSignup"
      />

      <!-- Step 2: OTP Verification -->
      <OTPVerifyForm
        v-else
        :email="email"
        :loading="loading"
        :error="error"
        @submit="handleOTPVerify"
        @back="backToEmailStep"
      />
    </template>
  </div>
</template>
