<script setup lang="ts">
import HouseholdInfo from '@/components/household/HouseholdInfo.vue';
import { useRoute } from 'vue-router';
import { computed, onMounted, ref } from 'vue';
import { useUserStore } from '@/stores/user';
import router from '@/router';
import axios from 'axios';

const route = useRoute()
const userStore = useUserStore()
const householdId = computed(() => route.params.id)
const household = ref<any | null>(null)
const isLoading = ref(true)

const fetchHousehold = async () => {
  try {
    isLoading.value = true
    const result = await axios.get(`/api/household/${householdId.value}`)
    
    if (result.data.data.owner_id !== userStore.id) {
      router.push({ name: 'my-households' })
      return
    }
    
    household.value = result.data.data
    console.log(household.value)
  } catch (err) {
    console.error(err)
    router.push({ name: 'my-households' })
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (userStore.role !== 'Regular') {
    router.push({ name: 'home' })
    return
  }
  
  fetchHousehold()
})

</script>

<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <HouseholdInfo v-if="!isLoading && household" :household="household"></HouseholdInfo>
      
      <!-- Placeholder for electricity usage component -->
      <div v-if="!isLoading && household" class="my-10">
        <div class="text-center text-xl text-gray-600 mb-5">
          Electricity Usage
        </div>
        <div class="bg-white shadow-lg rounded-lg p-8 text-center">
          <p class="text-gray-500">
            Electricity usage analytics will be displayed here
          </p>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
