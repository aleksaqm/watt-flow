<template>
  <div :class="[
    'flex items-center justify-center w-24 h-12 rounded-xl text-sm font-medium border cursor-pointer',
    isAvailable ? 'text-gray-800 hover:bg-violet-400 hover:text-white' : 'bg-gray-200 text-gray-500 hover:bg-gray-100'
  ]">
    {{ formattedTime }}
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

// Props
const props = defineProps({
  startHour: {
    type: Number,
    required: true,
  },
  startMinute: {
    type: Number,
    required: true,
  },
  isAvailable: {
    type: Boolean,
    default: true,
  },
});

const formattedTime = computed(() =>
  `${String(props.startHour).padStart(2, "0")}:${String(props.startMinute).padStart(2, "0")}`
);
</script>

<style scoped></style>
