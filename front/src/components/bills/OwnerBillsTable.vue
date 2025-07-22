<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import axios from 'axios'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shad/components/ui/table'
import { Button } from '@/shad/components/ui/button'
import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationListItem,
  PaginationNext,
  PaginationPrev,
} from '@/shad/components/ui/pagination'
import Input from '@/shad/components/ui/input/Input.vue';
import type { PriceListSearch, SearchBill } from './models'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import Spinner from '../Spinner.vue'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue

} from '@/shad/components/ui/select'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
}from '@/shad/components/ui/popover'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/shad/components/ui/dialog';
import { XCircleIcon } from '@heroicons/vue/24/solid';
import { useRouter } from 'vue-router';

const router = useRouter();
const bills = ref<SearchBill[]>([])

const searchQuery = ref({
  status: '',
  minPrice: '',
  maxPrice: '',
  billingDate: ''
});
const totalBills = ref(0);
const tempMonth = ref<string | undefined>();
const tempYear = ref<string | undefined>();

const months = Array.from({ length: 12 }, (_, i) => i + 1);

const currentYear = new Date().getFullYear();
const years = Array.from({ length: currentYear - 2000 }, (_, i) => currentYear - i);
const formatMonthName = (month: number) => {
  return new Date(0, month - 1).toLocaleString('en-EN', { month: 'long' });
};


watch([tempMonth, tempYear], ([month, year]) => {
  if (month && year) {
    const formattedMonth = month.padStart(2, '0');
    searchQuery.value.billingDate = `${year}-${formattedMonth}`;
    
    console.log(searchQuery.value.billingDate)
  } else {
    searchQuery.value.billingDate = '';
  }
});

const formattedPeriod = computed(() => {
  if (!searchQuery.value.billingDate) {
    return 'Choose month';
  }
  const temp = searchQuery.value.billingDate.split('-');
  const month = temp[1];
  const year = temp[0]
  const monthName = formatMonthName(+month);
  return `${monthName.charAt(0).toUpperCase() + monthName.slice(1)} ${year}`;
});


const pagination = ref({ page: 1, total: 0, perPage: 10, sortBy: '', sortOrder: '' })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))
const loading = ref(false);



watch(searchQuery, (newVal) => {
  Object.keys(newVal).forEach((key) => {
    
    if (newVal[key as keyof typeof newVal] === '' || newVal[key as keyof typeof newVal] === null) {
      delete newVal[key as keyof typeof newVal]
    }
  })
  if(newVal['minPrice'] === '0'){
    delete newVal['minPrice' as keyof typeof newVal]
  }
  if(newVal['maxPrice'] === '0'){
    delete newVal['maxPrice' as keyof typeof newVal]
  }
  
  console.log(newVal)
}, { deep: true })

const fetchBills = async () => {
  loading.value = true;

  const searchParams: Record<string, any> = {};
  
  if (searchQuery.value.status && searchQuery.value.status !== 'all') {
    searchParams.status = searchQuery.value.status;
  }
  if (searchQuery.value.minPrice !== null && +searchQuery.value.minPrice > 0) {
    searchParams.minPrice = +searchQuery.value.minPrice;
  }
  if (searchQuery.value.maxPrice !== null && +searchQuery.value.maxPrice > 0) {
    searchParams.maxPrice = +searchQuery.value.maxPrice;
  }
  if (searchQuery.value.billingDate) {
    searchParams.billingDate = searchQuery.value.billingDate;
  }
  const params = {
    page: pagination.value.page,
    pageSize: pagination.value.perPage,
    sortBy: pagination.value.sortBy,
    sortOrder: sortOrder.value[pagination.value.sortBy],

    search: JSON.stringify(searchParams)
  };

  try {
    const response = await axios.get('/api/bills/search', { params });

    if (response.data) {
      bills.value = response.data.bills;
      totalBills.value = response.data.total;
    }



  } catch (error: any) {
    bills.value = [];
    totalBills.value = 0;
  } finally {
    loading.value = false;
  }
};




function onPageChange(page: number) {
  pagination.value.page = page
  fetchBills()
}


function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : []
}

function formatDate(date: string): string {
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }
  return new Date(date).toLocaleString('sr-RS', options)
}
onMounted(() => {
  fetchBills()
})
watch(pagination, fetchBills, { deep: true });

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('sr-RS', { style: 'currency', currency: 'RSD' }).format(value);
};

const formatBillingPeriod = (billingDate: string) => {
  const [year, month] = billingDate.split('-');
  const monthName = new Date(parseInt(year), parseInt(month) - 1).toLocaleString('en-EN', { month: 'long' });
  return `${monthName.charAt(0).toUpperCase() + monthName.slice(1)} ${year}`;
};

const isPricelistDialogOpen = ref(false);
const selectedPricelist = ref<PriceListSearch | null>(null);
const showPricelistDetails = (pricelist: PriceListSearch) => {
  selectedPricelist.value = pricelist;
  isPricelistDialogOpen.value = true;
};

const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  billing_date: "desc",
  spent_power: "",
  price: "",
  issueDate: "",
})

function handleSort(field: string) {
  const temp = sortOrder.value[field]
  Object.keys(sortOrder.value).forEach((key) => (sortOrder.value[key] = ""))
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  pagination.value.sortBy = field
  fetchBills()
}

const clearPeriod = () => {
  searchQuery.value.billingDate = '';
  tempMonth.value = undefined;
  tempYear.value = undefined;
};


const onPay = (bill: number) => {
  router.push('/bills/pay/' + bill)
}

</script>

<template>

  <div class="p-7 flex flex-col bg-white shadow-lg">

    <div class="w-full text-center my-10 text-xl">
      Your bills
    </div>
    <div>
    <form class="flex flex-wrap items-center gap-8 border rounded-2xl p-10 mb-10" @submit.prevent="fetchBills">
      <FormField name="status">
        <FormItem class="items-center gap-6 pb-2">
          <FormLabel class="w-1/4 text-right">Status</FormLabel>
          <FormControl class="w-3/4">
            <Select v-model="searchQuery.status">
              <SelectTrigger>
                <SelectValue placeholder="Select status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem class="cursor-pointer" value="Paid">Paid</SelectItem>
                <SelectItem class="cursor-pointer" value="Delivered">Delivered</SelectItem>
                <SelectItem class="cursor-pointer" value="all">All</SelectItem>
              </SelectContent>
            </Select>
          </FormControl>
        </FormItem>
      </FormField>

      

      <FormField name="minprice">
        <FormItem class="pb-2">
          <FormLabel>Min price</FormLabel>
          <FormControl>
            <Input v-model="searchQuery.minPrice" type="number" step=".01" placeholder="Enter min price" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="maxprice">
        <FormItem class="pb-2">
          <FormLabel>Max price</FormLabel>
          <FormControl>
            <Input v-model="searchQuery.maxPrice" type="number" step=".01" placeholder="Enter max price" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="date">
        <FormItem class="items-center gap-4 pb-2">
          <FormLabel class="w-1/4 text-right">Month</FormLabel>
          <FormControl class="w-3/4">
            <Popover >
              <PopoverTrigger as-child>
                <Button variant="outline" class="w-full justify-between text-left font-normal">
                  <span>{{ formattedPeriod }}</span>
                  <button
                    class="rounded-full hover:bg-gray-200"
                    @click.stop.prevent="clearPeriod"
                  >
                    <XCircleIcon class="h-6 w-6 text-gray-500" />
                  </button>
                </Button>
              </PopoverTrigger>
              <PopoverContent class="w-auto p-4 flex flex-col gap-4">
                <div class="flex items-center gap-2">
                  <Select v-model="tempMonth" >
                    <SelectTrigger aria-label="Select month">
                      <SelectValue placeholder="Month" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem  class="cursor-pointer" v-for="month in months" :key="month" :value="month.toString()" >
                        {{ formatMonthName(month) }}
                      </SelectItem>
                    </SelectContent>
                  </Select>

                  <Select v-model="tempYear">
                    <SelectTrigger aria-label="Select year">
                      <SelectValue placeholder="Year" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem  class="cursor-pointer" v-for="year in years" :key="year" :value="year.toString()">
                        {{ year }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </PopoverContent>
            </Popover>
          </FormControl>
        </FormItem>
      </FormField>

      <Button type="submit" class="bg-indigo-500 text-white w-32 ml-10 hover:bg-gray-600 rounded-full mt-4">
        Search
      </Button>
    </form>
    </div>

    <Spinner v-if="loading" />
    <Table v-if="!loading" class="gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
        <TableHead @click="handleSort('billing_date')" :orientation="sortOrder.billing_date" class="cursor-pointer">
          Month
        </TableHead>
        <TableHead @click="handleSort('issue_date')" :orientation="sortOrder.issue_date" class="cursor-pointer">
          Issued At
        </TableHead>
        <TableHead>
          Owner
        </TableHead>
        <TableHead @click="handleSort('spent_power')" :orientation="sortOrder.spent_power" class="cursor-pointer">
          Power Spent (kWh)
        </TableHead>
        <TableHead @click="handleSort('price')" :orientation="sortOrder.price" class="cursor-pointer">
          Price
        </TableHead>
        <TableHead>
          Price List
        </TableHead>
        <TableHead>
          Status
        </TableHead>
        <TableHead>
          Pay Bill
        </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>

      <TableRow v-for="bill in bills" :key="bill.id" class="hover:bg-gray-50">
        <TableCell>
          {{ formatBillingPeriod(bill.billing_date) }}
        </TableCell>
        <TableCell>
          {{ formatDate(bill.issue_date) }}
        </TableCell>
        <TableCell>
          {{ bill.owner.username }}
        </TableCell>
        <TableCell clas="text-right">
          {{ bill.spent_power.toFixed(2) }}
        </TableCell>
         
        <TableCell>
          {{ formatCurrency(bill.price) }}
        </TableCell>

        <TableCell>
          <Button @click="showPricelistDetails(bill.pricelist)" variant="outline" size="sm">
            Show
          </Button>
        </TableCell>
        <TableCell>
          {{ bill.status }}
        </TableCell>
        <TableCell>
          <Button v-if="bill.status !== 'Paid'" size="sm" @click="onPay(bill.id)">
            Plati
          </Button>
        </TableCell>
      </TableRow>

    </TableBody>
  </Table>

    <div class="flex gap-20 pt-10">
      <Pagination v-slot="{ page }" :total="pagination.total" :sibling-count="1" show-edges
        :default-page="pagination.page" :items-per-page="pagination.perPage">
        <PaginationList v-slot="{ items }" class="flex items-center gap-1">
          <PaginationFirst @click="onPageChange(1)" :disabled="pagination.page === 1" />
          <PaginationPrev @click="onPageChange(pagination.page - 1)" :disabled="pagination.page === 1" />
          <template v-for="(item, index) in items">
            <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
              <Button class="w-10 h-10 p-0 hover:bg-indigo-300" :class="getButtonStyle(item.value === page)"
                :variant="item.value === page ? 'default' : 'outline'" @click="onPageChange(item.value)">
                {{ item.value }}
              </Button>
            </PaginationListItem>
            <PaginationEllipsis v-else :key="item.type" :index="index" />
          </template>

          <PaginationNext @click="onPageChange(pagination.page + 1)" :disabled="pagination.page === totalPages" />
          <PaginationLast @click="onPageChange(totalPages)" :disabled="pagination.page === totalPages" />
        </PaginationList>
      </Pagination>
      <div class="flex items-center gap-2">
        <span>Rows per page:</span>
        <Input v-model="pagination.perPage" type="number" class="w-20" min="1" placeholder="Rows per page" />
      </div>
    </div>
  </div>

  <Dialog v-model:open="isPricelistDialogOpen">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Price List Details</DialogTitle>
      </DialogHeader>
      
      <div v-if="selectedPricelist" class="grid gap-4 py-4">
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Valid from: </p>
          <p class="text-sm">{{ formatDate(selectedPricelist.valid_from) }}</p>
        </div>
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Billing power:</p>
          <p class="text-sm">{{ selectedPricelist.billing_power }} kWh</p>
        </div>
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Tax:</p>
          <p class="text-sm">{{ selectedPricelist.tax }}%</p>
        </div>
        
        <hr class="my-2">
        
        <p class="text-sm font-medium text-gray-500">Zone pricing (RSD/kWh):</p>
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Blue zone:</p>
          <p class="text-sm text-blue-600">{{ selectedPricelist.blue_zone }} RSD</p>
        </div>
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Green zone:</p>
          <p class="text-sm text-green-600">{{ selectedPricelist.green_zone }} RSD</p>
        </div>
        <div class="grid grid-cols-2 items-center gap-4">
          <p class="text-sm font-medium text-gray-500">Red zone:</p>
          <p class="text-sm text-red-600">{{ selectedPricelist.red_zone }} RSD</p>
        </div>
      </div>
    </DialogContent>
  </Dialog>
  <Toaster />
</template>
