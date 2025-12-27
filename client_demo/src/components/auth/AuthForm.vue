<script setup lang="ts">
import { ref } from "vue";
import FormInput from "./FormInput.vue";
import ErrorMessage from "./ErrorMessage.vue";
import SubmitButton from "./SubmitButton.vue";
import AuthModeToggle from "./AuthModeToggle.vue";

interface Props {
  isSignUp: boolean;
  loading: boolean;
  error: string;
}

defineProps<Props>();

const emit = defineEmits<{
  submit: [email: string, password: string];
  toggleMode: [];
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
        {{ isSignUp ? "新規登録" : "ログイン" }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <FormInput
          id="email"
          v-model="email"
          label="メールアドレス"
          type="email"
          placeholder="example@email.com"
          autocomplete="email"
          :required="true"
        />

        <FormInput
          id="password"
          v-model="password"
          label="パスワード"
          type="password"
          placeholder="6文字以上"
          autocomplete="current-password"
          :minlength="6"
          :required="true"
        />

        <ErrorMessage :message="error" />

        <SubmitButton :loading="loading" :is-sign-up="isSignUp" />
      </form>

      <AuthModeToggle :is-sign-up="isSignUp" @toggle="emit('toggleMode')" />
    </div>
  </div>
</template>
