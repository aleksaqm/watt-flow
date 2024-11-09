<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useToast } from 'vue-toastification'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import axios from 'axios'
import * as z from 'zod'

const toast = useToast()

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
    toast.success('Login successful!')
    console.log('Response:', response.data)
  } catch (error) {
    toast.error('Login failed!')
    console.error('Error:', error)
  }
}

const onSubmit = handleSubmit((values) => {
  submitForm(values)
})
</script>

<template>
  <div class="w-1/2 p-10 flex flex-col justify-center items-center bg-white">
    <div class="p-10 flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-lg">Sign In</span>
      <form class="w-full space-y-6" @submit="onSubmit">
        <FormField name="username" v-slot="{ field }">
          <FormItem>
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter your username" />
            </FormControl>
            <FormMessage v-if="errors.username">{{ errors.username }}</FormMessage>
          </FormItem>
        </FormField>

        <FormField name="password" v-slot="{ field }">
          <FormItem>
            <FormLabel>Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter your password" />
            </FormControl>
            <FormMessage v-if="errors.password">{{ errors.password }}</FormMessage>
          </FormItem>
        </FormField>

        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Sign In
        </Button>
      </form>
    </div>
  </div>
</template>

<style scoped></style>
