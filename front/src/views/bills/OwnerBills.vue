<script setup lang="ts">
import OwnerBillsTable from '@/components/bills/OwnerBillsTable.vue';
import { ref } from 'vue';


</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <OwnerBillsTable></OwnerBillsTable>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>