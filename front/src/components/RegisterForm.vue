<script setup lang="ts">
import { Button } from '@/shad/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import { Input } from '@/shad/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import axios from 'axios'
import * as z from 'zod'
import { ref } from 'vue'


const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
  confirmPassword: z.string().min(6, { message: "Password must be at least 6 characters" }),
  email: z.string().email()
}));

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const profilePicturePreview = ref<string | null>(null)

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files ? target.files[0] : null
//   formData.profilePicture = file
    profilePicturePreview.value = file ? URL.createObjectURL(file) : null
}

const submitForm = async (formData: { username: string; password: string }) => {
  try {
    const response = await axios.post('/login', formData)
    console.log('Response:', response.data)
  } catch (error) {
    console.error('Error:', error)
  }
}

const onSubmit = handleSubmit((values) => {
  submitForm(values)
})
</script>

<template>
  <div class="w-1/3 p-10 flex flex-col justify-center items-center bg-white shadow-lg">
    <div class="p-10 flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-lg">Sign Up</span>
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

        <FormField name="confirmPassword" v-slot="{ field }">
          <FormItem>
            <FormLabel>Confirm password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter same password" />
            </FormControl>
            <FormMessage v-if="errors.confirmPassword">{{ errors.confirmPassword }}</FormMessage>
          </FormItem>
        </FormField>
        
        <FormField name="email" v-slot="{ field }">
          <FormItem>
            <FormLabel>Email</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter your email" />
            </FormControl>
            <FormMessage v-if="errors.email">{{ errors.email }}</FormMessage>
          </FormItem>
        </FormField>

        <h3 class="pt-2 text-black b">Profile picture:</h3>
        <FormField name="profilePicture">
          <FormItem>
            <FormControl>
              <input type="file" @change="onFileChange" accept="image/*" />
            </FormControl>
          </FormItem>
        </FormField>

        <div v-if="profilePicturePreview" class="mt-4">
          <img :src="profilePicturePreview" alt="Profile preview" class="max-w-28 h-28 rounded-full object-cover" />
        </div>

        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Sign In
        </Button>
      </form>
    </div>
  </div>
</template>

<style scoped></style>
