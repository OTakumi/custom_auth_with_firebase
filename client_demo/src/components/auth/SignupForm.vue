<script setup lang="ts">
import { ref } from "vue";
import FormInput from "./FormInput.vue";
import ErrorMessage from "./ErrorMessage.vue";
import SubmitButton from "./SubmitButton.vue";

interface Props {
  loading: boolean;
  error: string;
}

defineProps<Props>();

const emit = defineEmits<{
  submit: [email: string, password: string];
  switchToLogin: [];
}>();

const email = ref("");
const password = ref("");

const handleSubmit = () => {
  emit("submit", email.value, password.value);
};
</script>

<template>
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-3xl font-bold justify-center mb-4">
        新規登録
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <FormInput
          id="signup-email"
          v-model="email"
          label="メールアドレス"
          type="email"
          placeholder="example@email.com"
          autocomplete="email"
          :required="true"
        />

        <FormInput
          id="signup-password"
          v-model="password"
          label="パスワード"
          type="password"
          placeholder="6文字以上"
          autocomplete="new-password"
          :minlength="6"
          :required="true"
        />

        <ErrorMessage :message="error" />

        <SubmitButton :loading="loading" :is-sign-up="true" />
      </form>

      <div class="divider"></div>

      <div class="text-center">
        <p class="text-sm text-base-content/70">
          すでにアカウントをお持ちですか？
        </p>
        <button
          type="button"
          class="btn btn-link btn-sm"
          @click="emit('switchToLogin')"
        >
          ログインはこちら
        </button>
      </div>
    </div>
  </div>
</template>
