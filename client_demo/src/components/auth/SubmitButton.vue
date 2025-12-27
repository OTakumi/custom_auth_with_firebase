<script setup lang="ts">
interface Props {
  loading: boolean;
  isSignUp: boolean;
}

defineProps<Props>();
</script>

<template>
  <button type="submit" class="btn btn-primary w-full" :disabled="loading">
    <span v-if="loading" class="loading loading-spinner"></span>
    {{ loading ? "処理中..." : isSignUp ? "登録" : "ログイン" }}
  </button>
</template>
