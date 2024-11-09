<script setup lang="ts">
import Button from '../shad/components/ui/button/Button.vue';
import Input from '../shad/components/ui/input/Input.vue';
</script>
<template>
  <div class="w-1/4">
    <div class="border-gray-500 border-2 p-10 m-5 rounded-2xl flex flex-col justify-center items-center gap-5">
      <span class="text-gray-800 text-lg">Welcome</span>
      <Input placeholder="Email" type="email"></Input>
      <Input placeholder="Password" type="password"></Input>

      <Button>Sign In</Button>
    </div>

  </div>
</template>

<style scoped></style>
