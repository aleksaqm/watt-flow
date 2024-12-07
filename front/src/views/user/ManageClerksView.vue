<script setup lang="ts">
import ClerksTable from '@/components/user/ClerksTable.vue';
import Button from '@/shad/components/ui/button/Button.vue';
import { RouterLink } from 'vue-router';


</script>

<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <span class="text-gray-800 text-2xl">Clerks</span>
      <RouterLink to="/manage/clerks/new"><Button variant="outline">Add new clerk</Button></RouterLink>
      <ClerksTable></ClerksTable>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
  display: flex;
  flex: flex 1;
  justify-content: center;
  align-items: center;
  flex-direction: column;

}
</style>
