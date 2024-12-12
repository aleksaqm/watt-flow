<script setup lang="ts">
import { Calendar } from '@/shad/components/ui/calendar'
import { type DateValue, getLocalTimeZone, today } from '@internationalized/date'
import { type Ref, ref } from 'vue'
import TimeSlot from './TimeSlot.vue';

const value = ref(today(getLocalTimeZone())) as Ref<DateValue>

const generateSlots = () => {
  const slots = [];
  for (let hour = 8; hour < 16; hour++) {
    for (let minute of [0, 30]) {
      if (hour === 12 && minute === 0) continue; // Skip 12:00 to 12:30
      slots.push({ hour, minute, isAvailable: true });
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

const availableSlots = generateSlots();
const columns = splitIntoColumns(availableSlots, 3);
</script>

<template>
  <div class="w-fit flex flex-row justify-center gap-10">
    <div class="flex-1 bg-white rounded-md p-4 flex flex-col" id="1">
      <span class="text-gray-900 text-lg text-center my-5 w-full">Date</span>
      <Calendar v-model="value" :weekday-format="'short'" class="rounded-md border h-full" :week-starts-on="1" />
    </div>

    <div class="flex bg-white rounded-md p-4 flex-col" id="2">
      <span class="text-gray-900 text-lg text-center my-5 w-full">Slots</span>
      <div class="grid grid-cols-3 gap-4 border border-gray-200 p-10 rounded-md">
        <div v-for="(column, colIndex) in columns" :key="colIndex" class="flex flex-col gap-4">
          <TimeSlot v-for="(slot, index) in column" :key="index" :startHour="slot.hour" :startMinute="slot.minute"
            :isAvailable="slot.isAvailable" />
        </div>
      </div>
    </div>
  </div>

</template>
