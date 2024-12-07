<script setup lang="ts">

import Button from '@/shad/components/ui/button/Button.vue';
import Input from '@/shad/components/ui/input/Input.vue';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/shad/components/ui/form'
import { onMounted, ref } from 'vue';
import { useUserStore } from '@/stores/user';

const userStore = useUserStore();
const searchQuery = ref<{ City?: string; Street?: string; Number?: string; Id?: string ; WithoutOwner?: boolean}>({})
const emit = defineEmits(['search'])
const isAdmin = ref<boolean>(true)

const onSubmit = () => {
  if (!isAdmin){
    searchQuery.value.WithoutOwner = true;
  }
  emit('search', searchQuery.value)
}

onMounted(()=>{
  if (userStore.role == "Regular"){
    isAdmin.value = false;
  }
})

</script>
<template>
  <div>
    <div v-if="isAdmin" class="w-full text-center my-10 text-xl">
      Search household
    </div>
    <div v-if="!isAdmin" class="w-full text-center my-10 text-xl">
      Find your household 
    </div>

    <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10">
      <FormField name="city" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">City:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-model="searchQuery.City" placeholder="Enter city" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="street" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">Street:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-model="searchQuery.Street" placeholder="Enter street" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="number" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">Number:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-model="searchQuery.Number" placeholder="Enter number" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="id" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">ID:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-model="searchQuery.Id" placeholder="Enter household id" />
          </FormControl>
        </FormItem>
      </FormField>

      <Button type="button" @click="onSubmit"
        class="bg-indigo-500 text-white w-32 ml-10 hover:bg-gray-600 rounded-full py-2">
        Search
      </Button>
    </form>


    <div class="w-full my-7 text-lg">
      Results:
    </div>
  </div>
</template>
<style scoped></style>
