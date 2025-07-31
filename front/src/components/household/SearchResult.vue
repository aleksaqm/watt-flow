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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shad/components/ui/dialog'
import Input from '@/shad/components/ui/input/Input.vue';
import router from '@/router'
import type { Household } from './household'
import { useUserStore } from '@/stores/user'
import OwnershipRequestForm from '@/components/household/OwnershipRequestForm.vue'
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
  },
  mode: {
    type: String,
    default: 'search' // 'search' | 'my-households'
  }
})
const dialogKey = ref(0)


// Watch for changes in triggerSearch to run fetch
watch(() => props.triggerSearch, () => {
  pagination.value.page = 1
  fetchHouseholds()
})
const households = ref<Household[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 10 })

const sortBy = ref("city")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  city: "",
  street: "",
  number: "",
  status: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))

async function fetchHouseholds() {

  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value[sortBy.value],
    }
    console.log(props.query)

    const response = await axios.post('/api/household/query', props.query, { params: params })

    if (response.data && response.data.households) {
      households.value = response.data.households.map((household: any) => mapToHousehold(household))

      pagination.value.total = response.data.total
    }
  } catch (error) {
    console.error('Failed to fetch households:', error)
  }
}

function mapToHousehold(data: any): Household {
  if (data.status == 0) {
    data.status = "Pending"
  }
  return {
    id: data.id,
    city: data.city,
    street: data.street,
    number: data.number,
    status: data.status.toString(),
    cadastral_number: data.cadastral_number,
    floor: data.floor,
    suite: data.suite,
  }
}


function onPageChange(page: number) {
  pagination.value.page = page
  fetchHouseholds()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.city = ""
  sortOrder.value.street = ""
  sortOrder.value.number = ""
  sortOrder.value.status = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  console.log(sortOrder.value)
  sortBy.value = field
  fetchHouseholds()
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

const userStore = useUserStore();
const isAdmin = ref<boolean>(true)
onMounted(()=>{
  if (userStore.role === "Regular"){
    isAdmin.value = false;
  }
})

function viewHousehold(id: number) {
  if (isAdmin.value) {
    router.push({ name: "household", params: { id: id } })
  } else if (props.mode === 'my-households') {
    router.push({ name: "my-household", params: { id: id } })
  } else {
    console.log("Nothing");
  }
}

function handleRequestSent(){
  toast({
      title: 'Request Sent Successfully',
      variant: 'default'
  })
  dialogKey.value ++;
}

</script>

<template>

  <div class="p-7 flex flex-col bg-white shadow-lg">

    <Table class="gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('city')" :orientation="sortOrder.city">City</TableHead>
          <TableHead @click="onSortChange('street')" :orientation="sortOrder.street">Street</TableHead>
          <TableHead @click="onSortChange('number')" :orientation="sortOrder.number">Number</TableHead>
          <TableHead @click="onSortChange('status')" :orientation="sortOrder.status">Status</TableHead>
          <TableHead>Floor</TableHead>
          <TableHead>Suite</TableHead>
          <TableHead>C-Number</TableHead>
          <TableHead v-if="!isAdmin && props.mode === 'search'">Ownership</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="household in households" :key="household.id" @click="viewHousehold(household.id)">
          <TableCell>{{ household.city }}</TableCell>
          <TableCell>{{ household.street }}</TableCell>
          <TableCell>{{ household.number }}</TableCell>
          <TableCell>{{ household.status }}</TableCell>
          <TableCell>{{ household.floor }}</TableCell>
          <TableCell>{{ household.suite }}</TableCell>
          <TableCell>{{ household.cadastral_number }}</TableCell>
          <TableCell v-if="!isAdmin && props.mode === 'search'">
            <Dialog :key="dialogKey">
              <DialogTrigger>
                <Button class="bg-indigo-500 hover:bg-gray-600">Claim</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Prove your ownership</DialogTitle>
                </DialogHeader>
                <OwnershipRequestForm :household="household.id" @requestSent="handleRequestSent"/>
              </DialogContent>
            </Dialog>
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
