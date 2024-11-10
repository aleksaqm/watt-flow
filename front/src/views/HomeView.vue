<script setup lang="ts">
import { getUsernameFromToken } from '@/utils/jwtDecoder';
import { ref, onMounted } from 'vue';

// Create a reactive username variable
const username = ref<string | null>(null);

onMounted(() => {
//   const decodedUsername = getUsernameFromToken();
  console.log(1)
  // Check if the username is defined, then assign it to the ref
//   if (decodedUsername) {
//     username.value = decodedUsername;
//   } else {
//     console.log("No username found in token");
//   }
});
</script>

<template>
  <main>
    <h1>HOME PAGE</h1>
    <p v-if="username">Welcome, {{ username }}!</p>
    <p v-else>No username found in token.</p>
  </main>    
</template>

<style></style>
