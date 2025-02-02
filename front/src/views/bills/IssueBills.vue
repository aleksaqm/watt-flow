<script setup lang="ts">
import UnsentBills from '@/components/bills/UnsentBills.vue';
import Bills from '@/components/bills/Bills.vue';
import { ref } from 'vue';


const triggerQuery = ref(0)
const onSent = () => {
  triggerQuery.value++
}
</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <div class="w-full text-center my-10 text-xl">
        Unsent bills
      </div>
      <UnsentBills @bill-sent="onSent"></UnsentBills>
      <div class="w-full text-center my-10 text-xl">
        Sent bills
      </div>
      <Bills :trigger-query="triggerQuery"></Bills>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
