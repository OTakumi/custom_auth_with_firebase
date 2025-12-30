<script setup lang="ts">
interface Props {
  isSignUp: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  toggle: [];
}>();
</script>

<template>
  <div class="text-center mt-4">
    <p class="text-sm">
      {{ isSignUp ? "既にアカウントをお持ちですか？" : "アカウントをお持ちでないですか？" }}
      <button type="button" @click="emit('toggle')" class="link link-primary ml-1">
        {{ isSignUp ? "ログイン" : "新規登録" }}
      </button>
    </p>
  </div>
</template>
