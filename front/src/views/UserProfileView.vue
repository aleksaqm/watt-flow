<script setup lang="ts">
import UserProfileCard from '@/components/UserProfileCard.vue';
import Label from '../shad/components/ui/label/Label.vue'
import { getUsernameFromToken } from '@/utils/jwtDecoder';
import { getEmailFromToken } from '@/utils/jwtDecoder';
import { getRoleFromToken } from '@/utils/jwtDecoder';
import { onMounted, ref } from 'vue';

const username = ref<string>("")
const email = ref<string>("")
const role = ref<string>("")

onMounted(() => {
  const user = getUsernameFromToken();
  const userEmail = getEmailFromToken();
  const userRole = getRoleFromToken();
  
  if (user) {
    username.value = user;
  }
  if (userEmail) {
    email.value = userEmail;
  }
  if (userRole) {
    role.value = userRole;
  }
})

</script>

<template>
  <UserProfileCard :username="username" :email="email" :role="role" />
</template>