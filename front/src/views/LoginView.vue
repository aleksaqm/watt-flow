<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import RegisterLink from '@/components/RegisterLink.vue'
</script>

<template>
  <main class="flex h-screen items-center justify-center">
    <div class="flex w-full max-w-4xl shadow-lg">
      <LoginForm />
      <RegisterLink />
    </div>
  </main>
</template>

<style scoped>
html, body {
  height: 100%;
}
</style>
