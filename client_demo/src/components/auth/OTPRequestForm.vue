<script setup lang="ts">
import { ref } from "vue";
import FormInput from "./FormInput.vue";
import ErrorMessage from "./ErrorMessage.vue";

interface Props {
  loading: boolean;
  error: string;
}

defineProps<Props>();

const emit = defineEmits<{
  submit: [email: string];
  switchToSignup: [];
}>();

const email = ref("");

const handleSubmit = () => {
  emit("submit", email.value);
};
</script>

<template>
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-3xl font-bold justify-center mb-4">
        ログイン
      </h2>

      <p class="text-sm text-base-content/70 mb-4">
        メールアドレスを入力してください。ワンタイムパスワードを送信します。
      </p>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <FormInput
          id="login-email"
          v-model="email"
          label="メールアドレス"
          type="email"
          placeholder="example@email.com"
          autocomplete="email"
          :required="true"
        />

        <ErrorMessage :message="error" />

        <button
          type="submit"
          class="btn btn-primary w-full"
          :disabled="loading"
        >
          <span v-if="loading" class="loading loading-spinner"></span>
          {{ loading ? "送信中..." : "ワンタイムパスワードを送信" }}
        </button>
      </form>

      <div class="divider"></div>

      <div class="text-center">
        <p class="text-sm text-base-content/70">
          アカウントをお持ちでない方
        </p>
        <button
          type="button"
          class="btn btn-link btn-sm"
          @click="emit('switchToSignup')"
        >
          新規登録はこちら
        </button>
      </div>
    </div>
  </div>
</template>
