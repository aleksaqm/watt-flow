
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Spinner from './Spinner.vue'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../shad/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../shad/components/ui/dialog'
import { Button } from '../shad/components/ui/button'
import NewUserForm from './NewUserForm.vue'

interface User {
  id: number
  username: string
  email: string
  role: string
}

const users = ref<User[]>([])
// const dialogOpen = ref(false)
const dialogKey = ref(0)
const loading = ref(true)


async function fetchAdmins() {
  try {
    loading.value=true
    const response = await axios.get('/api/user/admins')
    loading.value=false
    users.value = response.data
  } catch (error) {
    loading.value=false
    console.error('Failed to fetch users:', error)
  }
}

async function handleUserCreated() {
  fetchAdmins()
  dialogKey.value++
}

onMounted(async () => {
  fetchAdmins()
})
</script>

<template>
  <div class="w-1/2 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <Spinner v-if="loading"/>
    <div v-if="!loading" class="flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-2xl">Admins</span>
      <Dialog :key="dialogKey">
            <DialogTrigger>
                <Button variant="outline" class="border-2">
                    Add New Admin
                </Button>
            </DialogTrigger>
            
            <DialogContent ref="dialogRef">
                <DialogHeader>
                    <DialogTitle>Add New Admin</DialogTitle>
                </DialogHeader>
                <NewUserForm url="/api/user/admin" role="Admin" @userCreated="handleUserCreated"></NewUserForm>
            </DialogContent>
        </Dialog>
      <!-- <Button @click="showDialog = true" class="bg-gray-500 text-white rounded-xl px-4 py-2 mb-4">Add New Admin</Button> -->
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
          <TableRow v-for="user in users" :key="user.id">
            <TableCell class="font-medium">{{ user.id }}</TableCell>
            <TableCell>{{ user.username }}</TableCell>
            <TableCell>{{ user.email }}</TableCell>
            <TableCell>{{ user.role }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>