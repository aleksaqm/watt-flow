<script setup lang="ts">
import Label from "@/shad/components/ui/label/Label.vue";

const props = defineProps({
  username: String,
  email: String,
  role: String,
  firstName: String,
  lastName: String,
});
</script>

<template>

  <div
    class="flex justify-center items-center min-w-screen min-h-screen bg-gradient-to-r from-[#1d1e26] to-[#4d596a]" v-if="props.username"
  >
    <div class="relative bg-white rounded-lg w-4/5 shadow-lg p-8 max-w-lg">
      <!-- Profile Picture -->
      <div class="absolute top-[-70px] left-1/2 transform -translate-x-1/2">
        <img
          :src="`http://localhost:8080/profile_pictures/${props.username}.jpg`"
          alt="User Profile"
          class="w-40 h-40 rounded-full border-4 border-white shadow-md object-cover"
        />
      </div>

      <!-- User Details -->
      <div class="mt-24 text-gray-800">
        <h2 class="text-2xl mb-10 font-bold text-center">
          {{ props.username }}
        </h2>

        <div class="space-y-5">
          <div class="grid grid-cols-2 gap-2 items-center">
            <Label class="text-lg text-left">First Name:</Label>
            <Label class="text-xl font-medium">{{ props.firstName }}</Label>
          </div>
          <div class="grid grid-cols-2 gap-2 items-center">
            <Label class="text-lg text-left">Last Name:</Label>
            <Label class="text-xl font-medium">{{ props.lastName }}</Label>
          </div>
          <div class="grid grid-cols-2 gap-2 items-center">
            <Label class="text-lg text-left">Email:</Label>
            <Label class="text-xl font-medium">{{ props.email }}</Label>
          </div>
          <div class="grid grid-cols-2 gap-2 items-center">
            <Label class="text-lg text-left">Role:</Label>
            <Label class="text-xl font-medium">{{ props.role }}</Label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
