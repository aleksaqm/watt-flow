<script setup lang="ts">
import HouseholdInfo from '@/components/household/HouseholdInfo.vue';
import SimulatorAvailabilty from '../../components/household/SimulatorAvailabilty.vue';

</script>


<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <HouseholdInfo> </HouseholdInfo>
      <SimulatorAvailabilty></SimulatorAvailabilty>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
