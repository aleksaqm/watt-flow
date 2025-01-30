<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/shad/components/ui/alert-dialog'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import Spinner from '../Spinner.vue'

const { toast } = useToast()


const unsentBills = ref<string[]>([])

const loading = ref(false);

async function fetchUnsentBills() {

  try {
    loading.value = true;
    const response = await axios.get('/api/bills/unsent')

    if (response.data && response.data.data) {
      unsentBills.value = response.data.data
    }
    loading.value = false;
  } catch (error) {
    console.error('Failed to fetch unsent bills:', error)
    loading.value = false;
  }
}

onMounted(() => {
  fetchUnsentBills()
})

const sendBill = (month: string) => {

}

</script>

<template>

  <div class="p-7 flex flex-col bg-white shadow-lg gap-10 mt-10">
    <Spinner v-if="loading" />
    <Table v-if="!loading" class=" gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10">
      <TableHeader>
        <TableRow>
          <TableHead>Month</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="bill in unsentBills" :key="bill">
          <TableCell>{{ bill }}</TableCell>
          <TableCell>
            <AlertDialog>
              <AlertDialogTrigger><Button class="bg-indigo-500">Send</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction><Button @click="sendBill(bill)">Send</Button></AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </TableCell>

        </TableRow>
      </TableBody>
    </Table>
  </div>
  <Toaster />
</template>
