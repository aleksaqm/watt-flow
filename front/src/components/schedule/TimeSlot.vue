<template>
  <div :class="[
    'flex items-center justify-center w-24 h-12 rounded-xl text-sm font-medium border cursor-pointer', buttonFormat()
  ]">
    {{ formattedTime }}
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { any, number } from "zod";

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
  past: {
    type: Boolean,
    default: false,
  },
  hasClerk: {
    type: Boolean,
    default: false,
  }
});

const formattedTime = computed(() =>
  `${String(props.startHour).padStart(2, "0")}:${String(props.startMinute).padStart(2, "0")}`
);

const buttonFormat = () => {
  if(props.hasClerk){
    if (!props.past) {
      if (props.isAvailable) {
        return 'text-gray-800 hover:bg-violet-400 hover:text-white'
      } else {
        return 'bg-gray-200 text-gray-500 hover:bg-gray-100'
      }
    } else {
      if (props.isAvailable) {
        return 'cursor-not-allowed bg-gray-100 hover:bg-gray-100'
      } else {
        return 'bg-gray-200 text-gray-500 hover:bg-violet-100'
      }
    }
  } else{
    return 'cursor-not-allowed bg-gray-100 hover:bg-gray-100'
  }
  
}
</script>

<style scoped></style>
