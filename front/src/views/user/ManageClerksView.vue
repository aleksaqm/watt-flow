<script setup lang="ts">
import ClerksTable from '@/components/user/ClerksTable.vue';
import UserSearchForm from '@/components/user/UserSearchForm.vue';
import Button from '@/shad/components/ui/button/Button.vue';
import { ref } from 'vue';
import { RouterLink } from 'vue-router';


const triggerSearch = ref(0)
const searchQuery = ref<{ FirstName?: string; LastName?: string; Username?: string; Role: string; }>({ Role: "Clerk" })
const onSearch = (query: any) => {
  searchQuery.value = query
  triggerSearch.value++
}

</script>

<template>
  <main>
    <div class="w-10/12 h-screen wrapper">
      <div class="w-full text-center mt-10 mb-4 text-2xl">
        Clerks
      </div>
      <div class="flex justify-center mb-6">
        <RouterLink to="/manage/clerks/new">
          <Button variant="outline">Add new clerk</Button>
        </RouterLink>
      </div>
      <UserSearchForm @clerk-search=onSearch></UserSearchForm>
      <ClerksTable :query="searchQuery" :trigger-search="triggerSearch"></ClerksTable>
    </div>
  </main>
</template>
<style scoped>
.wrapper {
  margin: 0 auto;
}
</style>
