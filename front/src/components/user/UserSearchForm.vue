<script setup lang="ts">

import Button from '@/shad/components/ui/button/Button.vue';
import Input from '@/shad/components/ui/input/Input.vue';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/shad/components/ui/form'
import { ref } from 'vue';


const searchQuery = ref<{ FirstName?: string; LastName?: string; Username?: string; Role: string; }>({ Role: "Clerk" })
const emit = defineEmits(['clerk-search'])

const onSubmit = () => {
  emit('clerk-search', searchQuery.value)
}

</script>
<template>
  <div class="mx-auto w-full p-7 flex flex-col bg-white shadow-lg">

    <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10">
      <FormField name="firstName" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/2 text-right">First name:</FormLabel>
          <FormControl class="w-1/2">
            <Input type="text" v-model="searchQuery.FirstName" placeholder="Enter first name" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="street" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/2 text-right">Last name:</FormLabel>
          <FormControl class="w-1/2">
            <Input type="text" v-model="searchQuery.LastName" placeholder="Enter last name" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="number" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/2 text-right">Username:</FormLabel>
          <FormControl class="w-1/2">
            <Input type="text" v-model="searchQuery.Username" placeholder="Enter username" />
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
