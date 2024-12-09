<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import axios from 'axios'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shad/components/ui/table'
import { Button } from '@/shad/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
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
import { useUserStore } from '@/stores/user';
import { getUserIdFromToken } from '@/utils/jwtDecoder'
import type { string } from 'zod'


interface OwnershipRequest {
  id: number
  images: string[]
  documents: string[]
  created_at: string
  closed_at: string
  city: string
  street: string
  number: string
  floor: number
  suite: string
  status: string
  username: string
}

const requests = ref<OwnershipRequest[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 5 })
const searchQuery = ref<{ city?: string; street?: string; number?: string; floor?: number, suite?: string}>({})

const sortBy = ref("created_at")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  city: "",
  street: "",
  number: "",
  status: "",
  created_at: "",
  closed_at: "",
  floor: "",
  suite: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))

function fetchPropertiesForm(event: Event){
  const submitEvent = event as SubmitEvent;
  submitEvent.preventDefault();
  fetchProperties();
  (submitEvent.submitter as HTMLButtonElement).blur();
}

async function fetchProperties() {
  try {
    const userStore = useUserStore()
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value[sortBy.value],
      search: JSON.stringify(searchQuery.value),
    }
    console.log(params)

    const response = await axios.get('/api/property/query', { params })
    console.log(response)

    if (response.data) {
        requests.value = response.data.requests.map((request: any) => mapToRequest(request))
        pagination.value.total = response.data.total
    }
    console.log(response.data.requests)
  } catch (error) {
    console.error('Failed to fetch properties:', error)
  }
}

function mapToRequest(data: any): OwnershipRequest {
  return {
    id: data.id,
    city: data.city,
    street: data.street,
    number: data.number,
    floor: data.floor,
    suite: data.suite,
    status: data.status,
    created_at: data.created_at,
    closed_at: data.closed_at,
    username: data.username,
    images: data.images,
    documents: data.documents,
  }
}

onMounted(fetchProperties)

function onPageChange(page: number) {
  pagination.value.page = page
  fetchProperties()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.city = ""
  sortOrder.value.street = ""
  sortOrder.value.number = ""
  sortOrder.value.status = ""
  sortOrder.value.created_at = ""
  sortOrder.value.closed_at = ""
  sortOrder.value.floor = ""
  sortOrder.value.suite = ""
//   sortOrder.value.username = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  console.log(sortOrder.value)
  sortBy.value = field
  fetchProperties()
}

function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : []
}

function formatDate(date: string): string {
    if (date === "0001-01-01 00:00:00 +0000 UTC"){
        return "/"
    }
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }
  const dateObj = new Date(date)
  return dateObj.toLocaleString('en-US', options)
}

watch(searchQuery, (newVal) => {
  Object.keys(newVal).forEach((key) => {
    if (newVal[key as keyof typeof newVal] === '' || newVal[key as keyof typeof newVal] === null) {
      delete newVal[key as keyof typeof newVal]
    }
  })
}, { deep: true })

</script>

<template>
  
  <div class="p-7 flex flex-col bg-white shadow-lg">
    <div>
      <div class="w-full text-center my-10 text-xl">
        Your Ownership Requests
      </div>
      <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10"
        @submit.prevent="fetchPropertiesForm">

        <FormField name="city" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">City:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.city" placeholder="Enter city" />
            </FormControl>
          </FormItem>
        </FormField>

        <FormField name="street" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Street:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.street" placeholder="Enter street" />
            </FormControl>
          </FormItem>
        </FormField>

        <FormField name="number" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Number:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.number" placeholder="Enter number" />
            </FormControl>
          </FormItem>
        </FormField>

        <FormField name="floor" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Floor:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="number" v-model="searchQuery.floor" placeholder="Enter floor number" />
            </FormControl>
          </FormItem>
        </FormField>

        <FormField name="suite" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Suite:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="string" v-model="searchQuery.suite" placeholder="Enter suite" />
            </FormControl>
          </FormItem>
        </FormField>

        <Button
          type="submit"
          @click="(event: MouseEvent) => (event.target as HTMLButtonElement).blur()"
          class="bg-indigo-500 text-white w-32 ml-10 hover:bg-gray-600 rounded-full py-2"
        >
          Search
        </Button>
      </form>

    </div>
    <Table class="gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('city')" :orientation="sortOrder.city">City</TableHead>
          <TableHead @click="onSortChange('street')" :orientation="sortOrder.street">Street</TableHead>
          <TableHead @click="onSortChange('number')" :orientation="sortOrder.number">Number</TableHead>
          <TableHead @click="onSortChange('floor')" :orientation="sortOrder.floor">Floor</TableHead>
          <TableHead @click="onSortChange('suite')" :orientation="sortOrder.suite">Suite</TableHead>
          <TableHead @click="onSortChange('status')" :orientation="sortOrder.status">Status</TableHead>
          <TableHead @click="onSortChange('created_at')" :orientation="sortOrder.created_at">Creation Time</TableHead>
          <TableHead @click="onSortChange('closed_at')" :orientation="sortOrder.closed_at">Resolved Time</TableHead>
          <!-- <TableHead @click="onSortChange('username')" :orientation="sortOrder.username">Username</TableHead> -->
          <!-- <TableHead>Households</TableHead> -->
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="property in requests" :key="property.id">
          <TableCell>{{ property.city }}</TableCell>
          <TableCell>{{ property.street }}</TableCell>
          <TableCell>{{ property.number }}</TableCell>
          <TableCell>{{ property.floor }}</TableCell>
          <TableCell>{{ property.suite }}</TableCell>
          <TableCell>{{ property.status }}</TableCell>
          <TableCell>{{ formatDate(property.created_at) }}</TableCell>
          <TableCell>{{ formatDate(property.closed_at) }}</TableCell>
          <!-- <TableCell>{{ property.username }}</TableCell> -->
          <!-- <TableCell>{{ property.households }}</TableCell> -->
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
          <Input
            v-model="pagination.perPage"
            type="number"
            class="w-20"
            min="1"
            @input="fetchProperties"
            placeholder="Rows per page"
          />
        </div>
    </div>
  </div>
</template>
