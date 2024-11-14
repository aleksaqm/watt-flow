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
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import axios from 'axios'
import * as z from 'zod'
import { useRouter } from 'vue-router';

const { toast } = useToast()
const router = useRouter()


const formSchema = toTypedSchema(z.object({
  defaultPassword: z.string().min(2, { message: "Password must be at least 6 characters" }).max(50, { message: "Password cannot exceed 50 characters" }),
  newPassword: z.string().min(6, { message: "Password must be at least 6 characters" }),
  confirmPassword: z.string().min(6, { message: "Password must be at least 6 characters" }),
}))

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const submitForm = async (formData: { defaultPassword: string; newPassword: string, confirmPassword: string }) => {
  if (formData.newPassword != formData.confirmPassword) {
    toast({
      title: 'Password change failed',
      description: "Passwords don't match",
      variant: 'destructive'
    })
    return
  }
  console.log(formData)
  try {
    const response = await axios.post('/api/admin/password', { "old_password": formData.defaultPassword, "new_password": formData.newPassword })
    console.log('Response:', response.data)
    toast({
      title: 'Password changed successfully',
      variant: 'default'
    })
    router.push({ name: 'login' })
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
      <!-- <span class="text-gray-800 text-lg">Sign In</span> -->
      <form class="w-full space-y-6" @submit="onSubmit">
        <FormField name="defaultPassword" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Default Password</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter default password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.defaultPassword">{{
              errors.defaultPassword }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="newPassword" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>New Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter your new password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.newPassword">{{ errors.newPassword }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="confirmPassword" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Confirm Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Confirm your new password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.confirmPassword">{{
              errors.confirmPassword }}
            </FormMessage>
          </FormItem>
        </FormField>

        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Change
        </Button>
      </form>
    </div>
  </div>
  <Toaster />
</template>

<style scoped></style>
