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
import Input from '@/shad/components/ui/input/Input.vue';
import { useUserStore } from '@/stores/user';
import { getUserIdFromToken } from '@/utils/jwtDecoder'
import { z, type string } from 'zod'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { useToast } from '../../shad/components/ui/toast/use-toast'
import { toTypedSchema } from '@vee-validate/zod'
import Spinner from '../Spinner.vue'
import ImageDocumentsDisplay from '@/components/household/ImageDocumentsDisplay.vue'


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
  username: string
  status: string
}

const { toast } = useToast()

const requests = ref<OwnershipRequest[]>([])
const isAdmin = ref<boolean>(true)
const isLoading = ref<boolean>(false)
const currentReqyestId = ref<number | null>(null); 
const isDialogOpen = ref<boolean>(false);
const dialogsOpen = ref<boolean[]>([])

const pagination = ref({ page: 1, total: 0, perPage: 5 })
const searchQuery = ref<{ city?: string; street?: string; number?: string; floor?: number, suite?: string}>({})

const sortBy = ref("created_at")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  city: "",
  street: "",
  number: "",
  username: "",
  created_at: "",
  closed_at: "",
  floor: "",
  suite: "",
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))

function fetchRequestsForm(event: Event){
  const submitEvent = event as SubmitEvent;
  submitEvent.preventDefault();
  fetchRequests();
  (submitEvent.submitter as HTMLButtonElement).blur();
}

async function fetchRequests() {
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

    var requestUrl : string
    if(isAdmin.value){
        requestUrl = "/api/ownership/pending"
    }else{
        requestUrl = "/api/ownership/requests/"+ userStore.id
    }
    isLoading.value = true
    const response = await axios.get(requestUrl, { params })
    isLoading.value = false
    console.log(response)
    if (response.data) {
        requests.value = response.data.requests.map((request: any) => mapToRequest(request))
        pagination.value.total = response.data.total
    }
    console.log(response.data.requests)
  } catch (error) {
    isLoading.value = false
    console.error('Failed to fetch requests:', error)
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

onMounted(()=>{
    const userStore = useUserStore()
    if (userStore.role === "Regular"){
        isAdmin.value = false
    }
    fetchRequests()
})


function onPageChange(page: number) {
  pagination.value.page = page
  fetchRequests()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.city = ""
  sortOrder.value.street = ""
  sortOrder.value.number = ""
  sortOrder.value.username = ""
  sortOrder.value.created_at = ""
  sortOrder.value.closed_at = ""
  sortOrder.value.floor = ""
  sortOrder.value.suite = ""
//   sortOrder.value.status = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  console.log(sortOrder.value)
  sortBy.value = field
  fetchRequests()
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

async function handleAccept(id: number){
  try {
    isLoading.value = true
    const response = await axios.patch(`/api/ownership/accept/` + id)
    console.log(`Request accepteded successfully`, response.data)
    fetchRequests()
    isLoading.value = false
    toast({
      title: 'Request Accepted',
      description: `Request was accepted successfully.`,
      variant: "default",
    });
  } catch (error) {
    isLoading.value = false
    console.error(`Failed to accept request with ID ${id}:`, error)
  }
}
const formSchema = toTypedSchema(z.object({
  declineReason: z.string().min(2).max(50),
}))

async function handleDecline(values: any) {
  try {
    if (!currentReqyestId.value) {
      throw new Error("No request ID found for declining.");
    }
    console.log(`Declining request with ID: ${currentReqyestId.value}`);
    console.log(`Reason: ${values.declineReason}`);
    isLoading.value = true

    await axios.put(`/api/ownership/decline/${currentReqyestId.value}`, {
      message: values.declineReason,
    });

    fetchRequests();
    isLoading.value = false
    toast({
      title: 'Request Declined',
      description: `Request was declined successfully.`,
      variant: "default",
    });
  } catch (error) {
    isLoading.value = false
    console.error(`Failed to decline request with ID ${currentReqyestId.value}:`, error);
    toast({
      title: 'Error:',
      description: "Error declining  request.",
      variant: "destructive",
    });
  } finally {
    currentReqyestId.value = null; 
  }
}

function openDialog(){
  isDialogOpen.value = true;
}

function closeDialog(){
  isDialogOpen.value = false;
}

function onDialogUpdate(value: any) {
  isDialogOpen.value = value; 
}

</script>

<template>
  
  <div class="p-7 flex flex-col bg-white shadow-lg">
    <div>
      <div class="w-full text-center my-10 text-xl">
        Ownership Requests
      </div>
      <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10"
        @submit.prevent="fetchRequestsForm">
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
          <TableHead @click="onSortChange('username')" :orientation="sortOrder.username">username</TableHead>
          <TableHead @click="onSortChange('created_at')" :orientation="sortOrder.created_at">Creation Time</TableHead>
          <TableHead @click="onSortChange('closed_at')" :orientation="sortOrder.closed_at">Resolved Time</TableHead>
          <TableHead v-if="isAdmin">Documentation</TableHead>
          <TableHead v-if="isAdmin">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <Spinner  v-if="isLoading"/>
      
      <TableBody v-if="!isLoading">
        <TableRow v-for="request in requests" :key="request.id">
          <TableCell @click="openDialog">{{ request.city }}</TableCell>
          <TableCell @click="openDialog">{{ request.street }}</TableCell>
          <TableCell @click="openDialog">{{ request.number }}</TableCell>
          <TableCell @click="openDialog">{{ request.floor }}</TableCell>
          <TableCell @click="openDialog">{{ request.suite }}</TableCell>
          <TableCell @click="openDialog">{{ request.username }}</TableCell>
          <TableCell @click="openDialog">{{ formatDate(request.created_at) }}</TableCell>
          <TableCell @click="openDialog">{{ formatDate(request.closed_at) }}</TableCell>
          <TableCell v-if="isAdmin">
            <Dialog>
              <DialogTrigger>
                <Button class="bg-gray-600 text-white">Details</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogTitle>Photos & Documents for proof</DialogTitle>
                <DialogDescription>
                  Analize documentation for ownership
                </DialogDescription>
                <div class="flex justify-center items-center w-full h-full">
                  <ImageDocumentsDisplay :images="request.images" :documents="request.documents" />
                </div>
              </DialogContent>
            </Dialog>
          </TableCell>
            
          <TableCell v-if="isAdmin">
            <Button class="bg-indigo-700 text-white mr-2 hover:bg-indigo-300" @click="handleAccept(request.id)">Accept</Button>
            <Form v-slot="{ handleSubmit }" as="" :validation-schema="formSchema">
                <Dialog>
                    <DialogTrigger as-child>
                        <Button class="bg-gray-500 text-white" @click="currentReqyestId = request.id">Decline</Button>
                    </DialogTrigger>
                    <DialogContent class="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>Decline reason</DialogTitle>
                            <DialogDescription>
                                Reason for declining ownership request
                            </DialogDescription>
                        </DialogHeader>
                        <Spinner  v-if="isLoading"/>
                        <form v-if="!isLoading" id="dialogForm" @submit="handleSubmit($event, handleDecline)">
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
            @input="fetchRequests"
            placeholder="Rows per page"
          />
        </div>
    </div>
  </div>
  <Toaster />
</template>
