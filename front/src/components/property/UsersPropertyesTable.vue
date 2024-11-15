<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
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

interface Property {
  id: number
  city: string
  street: string
  number: string
  status: string
  created: string
  floors: number
  households: number
}

const properties = ref<Property[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 2 })
const searchQuery = ref<{ city?: string; street?: string; number?: string; floors?: number }>({})

const sortBy = ref("city")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  city: "",
  street: "",
  number: "",
  status: "",
  created_at: "",
  floors: "",
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

    if (response.data && response.data.properties) {
      properties.value = response.data.properties.map((property: any) => mapToProperty(property))

      pagination.value.total = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch properties:', error)
  }
}

function mapToProperty(data: any): Property {
  if (data.status == 0) {
    data.status = "Pending"
  }
  return {
    id: data.id,
    city: data.address.city,
    street: data.address.street,
    number: data.address.number,
    status: data.status.toString(),
    created: data.created_at,
    floors: data.floors,
    households: data.household.length,
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
  sortOrder.value.floors = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  console.log(sortOrder.value)
  sortBy.value = field
  fetchProperties()
}

function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : []
}

function formatDate(date: string): string {
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

</script>

<template>
  
  <div class="p-7 flex flex-col bg-white shadow-lg">
    <div>
      <div class="w-full text-center my-10 text-xl">
        Search Property
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

        <FormField name="floors" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Floors:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="number" v-model="searchQuery.floors" placeholder="Enter floors" />
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
          <TableHead @click="onSortChange('status')" :orientation="sortOrder.status">Status</TableHead>
          <TableHead @click="onSortChange('created_at')" :orientation="sortOrder.created_at">Creation Time</TableHead>
          <TableHead @click="onSortChange('floors')" :orientation="sortOrder.floors">Floors</TableHead>
          <TableHead>Households</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="property in properties" :key="property.id">
          <TableCell>{{ property.city }}</TableCell>
          <TableCell>{{ property.street }}</TableCell>
          <TableCell>{{ property.number }}</TableCell>
          <TableCell>{{ property.status }}</TableCell>
          <TableCell>{{ formatDate(property.created) }}</TableCell>
          <TableCell>{{ property.floors }}</TableCell>
          <TableCell>{{ property.households }}</TableCell>
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
