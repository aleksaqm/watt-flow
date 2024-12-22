<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import Spinner from '../Spinner.vue'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shad/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shad/components/ui/dialog'
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
import { Button } from '@/shad/components/ui/button'
import NewUserForm from '@/components/auth/NewUserForm.vue'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';

const { toast } = useToast()

interface User {
  id: number
  username: string
  email: string
  role: string
}

const users = ref<User[]>([])
const dialogKey = ref(0)
const loading = ref(true)

const pagination = ref({ page: 1, total: 0, perPage: 6 })

async function fetchAdmins() {
  try {
    loading.value = true
    const response = await axios.get('/api/user/admins')
    users.value = response.data
    pagination.value.total = users.value.length
    loading.value = false
  } catch (error) {
    loading.value = false
    console.error('Failed to fetch users:', error)
  }
}

async function handleUserCreated() {
  await fetchAdmins()
  dialogKey.value++
  toast({
      title: 'You Added Admin Successfully',
      variant: 'default'
  })
}

onMounted(() => {
  fetchAdmins()
})

const paginatedUsers = computed(() => {
  const start = (pagination.value.page - 1) * pagination.value.perPage
  const end = start + pagination.value.perPage
  return users.value.slice(start, end)
})

function onPageChange(newPage: number) {
  pagination.value.page = newPage
}

function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : []
}

</script>


<template>
  <div class="w-1/2 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <Spinner v-if="loading" />
    <div v-if="!loading" class="flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-2xl">Admins</span>
      <Dialog :key="dialogKey">
        <DialogTrigger>
          <Button variant="outline" class="border-2">Add New Admin</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New Admin</DialogTitle>
          </DialogHeader>
          <NewUserForm url="/api/user/admin" role="Admin" @userCreated="handleUserCreated" />
        </DialogContent>
      </Dialog>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[100px]">ID</TableHead>
            <TableHead>Username</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Role</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="user in paginatedUsers" :key="user.id">
            <TableCell class="font-medium">{{ user.id }}</TableCell>
            <TableCell>{{ user.username }}</TableCell>
            <TableCell>{{ user.email }}</TableCell>
            <TableCell>{{ user.role }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Pagination
        v-slot="{ page }"
        :total="Math.ceil(pagination.total / pagination.perPage)"
        :default-page="pagination.page"
        :sibling-count="1"
        show-edges
      >
        <PaginationList v-slot="{ items }" class="flex items-center gap-1">
          <PaginationFirst @click="onPageChange(1)" :disabled="pagination.page === 1" />
          <PaginationPrev @click="onPageChange(pagination.page - 1)" :disabled="pagination.page === 1" />
          <template v-for="(item, index) in items">
            <PaginationListItem
              v-if="item.type === 'page'"
              :key="index"
              :value="item.value"
              as-child
            >
              <Button
                class="w-10 h-10 p-0 hover:bg-indigo-300"
                :class="getButtonStyle(item.value === page)"
                :variant="item.value === page ? 'default' : 'outline'"
                @click="onPageChange(item.value)"
              >
                {{ item.value }}
              </Button>
            </PaginationListItem>
            <PaginationEllipsis v-else :key="item.type" />
          </template>
          <PaginationNext
            @click="onPageChange(pagination.page + 1)"
            :disabled="pagination.page === Math.ceil(pagination.total / pagination.perPage)"
          />
          <PaginationLast
            @click="onPageChange(Math.ceil(pagination.total / pagination.perPage))"
            :disabled="pagination.page === Math.ceil(pagination.total / pagination.perPage)"
          />
        </PaginationList>
      </Pagination>
    </div>
  </div>
  <Toaster />
</template>
