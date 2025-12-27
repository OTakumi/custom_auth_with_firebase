<script setup lang="ts">
interface Props {
  id: string;
  label: string;
  type: string;
  placeholder: string;
  modelValue: string;
  autocomplete?: string;
  minlength?: number;
  required?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  required: false,
});

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit("update:modelValue", target.value);
};
</script>

<template>
  <div class="form-control w-full">
    <label class="label" :for="id">
      <span class="label-text">{{ label }}</span>
    </label>
    <input
      :id="id"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :minlength="minlength"
      :required="required"
      @input="handleInput"
      class="input input-bordered w-full"
    />
  </div>
</template>
