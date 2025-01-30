<script setup lang="ts">
import UnsentBills from '@/components/bills/UnsentBills.vue';
import Bills from '@/components/bills/Bills.vue';

</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <div class="w-full text-center my-10 text-xl">
        Unsent bills
      </div>
      <UnsentBills></UnsentBills>
      <div class="w-full text-center my-10 text-xl">
        Sent bills
      </div>
      <Bills></Bills>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
