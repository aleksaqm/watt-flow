<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Button } from '@/shad/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shad/components/ui/card';
import { Input } from '@/shad/components/ui/input';
import { Label } from '@/shad/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shad/components/ui/select';
import { Checkbox } from '@/shad/components/ui/checkbox';
import { CreditCard, Lock, Mail, User, Shield } from 'lucide-vue-next';
import axios from 'axios';
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { useToast } from '../../shad/components/ui/toast/use-toast';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { z } from 'zod';

const { toast } = useToast();

const formSchema = toTypedSchema(
  z.object({
    firstName: z
      .string()
      .min(1, { message: 'First name is required' })
      .min(2, { message: 'First name must be at least 2 characters' }),
    lastName: z
      .string()
      .min(1, { message: 'Last name is required' })
      .min(2, { message: 'Last name must be at least 2 characters' }),
    cardNumber: z
      .string()
      .min(1, { message: 'Card number is required' })
      .regex(/^\d{4}\s\d{4}\s\d{4}\s\d{4}$/, { 
        message: 'Card number must be in format 1234 5678 9012 3456' 
      }),
    expiryMonth: z
      .string()
      .min(1, { message: 'Expiry month is required' })
      .regex(/^(0[1-9]|1[0-2])$/, { message: 'Invalid month' }),
    expiryYear: z
      .string()
      .min(1, { message: 'Expiry year is required' })
      .refine((year) => {
        const currentYear = new Date().getFullYear();
        const selectedYear = parseInt(year);
        return selectedYear >= currentYear;
      }, { message: 'Card has expired' }),
    cvc: z
      .string()
      .min(1, { message: 'CVC is required' })
      .regex(/^\d{3,4}$/, { message: 'CVC must be 3 or 4 digits' }),
  })
);

const { handleSubmit, errors, setFieldValue, values, defineField } = useForm({
  validationSchema: formSchema,
  initialValues: {
    firstName: '',
    lastName: '',
    cardNumber: '',
    expiryMonth: '',
    expiryYear: '',
    cvc: '',
  }
});

const [firstName] = defineField('firstName');
const [lastName] = defineField('lastName');
const [cardNumber] = defineField('cardNumber');
const [expiryMonth] = defineField('expiryMonth');
const [expiryYear] = defineField('expiryYear');
const [cvc] = defineField('cvc');

interface Bill {
  id: number;
  billing_date: string;
  price: number;
  spent_power: number;
  owner: {
    username: string;
    email: string;
  };
}

const route = useRoute();
const router = useRouter();

const isLoading = ref(false);
const bill = ref<Bill | null>(null);
const taxToPay = ref<number>(0);

const months = new Map<string, string>([
    ['01', 'January'],
    ['02', 'February'],
    ['03', 'March'],
    ['04', 'April'],
    ['05', 'May'],
    ['06', 'June'],
    ['07', 'July'],
    ['08', 'August'],
    ['09', 'September'],
    ['10', 'October'],
    ['11', 'November'],
    ['12', 'December'],
])

onMounted(async () => {
  const id = route.params.billId;
  if (!id) {
    throw new Error('Bill ID is required');
  }
  
  try {
    const params = {
      id: id,
    };
    const response = await axios.get('/api/bills', { params });
    if (response.data) {
      console.log(response.data);
      var billingMonth = months.get(response.data.bill.billing_date.split('-')[1])
      console.log("MONTH", billingMonth)
      const data: Bill = {
        id: response.data.bill.id,
        billing_date: billingMonth + " " + response.data.bill.billing_date.split('-')[0],
        price: response.data.bill.price,
        spent_power: 150,
        owner: { username: 'vladimir', email: 'vladimir@example.com' }
      };

      taxToPay.value = (response.data.bill.price / 100) * response.data.bill.pricelist.tax;
      bill.value = data;
    }
  } catch (error: any) {
    let errorMessage = 'An unexpected error occurred';
    
    if (axios.isAxiosError(error)) {
      console.log(error)
      errorMessage = error.response?.data.error;
    } else if (error instanceof Error) {
      errorMessage = error.message;
    }
    toast({
      title: 'Error:',
      description: errorMessage,
      variant: "destructive",
    });
  }
});

watch(() => cardNumber.value, (newValue, oldValue) => {
  if (!newValue) return;
  
  const digitsOnly = newValue.replace(/\D/g, '');

  if (digitsOnly === oldValue?.replace(/\D/g, '')) {
    return;
  }

  const parts = [];
  for (let i = 0; i < digitsOnly.length; i += 4) {
    parts.push(digitsOnly.substring(i, i + 4));
  }
  
  const formattedValue = parts.join(' ').trim();

  if (cardNumber.value !== formattedValue) {
    cardNumber.value = formattedValue;
  }
});

const onSubmit = handleSubmit(async (formData) => {
  isLoading.value = true;
  
  try {
    await new Promise((resolve) => setTimeout(resolve, 1500));

    const response = await axios.put('/api/bills/pay/' + route.params.billId);
    
    console.log("Payment processed for bill:", bill.value?.id, "Details:", formData);
    
    router.push('/bills/pay/success');
  } catch (error) {
    let errorMessage = 'An unexpected error occurred';
    
    if (axios.isAxiosError(error)) {
      console.log(error)
      errorMessage = error.response?.data.error;
    } else if (error instanceof Error) {
      errorMessage = error.message;
    }

    toast({
      title: 'Error',
      description: errorMessage,
      variant: "destructive",
    });
  } finally {
    isLoading.value = false;
  }
});

const currentYear = new Date().getFullYear();
const years = Array.from({ length: 10 }, (_, i) => currentYear + i);

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
};
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-4 flex items-center justify-center">
    <div v-if="bill" class="w-full max-w-2xl space-y-6">
      
      <div class="text-center space-y-2">
        <div class="flex items-center justify-center gap-2 text-2xl font-bold text-slate-900">
          <Shield class="h-6 w-6 text-emerald-600" />
          Secure Payment
        </div>
      </div>

      <form @submit="onSubmit" class="space-y-6">
        
        <Card class="border-0 shadow-lg">
          <CardHeader class="pb-4">
            <CardTitle class="flex items-center gap-2">
              <CreditCard class="h-5 w-5" />
              Card Information
            </CardTitle>
            <CardDescription>Enter your credit or debit card details</CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="firstName">First Name</Label>
                <div class="relative">
                  <User class="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                  <Input 
                    id="firstName" 
                    v-model="firstName"
                    placeholder="John" 
                    class="pl-10 h-12" 
                    :class="{ 'border-red-500': errors.firstName }"
                  />
                </div>
                <p v-if="errors.firstName" class="text-sm text-red-500">{{ errors.firstName }}</p>
              </div>
              <div class="space-y-2">
                <Label for="lastName">Last Name</Label>
                <Input 
                  id="lastName" 
                  v-model="lastName"
                  placeholder="Doe" 
                  class="h-12" 
                  :class="{ 'border-red-500': errors.lastName }"
                />
                <p v-if="errors.lastName" class="text-sm text-red-500">{{ errors.lastName }}</p>
              </div>
            </div>
            <div class="space-y-2">
              <Label for="cardNumber">Card Number</Label>
              <div class="relative">
                <CreditCard class="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                <Input
                  id="cardNumber"
                  placeholder="1234 5678 9012 3456"
                  v-model="cardNumber"
                  class="pl-10 h-12 font-mono"
                  maxlength="19"
                  :class="{ 'border-red-500': errors.cardNumber }"
                />
              </div>
              <p v-if="errors.cardNumber" class="text-sm text-red-500">{{ errors.cardNumber }}</p>
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div class="space-y-2">
                <Label for="month">Month</Label>
                <Select 
                  v-model="expiryMonth"
                >
                  <SelectTrigger 
                    id="month" 
                    class="h-12"
                    :class="{ 'border-red-500': errors.expiryMonth }"
                  >
                    <SelectValue placeholder="MM" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem v-for="i in 12" :key="i" :value="String(i).padStart(2, '0')">
                      {{ String(i).padStart(2, "0") }}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p v-if="errors.expiryMonth" class="text-sm text-red-500">{{ errors.expiryMonth }}</p>
              </div>
              <div class="space-y-2">
                <Label for="year">Year</Label>
                <Select 
                  v-model="expiryYear"
                >
                  <SelectTrigger 
                    id="year" 
                    class="h-12"
                    :class="{ 'border-red-500': errors.expiryYear }"
                  >
                    <SelectValue placeholder="YYYY" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem v-for="year in years" :key="year" :value="String(year)">
                      {{ year }}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p v-if="errors.expiryYear" class="text-sm text-red-500">{{ errors.expiryYear }}</p>
              </div>
              <div class="space-y-2">
                <Label for="cvc">CVC</Label>
                <Input 
                  id="cvc" 
                  v-model="cvc"
                  placeholder="123" 
                  class="h-12 font-mono" 
                  maxlength="4" 
                  :class="{ 'border-red-500': errors.cvc }"
                />
                <p v-if="errors.cvc" class="text-sm text-red-500">{{ errors.cvc }}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card class="border-0 shadow-lg bg-slate-50">
          <CardHeader class="pb-4">
            <CardTitle>Order Summary</CardTitle>
          </CardHeader>
          <CardContent class="space-y-3">
            <div class="flex justify-between text-sm">
              <span>Bill for {{ bill.billing_date }}</span>
              <span>{{ formatCurrency(bill.price - taxToPay) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span>Tax</span>
              <span>{{ formatCurrency(taxToPay) }}</span>
            </div>
            <div class="flex justify-between font-semibold text-lg">
              <span>Total</span>
              <span>{{ formatCurrency(bill.price) }}</span>
            </div>
          </CardContent>
        </Card>

        <Button
          type="submit"
          class="w-full h-14 text-lg font-semibold bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-700 hover:to-emerald-800 shadow-lg text-white"
          :disabled="isLoading"
        >
          <div v-if="isLoading" class="flex items-center gap-2">
            <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            Processing...
          </div>
          <div v-else class="flex items-center gap-2">
            <Lock class="h-5 w-5" />
            Pay {{ formatCurrency(bill.price) }}
          </div>
        </Button>

        <div class="text-center text-xs text-slate-500 space-y-1">
          <p>🔒 Your payment is secured.</p>
          <p>We accept Visa, Mastercard and American Express.</p>
        </div>
      </form>
    </div>
    
    <div v-if="!bill" class="text-center">Loading bill details...</div>
  </div>
  <Toaster />
</template>