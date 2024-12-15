<script setup lang="ts">
import { Button } from '@/shad/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import { Input } from '@/shad/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import axios from 'axios'
import { Field, useForm } from 'vee-validate'
import * as z from 'zod'
import { useToast } from '@/shad/components/ui/toast/use-toast'
import { defineEmits, ref, watch } from 'vue'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import Spinner from '../Spinner.vue'
import { type DateValue, getLocalTimeZone, today } from '@internationalized/date'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shad/components/ui/select'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shad/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shad/components/ui/popover'

const props = defineProps<{ date: DateValue; clerkId: number | null | undefined; hour: number; minute: number; availableDuration: number[]; slotNumber: number; }>()

const formSchema = toTypedSchema(z.object({
  username: z.string().min(2, { message: "Username must be at least 2 characters" }).max(50, { message: "Username cannot exceed 50 characters" }),
  duration: z.string(),
  user_id: z.number(),
}));

const { handleSubmit, setFieldValue, values, errors } = useForm({
  validationSchema: formSchema,
  initialValues: {
    username: '',
  }
})

const { toast } = useToast()
const emit = defineEmits(['meetingCreated'])
const loading = ref(false)


const generateTimeslot = (date: string, occupiedIds: number[]) => {
  const slot = {
    Date: date + "T00:00:00Z",
    ClerkId: props.clerkId,
    Occupied: occupiedIds,
  }
  return slot
}

const submitForm = async (formData: { username: string; duration: string; user_id: number }) => {
  const slotSpan = Number(formData.duration) / 30
  const occupiedIds = []
  for (let i = props.slotNumber; i < props.slotNumber + slotSpan; i++) {
    occupiedIds.push(i)
  }
  const slot = generateTimeslot(props.date.toString(), occupiedIds)
  const meeting = {
    user_id: formData.user_id,
    duration: Number(formData.duration),
    clerk_id: props.clerkId,
    start_time: new Date(props.date.year, props.date.month - 1, props.date.day, props.hour, props.minute, 0),
    time_slot_id: props.slotNumber
  }
  const data = { timeslot: slot, meeting: meeting }

  try {
    loading.value = true
    const response = await axios.post("/api/meeting", data)
    loading.value = false
    console.log('Response:', response.data)
    emit('meetingCreated')
  } catch (error: any) {
    loading.value = false
    const errorMessage = error.response?.data?.error || 'Please check your information again and try again.'
    console.error('Error:', error)
    toast({
      title: 'Creation Failed',
      description: errorMessage,
      variant: 'destructive'
    })
  }
}

const onSubmit = handleSubmit((values) => {
  submitForm(values)
})

const formattedTime = (hour: number, minute: number) =>
  `${String(hour).padStart(2, "0")}:${String(minute).padStart(2, "0")}`


interface User {
  id: number
  username: string
}
const users = ref<(User[])>([])
const searchInput = ref("")
const updateSearchTerm = (search: any) => { searchInput.value = search }

watch(searchInput, (newValue) => {
  if (newValue != '') {
    fetchUsers(newValue)
  }
})

const fetchUsers = async (searchTerm: string) => {
  const searchQuery = { Role: "Regular", Username: searchTerm }
  const params = {
    sortBy: "username",
  }
  axios.post("/api/user/query", searchQuery, { params: params }).then((result) => {
    users.value = result.data.users
  }).catch(error => { console.log(error) })
}

</script>

<template>
  <main>
    <Spinner v-if="loading" />
    <form class="w-full space-y-6" @submit="onSubmit">
      <FormField name="username" v-slot="{ field }">
        <FormItem class="relative pb-2 flex flex-col">
          <FormLabel>Meeting with</FormLabel>
          <Popover>
            <PopoverTrigger as-child>
              <FormControl>
                <Button variant="outline" role="combobox"
                  :class="['w-full justify-between', !values.username && 'text-muted-foreground']">
                  {{ values.username ? users.find(
                    (user) => user.username === values.username,
                  )?.username : 'Type to search' }}
                </Button>
              </FormControl>
            </PopoverTrigger>
            <PopoverContent class="w-[200px] p-0">
              <Command v-on:update:search-term="updateSearchTerm">
                <CommandInput placeholder="Search users..." />
                <CommandEmpty>Nothing found.</CommandEmpty>
                <CommandList>
                  <CommandGroup>
                    <CommandItem v-for="user in users" :key="user.username" :value="user.username" @select="() => {
                      setFieldValue('username', user.username)
                      setFieldValue('user_id', user.id)
                    }">
                      {{ user.username }}
                    </CommandItem>
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          <FormMessage class="absolute -bottom-2 left-0 text-xs" v-if="errors.username">{{ errors.username }}
          </FormMessage>
        </FormItem>
      </FormField>

      <FormField name="date" v-slot="{ field }">
        <FormItem class="relative pb-2">
          <FormLabel>Date</FormLabel>
          <FormControl>
            <Input type="text" v-bind="field" placeholder="Date" :default-value="date.toString()" disabled />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="start" v-slot="{ field }">
        <FormItem class="relative pb-2">
          <FormLabel>Start time</FormLabel>
          <FormControl>
            <Input type="text" v-bind="field" placeholder="Start time" :default-value="formattedTime(hour, minute)"
              disabled />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField name="duration" v-slot="{ field }">
        <FormItem class="relative pb-2">
          <FormLabel>Duration</FormLabel>
          <Select v-bind="field">
            <FormControl>
              <SelectTrigger>
                <SelectValue placeholder="Select duration" />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="(duration, index) in availableDuration" :key="index" :value="duration.toString()">
                  {{ duration + " minutes" }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormItem>
      </FormField>

      <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
        Create
      </Button>
    </form>
    <Toaster />
  </main>
</template>

<style scoped></style>
