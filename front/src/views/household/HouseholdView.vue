<script setup lang="ts">
import HouseholdInfo from '@/components/household/HouseholdInfo.vue';
import SimulatorAvailabilty from '../../components/household/SimulatorAvailabilty.vue';
import { useRoute } from 'vue-router';
import { computed, onMounted, ref } from 'vue';
import axios from 'axios';


const route = useRoute()
const householdId = computed(() => route.params.id)
const household = ref<any | null>(null)
const isLoading = ref(true)


const fetchHousehold = async () => {
  try {
    isLoading.value = true
    const result = await axios.get(`/api/household/${householdId.value}`)
    household.value = result.data.data
    console.log(household.value)
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchHousehold()
})

</script>


<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <HouseholdInfo v-if="!isLoading" :household="household"> </HouseholdInfo>
      <SimulatorAvailabilty v-if="!isLoading " :device-id="household.device_address"></SimulatorAvailabilty>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
