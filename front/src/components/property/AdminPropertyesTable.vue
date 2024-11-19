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
  Form,
  FormControl,
  FormDescription,
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
import Input from '@/shad/components/ui/input/Input.vue'
import { toTypedSchema } from '@vee-validate/zod'
import { h, watch } from 'vue'
import * as z from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose
} from '@/shad/components/ui/dialog'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { useToast } from '../../shad/components/ui/toast/use-toast'

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


const { toast } = useToast()

const properties = ref<Property[]>([])

const currentPropertyId = ref<number | null>(null); 
const pagination = ref({ page: 1, total: 0, perPage: 5 })
const searchQuery = ref<{ city?: string; street?: string; number?: string; floors?: number }>({})
const sortBy = ref("created_at")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  city: "",
  street: "",
  number: "",
  status: "",
  created_at: "desc",
  floors: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))


async function fetchProperties() {
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value[sortBy.value],
      search: JSON.stringify(searchQuery.value),
    }

    const response = await axios.get('/api/property/query', { params })

    if (response.data && response.data.properties) {
      properties.value = response.data.properties.map((property: any) => mapToProperty(property))
      pagination.value.total = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch properties:', error)
  }
}

function mapToProperty(data: any): Property {
  return {
    id: data.id,
    city: data.address.city,
    street: data.address.street,
    number: data.address.number,
    status: data.status === 0 ? "Pending" : (data.status === 1 ? "Declined" : (data.status === 2 ? "Accepted" : data.status.toString())),
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
  const temp = sortOrder.value[field]
  Object.keys(sortOrder.value).forEach((key) => (sortOrder.value[key] = ""))
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
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
  return new Date(date).toLocaleString('en-US', options)
}

async function handleAccept(id: number) {
  try {
    const response = await axios.put(`/api/property/` + id +`/accept`)
    console.log(`Property accepteded successfully`, response.data)
    fetchProperties()
    toast({
      title: 'Property Accepted',
      description: `Property was accepted successfully.`,
      variant: "default",
    });
  } catch (error) {
    console.error(`Failed to accept property with ID ${id}:`, error)
  }
}

const formSchema = toTypedSchema(z.object({
  declineReason: z.string().min(2).max(50),
}))

async function handleDecline(values: any) {
  try {
    if (!currentPropertyId.value) {
      throw new Error("No property ID found for declining.");
    }
    console.log(`Declining property with ID: ${currentPropertyId.value}`);
    console.log(`Reason: ${values.declineReason}`);

    await axios.put(`/api/property/${currentPropertyId.value}/decline`, {
      message: values.declineReason,
    });

    fetchProperties();
    toast({
      title: 'Property Declined',
      description: `Property was declined successfully.`,
      variant: "default",
    });
  } catch (error) {
    console.error(`Failed to decline property with ID ${currentPropertyId.value}:`, error);
    toast({
      title: 'Error:',
      description: "Error declining property request.",
      variant: "destructive",
    });
  } finally {
    currentPropertyId.value = null; 
  }
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
        Search Property
      </div>

      <form class="flex flex-wrap gap-5 items-center border rounded-2xl p-10 mb-10" @submit.prevent="fetchProperties">
        <FormField name="city">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">City:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.city" placeholder="Enter city" />
            </FormControl>
          </FormItem>
        </FormField>
        <FormField name="street">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Street:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.street" placeholder="Enter street" />
            </FormControl>
          </FormItem>
        </FormField>
        <FormField name="number">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Number:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.number" placeholder="Enter number" />
            </FormControl>
          </FormItem>
        </FormField>
        <FormField name="floors">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Floors:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="number" v-model="searchQuery.floors" placeholder="Enter floors" />
            </FormControl>
          </FormItem>
        </FormField>
        <Button type="submit" class="bg-indigo-500 text-white w-32 ml-10 hover:bg-gray-600 rounded-full py-2">Search</Button>
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
            <TableHead>Actions</TableHead>
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
          <TableCell>
            <Button class="bg-indigo-500 text-white mr-2 hover:bg-indigo-300" @click="handleAccept(property.id)" v-if="property.status === 'Pending'">Accept</Button>

            <Form v-slot="{ handleSubmit }" as="" :validation-schema="formSchema">
                <Dialog>
                    <DialogTrigger as-child>
                        <Button class="bg-red-500 text-white" v-if="property.status === 'Pending'" @click="currentPropertyId = property.id">Decline</Button>
                    </DialogTrigger>
                    <DialogContent class="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>Decline reason</DialogTitle>
                            <DialogDescription>
                                Reason for declining property request
                            </DialogDescription>
                        </DialogHeader>
                        <form id="dialogForm" @submit="handleSubmit($event, handleDecline)">
                            <FormField v-slot="{ componentField }" name="declineReason">
                                <FormItem>
                                    <FormLabel>Decline reason</FormLabel>
                                    <FormControl>
                                        <Input type="text" placeholder="Reason" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                        </form>
                        <DialogFooter>
                            <Button type="submit" form="dialogForm">
                                Submit
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </Form>
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
  <Toaster />
</template>
