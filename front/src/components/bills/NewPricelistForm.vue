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
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import axios from 'axios'
import * as z from 'zod'
import { ref } from 'vue'

const formSchema = toTypedSchema(z.object({
  month: z.number().min(1, { message: "Month must be between 1 and 12" }).max(12, { message: "Month must be between 1 and 12" }),
  year: z.number().min(1970, "Year not valid!").max(2050, "Year not valid!"),
  red: z.number().min(0.01),
  blue: z.number().min(0.01),
  green: z.number().min(0.01),
  tax: z.number().min(0.01),
  bill_power: z.number().min(0.01),
}));

const { handleSubmit, errors } = useForm({
  validationSchema: formSchema,
})

const { toast } = useToast()
const loading = ref(false)

const emit = defineEmits(['pricelistCreated']);


const submitForm = async (formData: { month: number; year: number; red: number; green: number; blue: number; tax: number; bill_power: number; }) => {
  try {
    const now = new Date()
    if (now.getMonth() + 1 >= formData.month || now.getFullYear() > formData.year) {
      toast({
        title: 'Pricelist creation failed',
        description: "Invalid dates. Future dates must be used!",
        variant: 'destructive'
      })
      loading.value = false;
      return;
    }
    const response = await axios.post('/api/pricelist', formData)
    console.log(response)
    loading.value = false
    emit('pricelistCreated');
  } catch (error: any) {
    loading.value = false
    const errorMessage = 'Invalid date. Pricelist already exists for entered date.'
    console.error('Error:', error)
    toast({
      title: 'Pricelist creation failed',
      description: errorMessage,
      variant: 'destructive'
    })
  }
}

const onSubmit = handleSubmit((values) => {
  loading.value = true
  submitForm(values)
  console.log("submit")
})
</script>

<template>
  <div class="w-full flex flex-col justify-center items-center">
    <Spinner v-if="loading" />
    <div class="flex flex-col justify-center items-center gap-5 w-full" v-if="!loading">
      <form class="w-full space-y-6" @submit.prevent="onSubmit">
        <span class="text-sm font-semibold text-gray-900">Valid from</span>
        <div class="flex flex-row gap-10">
          <FormField name="month" v-slot="{ field }">
            <FormItem class="relative pb-2">
              <FormControl>
                <Input type="number" v-bind="field" placeholder="Enter a month" />
              </FormControl>
              <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.month">{{ errors.month }}
              </FormMessage>
            </FormItem>
          </FormField>

          <FormField name="year" v-slot="{ field }">
            <FormItem class="relative pb-2">
              <FormControl>
                <Input type="number" v-bind="field" placeholder="Enter a year" />
              </FormControl>
              <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.year">{{ errors.year }}
              </FormMessage>
            </FormItem>
          </FormField>
        </div>

        <FormField name="red" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Red zone price</FormLabel>
            <FormControl>
              <Input type="number" step=".01" min='0' v-bind="field" placeholder="Enter a red zone price" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.red">{{ errors.red }}
            </FormMessage>
          </FormItem>
        </FormField>
        <FormField name="blue" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Blue zone price</FormLabel>
            <FormControl>
              <Input type="number" step=".01" min='0' v-bind="field" placeholder="Enter a blue zone price" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.blue">{{ errors.blue }}
            </FormMessage>
          </FormItem>
        </FormField>
        <FormField name="green" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Green zone price</FormLabel>
            <FormControl>
              <Input type="number" step=".01" min='0' v-bind="field" placeholder="Enter a green zone price" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.green">{{ errors.green }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="bill_power" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Billing power</FormLabel>
            <FormControl>
              <Input type="number" step=".01" min='0' v-bind="field" placeholder="Enter a billing power factor" />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.bill_power">{{ errors.bill_power }}
            </FormMessage>
          </FormItem>
        </FormField>

        <FormField name="tax" v-slot="{ field }">
          <FormItem class="relative pb-2">
            <FormLabel>Tax</FormLabel>
            <FormControl>
              <Input type="number" step=".01" min='0' v-bind="field" placeholder="Enter tax percentage  " />
            </FormControl>
            <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.tax">{{ errors.tax }}
            </FormMessage>
          </FormItem>
        </FormField>


        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Submit
        </Button>
      </form>
    </div>
  </div>
</template>

<style scoped></style>
