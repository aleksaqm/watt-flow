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
import Spinner from '../Spinner.vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import axios from 'axios'
import * as z from 'zod'
import { ref } from 'vue'

const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  jmbg: z.string().length(13, "JMBG must have exactly 13 digits.")
    .regex(/^\d{13}$/, "JMBG can have only digits."),
  email: z.string().email(),
}));

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const { toast } = useToast()
const router = useRouter()
const loading = ref(false)

const profilePicture = ref<File | null>(null)
const profilePicturePreview = ref<string | null>(null)

const MAX_FILE_SIZE = 1 * 1024 * 1024

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files ? target.files[0] : null

  if (file && file.size > MAX_FILE_SIZE) {
    toast({
      title: 'File too large',
      description: `Please upload an image smaller than ${MAX_FILE_SIZE / (1024 * 1024)} MB.`,
      variant: 'destructive'
    })
    profilePicture.value = null
    profilePicturePreview.value = null
    return
  }

  profilePicture.value = file
  profilePicturePreview.value = file ? URL.createObjectURL(file) : null
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

const submitForm = async (formData: { username: string; jmbg: string; email: string }) => {
  try {
    let profileImageBase64 = ''
    if (profilePicture.value) {
      loading.value = true
      profileImageBase64 = await convertToBase64(profilePicture.value)
    } else {
      console.log("No image selected")
      toast({
        title: 'No image selected!',
        description: "No profile image was selected during registration",
        variant: 'default'
      })
      return
    }

    const data = {
      username: formData.username,
      jmbg: formData.jmbg,
      email: formData.email,
      profile_image: profileImageBase64,
    }
    console.log(data)

    const response = await axios.post('/api/user/clerk/new', data)
    console.log(response)
    loading.value = false
    toast({
      title: 'Operation Successful',
      description: 'You successfully added new clerk account!',
      variant: 'default'
    })
    router.push({ name: 'manageClerks' })
  } catch (error: any) {
    loading.value = false
    const errorMessage = 'Server error!. Please check entered data and try again.'
    console.error('Error:', error)
    toast({
      title: 'Clerk Registration Failed',
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

  <div class="w-1/3 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <Spinner v-if="loading" />
    <div class="flex flex-col justify-center items-center gap-5 w-full" v-if="!loading">
      <span class="text-gray-800 text-2xl">New Clerk</span>
      <form class="w-full space-y-6" @submit="onSubmit">
        <FormField name="username" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter clerk username" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.username">{{ errors.username }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="jmbg" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>JMBG</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter clerk JMBG" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.jmbg">{{ errors.jmbg }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="email" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Email</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter clerk email" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.email">{{ errors.email }}
            </FormMessage>
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
          Submit
        </Button>
      </form>
    </div>
  </div>
  <Toaster />
</template>

<style scoped></style>
