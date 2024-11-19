<script setup lang="ts">
import SearchForm from '@/components/household/SearchForm.vue';
import SearchResult from '@/components/household/SearchResult.vue';
import { ref } from 'vue';



const triggerSearch = ref(0)
const searchQuery = ref<{ City?: string; Street?: string; Number?: string; Id?: string }>({})
const onSearch = (query: any) => {
  searchQuery.value = query
  triggerSearch.value++
}

</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <SearchForm @search=onSearch></SearchForm>
      <SearchResult :query="searchQuery" :trigger-search="triggerSearch"></SearchResult>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
