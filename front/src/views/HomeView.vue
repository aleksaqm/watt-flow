<script setup lang="ts">
import LoginForm from '../components/LoginForm.vue';
</script>


<template>
  <main>
    <div class="flex flew-row items-center justify-center h-full">
      <LoginForm class=""></LoginForm>
      <RouterLink to="/register">Register</RouterLink>
    </div>

  </main>
</template>
