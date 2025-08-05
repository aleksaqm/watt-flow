<script setup lang="ts">
import SearchForm from '@/components/household/SearchForm.vue';
import SearchResult from '@/components/household/SearchResult.vue';
import { ref } from 'vue';



const triggerSearch = ref(0)
const searchQuery = ref<{ City?: string; Street?: string; Number?: string; Id?: string ; Withoutowner?: boolean}>({})
const onSearch = (query: any) => {
  searchQuery.value = query
  triggerSearch.value++
}

</script>
<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <div class="mx-auto w-full p-7 flex flex-col bg-white shadow-lg">
        <SearchForm @search=onSearch></SearchForm>
        <div class="w-full my-7 text-lg">
          Results:
        </div>
        <SearchResult :query="searchQuery" :trigger-search="triggerSearch"></SearchResult>
      </div>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
