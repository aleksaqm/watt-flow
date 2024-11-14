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
import axios from 'axios'
import { Field, useForm } from 'vee-validate'
import * as z from 'zod'
import { useToast } from '../shad/components/ui/toast/use-toast'
import { defineEmits } from 'vue'
import Toaster from '../shad/components/ui/toast/Toaster.vue';

const props = defineProps<{ url: string; role: string }>()

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
const emit = defineEmits(['userCreated'])


const submitForm = async (formData: { username: string; password: string; email: string }) => {
  try {
    const data = {
      username: formData.username,
      password: formData.password,
      email: formData.email,
      role: props.role
    }

    const response = await axios.post(props.url, data)
    console.log('Response:', response.data)
    emit('userCreated')
    toast({
      title: 'Creation Successful',
      variant: 'default'
    })
  } catch (error) {
    console.error('Error:', error)
    toast({
      title: 'Creation Failed',
      description: 'Please check information again and try again.',
      variant: 'destructive'
    })
  }
}

const onSubmit = handleSubmit((values) => {
  if (values.password != values.confirmPassword){
    toast({
      title: 'Creation Failed',
      description: "Passwords aren't the same.",
      variant: 'destructive'
    })
    return
  }
  submitForm(values)
})


</script>

<template>
    <main>
        <form class="w-full space-y-6" @submit="onSubmit">
        <FormField name="username" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter username" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.username">{{ errors.username }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="password" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.password">{{ errors.password }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="confirmPassword" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Confirm password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="field" placeholder="Enter same password" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.confirmPassword">{{
              errors.confirmPassword }}</FormMessage>
          </FormItem>
        </FormField>

        <FormField name="email" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Email</FormLabel>
            <FormControl>
              <Input type="text" v-bind="field" placeholder="Enter email" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.email">{{ errors.email }}
            </FormMessage>
          </FormItem>
        </FormField>
        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Create
        </Button>
      </form>
      <Toaster />
    </main>
</template>