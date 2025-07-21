<script setup lang="ts">
import SearchResult from '@/components/household/SearchResult.vue';
import { ref, onMounted } from 'vue';
import { useUserStore } from '@/stores/user';

const userStore = useUserStore();
const triggerSearch = ref(0)
const searchQuery = ref<{ ownerid?: string }>({})

onMounted(() => {
  if (userStore.id) {
    searchQuery.value = { ownerid: userStore.id.toString() }
    triggerSearch.value++
  }
})

</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <div class="w-full text-center my-10 text-xl">
        My Households
      </div>
      <SearchResult 
        :query="searchQuery" 
        :trigger-search="triggerSearch"
        mode="my-households"
      ></SearchResult>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
