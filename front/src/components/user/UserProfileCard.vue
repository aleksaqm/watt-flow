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
          :src="`http://localhost:80/profile_pictures/${props.username}.jpg`"
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
          <div class="flex flex-wrap justify-center items-center gap-2">
            <Label class="text-lg font-semibold">First Name:</Label>
            <Label class="text-xl font-medium break-words">{{ props.firstName }}</Label>
          </div>
          <div class="flex flex-wrap justify-center items-center gap-2">
            <Label class="text-lg font-semibold">Last Name:</Label>
            <Label class="text-xl font-medium break-words">{{ props.lastName }}</Label>
          </div>
          <div class="flex flex-wrap justify-center items-center gap-2">
            <Label class="text-lg font-semibold">Email:</Label>
            <Label class="text-xl font-medium break-words break-all">{{ props.email }}</Label>
          </div>
          <div class="flex flex-wrap justify-center items-center gap-2">
            <Label class="text-lg font-semibold">Role:</Label>
            <Label class="text-xl font-medium break-words">{{ props.role }}</Label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
