<script setup lang="ts">
import Label from '@/shad/components/ui/label/Label.vue'

const props = defineProps<{ username: string; email: string; role: string }>()

</script>

<template>
  <div class="flex justify-center items-center min-w-screen min-h-screen bg-gradient-to-r from-[#1d1e26] to-[#4d596a]">
    <div class="relative bg-white rounded-lg w-4/5 shadow-lg p-8 max-w-lg text-center">
      <div class="absolute top-[-70px] left-1/2 transform -translate-x-1/2">
        <img :src="`http://localhost:8080/profile_pictures/${props.username}.jpg`" alt="User Profile"
          class="w-40 h-40 rounded-full border-4 border-white shadow-md object-cover" />
      </div>

      <div class="mt-20"> <!-- Adjust margin-top to make space for the larger image -->
        <h2 class="text-2xl mb-10 font-bold text-gray-800">{{ props.username }}</h2>
        <div class="mb-5">
          <Label class="text-lg">Email: </Label>
          <Label class="text-xl">{{ props.email }}</Label>
        </div>
        <div class="mb-5">
          <Label class="text-lg">Role: </Label>
          <Label class="text-xl">{{ props.role }}</Label>
        </div>
      </div>
    </div>
  </div>
</template>
