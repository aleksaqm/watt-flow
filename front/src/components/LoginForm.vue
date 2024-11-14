<script setup lang="ts">
import Button from '../shad/components/ui/button/Button.vue';
import Input from '../shad/components/ui/input/Input.vue';
import { useToast } from '../shad/components/ui/toast/use-toast'
import Toaster from '../shad/components/ui/toast/Toaster.vue';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import axios from 'axios'
import * as z from 'zod'
import { useRouter } from 'vue-router';

const { toast } = useToast()
const router = useRouter()


const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
}))

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const submitForm = async (formData: { username: string; password: string }) => {
  try {
    const response = await axios.post('/api/login', formData)
    console.log('Response:', response.data)
    localStorage.setItem("authToken", response.data['token'])
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
  <div class="w-1/2 p-10 flex flex-col justify-center items-center bg-white">
    <div class="p-10 flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-lg">Sign In</span>
      <form class="w-full space-y-6" @submit="onSubmit">
        <FormField name="username" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter your username" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.username">{{ errors.username }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="password" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter your password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.password">{{ errors.password }}
            </FormMessage>
          </FormItem>
        </FormField>

        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Sign In
        </Button>
      </form>
    </div>
  </div>
  <Toaster />
</template>

<style scoped></style>
