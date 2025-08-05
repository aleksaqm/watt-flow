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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/shad/components/ui/alert-dialog'
import Input from '@/shad/components/ui/input/Input.vue';
import type { Pricelist } from './models'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import NewPricelistForm from './NewPricelistForm.vue'
import Spinner from '../Spinner.vue'

const { toast } = useToast()
const dialogKey = ref(0)


const pricelists = ref<Pricelist[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 10 })

const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  valid_from: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))
const loading = ref(false);

async function fetchPricelists() {

  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: 'valid_from',
      sortOrder: sortOrder.value['valid_from'],
    }

    const response = await axios.get('/api/pricelist/query', { params: params })

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
    date: data.valid_from,
    red: data.red_zone,
    blue: data.blue_zone,
    green: data.green_zone,
    tax: data.tax,
    bill_power: data.billing_power,
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
  console.log(date)
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }
  return new Date(date).toLocaleString('sr-RS', options)
}

const pricelistCreated = () => {
  toast({
    title: 'Operation Successful',
    description: 'You successfully added new pricelist!',
    variant: 'default'
  })
  dialogKey.value++;
  fetchPricelists()
}

const deletePricelist = async (id: number) => {
  try {
    loading.value = true;
    const response = await axios.delete('/api/pricelist/' + id)
    console.log(response)
    toast({
      title: 'Operation succeded!',
      description: "Pricelist deleted successfully!",
      variant: 'default'
    })
    loading.value = false;
    fetchPricelists()
  } catch (error: any) {
    const errorMessage = 'Try again later.'
    console.error('Error:', error)
    loading.value = false;
    toast({
      title: 'Pricelist deletion failed',
      description: errorMessage,
      variant: 'destructive'
    })
  }

}
onMounted(() => {
  fetchPricelists()
})

const isDeleteDisabled = (date: Date) => {
  return new Date(date) <= new Date();
}

</script>

<template>

  <div class="mx-auto w-full p-7 flex flex-col bg-white shadow-lg gap-10 mt-10">
    <Dialog :key="dialogKey">
      <DialogTrigger>
        <Button variant="outline">Add new pricelist</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New pricelist</DialogTitle>
        </DialogHeader>
        <NewPricelistForm @pricelist-created="pricelistCreated"></NewPricelistForm>
      </DialogContent>

    </Dialog>
    <Spinner v-if="loading" />
    <Table v-if="!loading" class="w-full border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('date')" :orientation="sortOrder.date">Valid from</TableHead>
          <TableHead>Bill. power</TableHead>
          <TableHead>Red</TableHead>
          <TableHead>Blue</TableHead>
          <TableHead>Green</TableHead>
          <TableHead>Tax</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="pricelist in pricelists  " :key="pricelist.id">
          <TableCell>{{ formatDate(pricelist.date.toString()) }}</TableCell>
          <TableCell>{{ pricelist.bill_power }}</TableCell>
          <TableCell>{{ pricelist.red }}</TableCell>
          <TableCell>{{ pricelist.blue }}</TableCell>
          <TableCell>{{ pricelist.green }}</TableCell>
          <TableCell>{{ pricelist.tax }}</TableCell>
          <TableCell>
            <AlertDialog>
              <AlertDialogTrigger :disabled="isDeleteDisabled(pricelist.date)"><Button
                  :disabled="isDeleteDisabled(pricelist.date)" class="bg-red-300">Delete</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction><Button @click="deletePricelist(pricelist.id)">Delete</Button></AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>

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
  <Toaster />
</template>
