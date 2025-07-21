<script setup lang="ts">
import { Calendar } from '@/shad/components/ui/calendar'
import { type DateValue, getLocalTimeZone, today } from '@internationalized/date'
import { computed, onMounted, type Ref, ref, watch } from 'vue'
import TimeSlot from './TimeSlot.vue';
import axios from 'axios';
import { useUserStore } from '@/stores/user';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shad/components/ui/dialog'
import NewMeetingFrom from './NewMeetingForm.vue';
import { useToast } from '@/shad/components/ui/toast';
import Toaster from '@/shad/components/ui/toast/Toaster.vue';

const userStore = useUserStore()
const isClerk = ref<boolean>(false)

interface TimeSlot {
  Date: string,
  Slots: Boolean[]
  ClerkId: number | null | undefined
  Id?: number
}

const emit = defineEmits(['meeting-id'])
const props = defineProps<{ userId: number | null, username: string | null }>();

const dateValue = ref(today(getLocalTimeZone())) as Ref<DateValue>
const hourValue = ref(0)
const minuteValue = ref(0)
const slotNumber = ref(-1)
const availableDuration = ref([0])
const isDialogOpen = ref(false)
const updateDialog = (open: boolean) => {
  isDialogOpen.value = open
}

const isPast = computed(() => {

  if (!dateValue.value) {
    return false
  }
  return dateValue.value.compare(today(getLocalTimeZone())) < 0
})

watch(() => props.userId, (newUserId) => {
  if (newUserId !== null) {
    console.log('User ID updated:', newUserId)
    fetchTimeTable(dateValue.value.toString(), newUserId);
  }
});

watch(dateValue, (newDate) => {
  if (newDate == undefined)
    newDate = today(getLocalTimeZone())
  fetchTimeTable(newDate.toString().trim(), props.userId)
})

onMounted(() => {
  if (userStore.role == "Clerk") {
    isClerk.value = true;
  }
  fetchTimeTable(today(getLocalTimeZone()).toString(), props.userId)
})

const fetchTimeTable = async (date: string, userId: number | null) => {
  try {
    const response = await axios.get("/api/timeslot", { params: { date: date, clerk_id: userId } })
    for (let i = 0; i < 15; i++) {
      slots.value[i] = { ...slots.value[i], meetingId: response.data.data.slots[i], past: isPast.value, id: response.data.data.id }
    }

  } catch (error: any) {
    if (error.status == 404) {
      console.log(isPast.value)
      if (isPast.value) {
        slots.value = generateSlots()
        return
      } else {
        slots.value = generateSlots()
      }
    } else {
      console.log("Error fetching timetable", error)
    }

  }

  const cmp = dateValue.value.compare(today(getLocalTimeZone()))
  if (cmp == 0) {  //today
    const todayDate = new Date()
    for (let i = 0; i < 15; i++) {
      if (slots.value[i].hour <= todayDate.getHours()) {
        slots.value[i] = { ...slots.value[i], past: true }
      }
    }
  }


}

const openSlot = async (slot: any) => {
  if (dateValue.value == undefined || props.userId == 0 || props.userId == null)
    return
  if (slot.meetingId != 0) {
    emit('meeting-id', slot.meetingId)
  } else {
    if (!slot.past) {
      console.log(slot.number)
      isDialogOpen.value = true
      hourValue.value = slot.hour
      minuteValue.value = slot.minute
      availableDuration.value = [30]
      slotNumber.value = slot.number
      let last = 30;
      console.log(slot.number)
      for (let i = slot.number + 1; i < 15; i++) {
        if (slots.value[i].meetingId == 0 && i != 8 && last != 180) {
          last += 30
          availableDuration.value.push(last)
        } else {
          break
        }
      }

    }
  }
}

const { toast } = useToast()
const generateSlots = () => {
  const slots = [];
  let i = 0;
  for (let hour = 8; hour < 16; hour++) {
    for (let minute of [0, 30]) {
      if (hour === 12 && minute === 0) continue; // Skip 12:00 to 12:30
      slots.push({ hour, minute, meetingId: 0, past: isPast.value, number: i, id: -1 });
      i++;
    }
  }
  return slots;
};

const splitIntoColumns = (slots: any, columnCount: any) => {
  const rowsPerColumn = Math.ceil(slots.length / columnCount);
  const columns = [];
  for (let i = 0; i < columnCount; i++) {
    columns.push(slots.slice(i * rowsPerColumn, (i + 1) * rowsPerColumn));
  }
  return columns;
};

const slots = ref(generateSlots());
const availableSlots = computed(() => splitIntoColumns(slots.value, 3))
const closeDialog = () => {
  updateDialog(false)
  fetchTimeTable(dateValue.value.toString(), props.userId)
  toast({
    title: 'Creation Successful',
    variant: 'default'
  })
}

const confirmMeeting = async () => {
  const occupiedIds = []
  for (let i = slotNumber.value; i < slotNumber.value + 1; i++) {
    occupiedIds.push(i)
  }
  const slot = {
    Date: dateValue.value.toString() + "T00:00:00Z",
    ClerkId: props.userId,
    Occupied: occupiedIds,
  }
  const meeting = {
    user_id: userStore.id,
    duration: 30,
    clerk_id: props.userId,
    start_time: new Date(dateValue.value.year, dateValue.value.month - 1, dateValue.value.day, hourValue.value, minuteValue.value, 0),
    time_slot_id: slotNumber.value
  }
  const data = { timeslot: slot, meeting: meeting }
  try {
    const response = await axios.post("api/meeting", data);
    console.log("Response:", response.data)
    toast({
      title: 'Meeting Scheduled',
      description: 'Your meeting has been successfully scheduled.',
    });
    await fetchTimeTable(dateValue.value.toString().trim(), props.userId)
    updateDialog(false);
  } catch (error: any) {
    const errorMessage = error.response?.data?.error || 'Please refresh page and try again.'
    console.error('Error:', error)
    toast({
      title: 'Creation Failed',
      description: errorMessage,
      variant: 'destructive'
    })
    updateDialog(false);
  }

};

const cancelMeeting = () => {
  updateDialog(false);
};

</script>

<template>
  <div class="w-fit flex flex-row justify-center gap-10">
    <div class="flex-1 bg-white rounded-md p-4 flex flex-col" id="1">
      <span class="text-gray-900 text-lg text-center my-5 w-full">Date</span>
      <Calendar v-model="dateValue" :weekday-format="'short'" class="rounded-md border h-full" :week-starts-on="1" />
    </div>

    <div class="flex bg-white rounded-md p-4 flex-col" id="2">
      <span class="text-gray-900 text-lg text-center my-5 w-full">Slots</span>
      <div class="grid grid-cols-3 gap-4 border border-gray-200 p-10 rounded-md">
        <div v-for="(column, colIndex) in availableSlots" :key="colIndex" class="flex flex-col gap-4">
          <TimeSlot v-for="(slot, index) in column" :key="index" :startHour="slot.hour" :startMinute="slot.minute"
            :isAvailable="slot.meetingId == 0" :hasClerk="props.userId != 0 && props.userId != null" :past="slot.past"
            @click.prevent="openSlot(slot)" />
        </div>
      </div>
    </div>
  </div>
  <Dialog :open="isDialogOpen" v-on:update:open="updateDialog">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>New meeting</DialogTitle>
      </DialogHeader>
      <template v-if="isClerk">
        <NewMeetingFrom :clerk-id="userStore?.id" :date="dateValue" :hour="hourValue" :minute="minuteValue"
          :available-duration="availableDuration" :slot-number="slotNumber" @meeting-created="closeDialog" />
      </template>
      <template v-else>
        <div class="text-center">
          <p>Are you sure you want to schedule meeting </p>
          <p>with {{ username }} on {{ dateValue }} at {{ hourValue }}:{{ minuteValue }}?</p>
          <div class="mt-4 flex justify-center gap-4">
            <button class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600" @click="confirmMeeting">
              Yes
            </button>
            <button class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400" @click="cancelMeeting">
              No
            </button>
          </div>
        </div>
      </template>
    </DialogContent>
  </Dialog>

  <Toaster />
</template>
