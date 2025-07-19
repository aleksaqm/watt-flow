<script setup lang="ts">
import UserProfileCard from "@/components/user/UserProfileCard.vue";
import { getUserIdFromToken } from "@/utils/jwtDecoder";
import axios from "axios";
import { onMounted, ref } from "vue";
interface UserProfile {
  first_name: string;
  last_name: string;
  username: string;
  email: string;
  role: string;
  status: string;
}

const user = ref<UserProfile | null>(null);


onMounted(async () => {
  const userId = getUserIdFromToken();

  if (userId) {
    const response = await axios.get(`/api/user/${userId}`);
    if (response.status === 200) {
      user.value = response.data.data;
    } else {
      console.error("Failed to fetch user data");
    }
  }
});
</script>

<template>
  <UserProfileCard v-if="user" :username="user?.username" :email="user?.email" :role="user?.role" :first-name="user?.first_name" :last-name="user?.last_name" />
</template>
