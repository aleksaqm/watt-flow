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

import Input from '@/shad/components/ui/input/Input.vue';
import { useUserStore } from '@/stores/user';
import { getUserIdFromToken } from '@/utils/jwtDecoder'
import { z, type string } from 'zod'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { useToast } from '../../shad/components/ui/toast/use-toast'
import { toTypedSchema } from '@vee-validate/zod'
import Spinner from '../Spinner.vue'
import ImageDocumentsDisplay from '@/components/household/ImageDocumentsDisplay.vue'


interface Meeting {
  id: number
  start_time: string
  duration: number
  clerk: string
}

const { toast } = useToast()

const meetings = ref<Meeting[]>([])
const isLoading = ref<boolean>(false)

const pagination = ref({ page: 1, total: 0, perPage: 5 })
const searchQuery = ref<{ clerk?: string;}>({})

const sortBy = ref("start_time")
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
    clerk: "",
    start_time: ""
})
const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.perPage))

function fetchMeetingsForm(event: Event){
  const submitEvent = event as SubmitEvent;
  submitEvent.preventDefault();
  fetchMeetings();
  (submitEvent.submitter as HTMLButtonElement).blur();
}

async function fetchMeetings() {
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value[sortBy.value],
      search: JSON.stringify(searchQuery.value),
    }
    console.log(params)

    isLoading.value = true
    const userStore = useUserStore()
    const response = await axios.get("/api/user/meetings/" + userStore.id, { params })
    isLoading.value = false
    console.log(response)
    if (response.data) {
        meetings.value = response.data.meetings.map((meeting: any) => mapToMeeting(meeting))
        pagination.value.total = response.data.total
    }
    console.log(response.data.meetings)
  } catch (error) {
    isLoading.value = false
    console.error('Failed to fetch meetings:', error)
  }
}

function mapToMeeting(data: any): Meeting {
  return {
    id: data.id,
    start_time: data.start_time,
    duration: data.duration,
    clerk: data.clerk
  }
}

onMounted(()=>{
    fetchMeetings()
})


function onPageChange(page: number) {
  pagination.value.page = page
  fetchMeetings()
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field]
  sortOrder.value.clerk = ""
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc"
  console.log(sortOrder.value)
  sortBy.value = field
  fetchMeetings()
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
  
  <div class="mx-auto w-10/12 p-7 flex flex-col bg-white shadow-lg">
    <div>
      <div class="w-full text-center my-10 text-xl">
        My meetings
      </div>
      <form class="w-full flex flex-wrap gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10 text-center justify-center"
        @submit.prevent="fetchMeetingsForm">
        <FormField name="clerk" v-slot="{ field }">
          <FormItem class="flex items-center gap-5">
            <FormLabel class="w-1/4 text-right">Clerk:</FormLabel>
            <FormControl class="w-3/4">
              <Input type="text" v-model="searchQuery.clerk" placeholder="Enter clerk" />
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
    <Table class="w-full border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead @click="onSortChange('clerk')" :orientation="sortOrder.clerk">Clerk</TableHead>
          <TableHead @click="onSortChange('start_time')" :orientation="sortOrder.start_time">StartTime</TableHead>
          <TableHead>Duration</TableHead>
        </TableRow>
      </TableHeader>
      <Spinner  v-if="isLoading"/>
      <TableBody v-if="!isLoading">
        <TableRow v-for="meeting in meetings" :key="meeting.id">
          <TableCell>{{ meeting.clerk }}</TableCell>
          <TableCell>{{ formatDate(meeting.start_time) }}</TableCell>
          <TableCell>{{ meeting.duration }}</TableCell>
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
            @input="fetchMeetings"
            placeholder="Rows per page"
          />
        </div>
    </div>
  </div>
  <Toaster />
</template>
