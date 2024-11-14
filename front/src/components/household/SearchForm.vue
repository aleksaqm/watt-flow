<script setup lang="ts">

import Button from '@/shad/components/ui/button/Button.vue';
import Input from '@/shad/components/ui/input/Input.vue';
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import router from '@/router';
import { toast } from '@/shad/components/ui/toast';
import { toTypedSchema } from '@vee-validate/zod';
import axios from 'axios';
import { useForm } from 'vee-validate';
import { z } from 'zod';


const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
}))

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const submitForm = async (formData: { username: string; password: string }) => {
  try {
    const response = await axios.post('/login', formData)
    console.log('Response:', response.data)
    localStorage.setItem("authToken", response.data['token'])
    toast({
      title: 'Login Successful',
      description: 'You have successfully logged in!',
      variant: 'default'
    })
    router.push({ name: 'home' })
  } catch (error: any) {
    console.error('Error:', error)
    const errorMessage = error.response?.data?.error || 'An unexpected error occurred'
    toast({
      title: 'Login Failed',
      description: errorMessage,
      variant: 'destructive'
    })
  }
}

const onSubmit = handleSubmit((values) => {
  submitForm(values)
})

</script>
<template>
  <div>
    <div class="w-full text-center my-10 text-xl">
      Search household
    </div>

    <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10"
      @submit="onSubmit">

      <FormField name="city" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">City:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-bind="field" placeholder="Enter city" />
          </FormControl>
        </FormItem>
        <FormMessage v-if="errors.username">{{ errors.username }}</FormMessage>
      </FormField>

      <FormField name="street" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">Street:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-bind="field" placeholder="Enter street" />
          </FormControl>
        </FormItem>
        <FormMessage v-if="errors.username">{{ errors.username }}</FormMessage>
      </FormField>

      <FormField name="number" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">Number:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-bind="field" placeholder="Enter number" />
          </FormControl>
        </FormItem>
        <FormMessage v-if="errors.username">{{ errors.username }}</FormMessage>
      </FormField>

      <FormField name="id" v-slot="{ field }">
        <FormItem class="flex items-center gap-5">
          <FormLabel class="w-1/4 text-right">ID:</FormLabel>
          <FormControl class="w-3/4">
            <Input type="text" v-bind="field" placeholder="Enter household id" />
          </FormControl>
        </FormItem>
        <FormMessage v-if="errors.username">{{ errors.username }}</FormMessage>
      </FormField>

      <Button type="submit" class="bg-indigo-500 text-white w-32 ml-10 hover:bg-gray-600 rounded-full py-2">
        Search
      </Button>
    </form>


    <div class="w-full my-7 text-lg">
      Results:
    </div>
  </div>
</template>
<style scoped></style>
