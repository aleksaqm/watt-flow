<script setup lang="ts">
import LoginForm from '@/components/auth/LoginForm.vue'
import RegisterLink from '@/components/auth/RegisterLink.vue'
</script>

<template>
  <main class="flex-1 flex items-center justify-center">
    <div class="flex w-full max-w-4xl shadow-lg">
      <LoginForm />
      <RegisterLink />
    </div>
  </main>
</template>

<style scoped>
html,
body {
  height: 100%;
}
</style>
