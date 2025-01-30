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
import type { Bill } from './models.ts'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import Spinner from '../Spinner.vue'

const bills = ref<Bill[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 10 })

const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  month: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))
const loading = ref(false);

async function fetchBills() {

  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: 'date',
      sortOrder: sortOrder.value['date'],
    }
    loading.value = true;

    const response = await axios.get('/api/bills/sent', { params: params })

    if (response.data && response.data.bills) {
      bills.value = response.data.bills.map((bill: any) => mapToBill(bill))
      pagination.value.total = response.data.total
    }
    loading.value = false;
  } catch (error) {
    console.error('Failed to fetch bills:', error)
    loading.value = false;
  }
}

function mapToBill(data: any): Bill {
  return {
    id: data.id,
    date: data.date,
    issue_date: data.issue_date,
    status: data.status
  }
}


function onPageChange(page: number) {
  pagination.value.page = page
  fetchBills()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.date = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  fetchBills()
}

function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : []
}

function formatDate(date: string): string {
  console.log(date)
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


</script>

<template>

  <div class="p-7 flex flex-col bg-white shadow-lg gap-10 mt-10">
    <Spinner v-if="loading" />
    <Table v-if="!loading" class=" gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('date')" :orientation="sortOrder.date">Month</TableHead>
          <TableHead>Issue date</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="bill in bills  " :key="bill.id">
          <TableCell>{{ bill.date }}</TableCell>
          <TableCell>{{ formatDate(bill.issue_date.toString()) }}</TableCell>
          <TableCell>{{ bill.status }}</TableCell>
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
