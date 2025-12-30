<script setup lang="ts">
import { ref } from "vue";
import ErrorMessage from "./ErrorMessage.vue";

interface Props {
  email: string;
  loading: boolean;
  error: string;
}

defineProps<Props>();

const emit = defineEmits<{
  submit: [otp: string];
  back: [];
}>();

const otp = ref("");

const handleSubmit = () => {
  emit("submit", otp.value);
};
</script>

<template>
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-3xl font-bold justify-center mb-4">
        ワンタイムパスワード入力
      </h2>

      <p class="text-sm text-base-content/70 mb-2">
        <strong>{{ email }}</strong> に送信されたワンタイムパスワードを入力してください。
      </p>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div class="form-control">
          <label for="otp-code" class="label">
            <span class="label-text">ワンタイムパスワード (6桁)</span>
          </label>
          <input
            id="otp-code"
            v-model="otp"
            type="text"
            inputmode="numeric"
            pattern="[0-9]{6}"
            maxlength="6"
            placeholder="123456"
            class="input input-bordered w-full text-center text-2xl tracking-widest"
            required
            autocomplete="one-time-code"
          />
        </div>

        <ErrorMessage :message="error" />

        <button
          type="submit"
          class="btn btn-primary w-full"
          :disabled="loading || otp.length !== 6"
        >
          <span v-if="loading" class="loading loading-spinner"></span>
          {{ loading ? "検証中..." : "ログイン" }}
        </button>

        <button
          type="button"
          class="btn btn-ghost w-full"
          :disabled="loading"
          @click="emit('back')"
        >
          戻る
        </button>
      </form>
    </div>
  </div>
</template>
