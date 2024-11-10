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
import { Field, useForm } from 'vee-validate'
import {useRouter} from 'vue-router'
import { useToast } from '../shad/components/ui/toast/use-toast'
import Toaster  from '../shad/components/ui/toast/Toaster.vue';
import axios from 'axios'
import * as z from 'zod'
import { ref } from 'vue'

const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
  confirmPassword: z.string().min(6, { message: "Password must be at least 6 characters" }),
  email: z.string().email(),
}));

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const { toast } = useToast()
const router = useRouter()

const profilePicture = ref<File | null>(null)
const profilePicturePreview = ref<string | null>(null)

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  profilePicture.value = target.files ? target.files[0] : null
  profilePicturePreview.value = profilePicture.value ? URL.createObjectURL(profilePicture.value) : null
}

const convertToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => {
      if (reader.result) {
        resolve(reader.result as string)
      } else {
        reject('Failed to convert file to base64')
      }
    }
    reader.onerror = () => reject('Failed to read file')
    reader.readAsDataURL(file)
  })
}

const submitForm = async (formData: { username: string; password: string; email: string }) => {
  try {
    let profileImageBase64 = ''
    if (profilePicture.value) {
      profileImageBase64 = await convertToBase64(profilePicture.value)
    }

    const data = {
      username: formData.username,
      password: formData.password,
      email: formData.email,
      profile_image: profileImageBase64,
    }

    console.log(data)

    const response = await axios.post('/register', data)
    console.log('Response:', response.data)
    toast({
      title: 'Registration Successful',
      description: 'You will have to activate account before logging in!',
      variant: 'default'
    })
    router.push({name: 'home'})
  } catch (error) {
    console.error('Error:', error)
    toast({
      title: 'Registration Failed',
      description: 'Please check your information again and try again.',
      variant: 'destructive'
    })
  }
}

const onSubmit = handleSubmit((values) => {
  submitForm(values)
})
</script>

<template>
  <div class="w-1/3 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <div class="flex flex-col justify-center items-center gap-5 w-full">
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

        <h3 class="pt-2 text-black text-bold">Profile picture:</h3>
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
  <Toaster />
</template>

<style scoped></style>
