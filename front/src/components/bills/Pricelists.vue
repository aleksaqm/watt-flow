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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shad/components/ui/dialog'
import Input from '@/shad/components/ui/input/Input.vue';
import type { Pricelist } from './pricelist'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';

const { toast } = useToast()
const props = defineProps({
  query: {
    type: Object,
    default: () => ({})
  },
  triggerSearch: {
    type: Number,
    default: 0
  }
})
const dialogKey = ref(0)


const pricelists = ref<Pricelist[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 10 })

const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  date: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))

async function fetchPricelists() {

  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: 'date',
      sortOrder: sortOrder.value['date'],
    }

    const response = await axios.post('/api/bills/query', props.query, { params: params })

    if (response.data && response.data.pricelists) {
      pricelists.value = response.data.pricelists.map((pricelist: any) => mapToPricelist(pricelist))
      pagination.value.total = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch pricelist:', error)
  }
}

function mapToPricelist(data: any): Pricelist {
  return {
    id: data.id,
    date: data.date,
    red: data.red,
    blue: data.blue,
    green: data.green,
    tax: data.tax,
    bill_power: data.bill_power,
    status: data.status,
  }
}


function onPageChange(page: number) {
  pagination.value.page = page
  fetchPricelists()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.date = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  fetchPricelists()
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

  <div class="p-7 flex flex-col bg-white shadow-lg gap-10 mt-10">
    <Dialog :key="dialogKey">
      <DialogTrigger>
        <Button variant="outline">Add new pricelist</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New pricelist</DialogTitle>
        </DialogHeader>
      </DialogContent>
    </Dialog>
    <Table class="gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('date')" :orientation="sortOrder.date">Valid from</TableHead>
          <TableHead>Bill. power</TableHead>
          <TableHead>Red</TableHead>
          <TableHead>Blue</TableHead>
          <TableHead>Green</TableHead>
          <TableHead>Tax</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="pricelist in pricelists  " :key="pricelist.id">
          <TableCell>{{ pricelist.date }}</TableCell>
          <TableCell>{{ pricelist.bill_power }}</TableCell>
          <TableCell>{{ pricelist.red }}</TableCell>
          <TableCell>{{ pricelist.blue }}</TableCell>
          <TableCell>{{ pricelist.green }}</TableCell>
          <TableCell>{{ pricelist.tax }}</TableCell>
          <TableCell>{{ pricelist.status }}</TableCell>
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
  <Toaster />
</template>
