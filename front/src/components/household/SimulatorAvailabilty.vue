<script setup lang="ts">
import Card from '@/shad/components/ui/card/Card.vue';
import CardContent from '@/shad/components/ui/card/CardContent.vue';
import CardHeader from '@/shad/components/ui/card/CardHeader.vue';
import CardTitle from '@/shad/components/ui/card/CardTitle.vue';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/shad/components/ui/select'
import DateRangePicker from './DateRangePicker.vue';
import { ref, watch } from 'vue';
import Button from '@/shad/components/ui/button/Button.vue';
import LineChart from './LineChart.vue';

const isCalendarEnabled = ref(false)
const selectedTimePeriod = ref("")

watch(selectedTimePeriod, (new_value) => {
  isCalendarEnabled.value = new_value === "custom"
})

const handleFetch = async () => {
  console.log("clicked")
  console.log(selectedTimePeriod)

}



</script>
<template>
  <div class="flex flex-col gap-3 items-center justify-center w-full">
    <Card class="w-3/4">
      <CardHeader>
        <CardTitle>
          <span class="text-gray-600 text-xl">Power meter availability</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div class="flex justify-around items-center spaxe-x-10 gap-10">
          <Select v-model="selectedTimePeriod">
            <SelectTrigger>
              <SelectValue placeholder="Select time period" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="3h">
                  Last 3h
                </SelectItem>
                <SelectItem value="6h">
                  Last 6h
                </SelectItem>
                <SelectItem value="12h">
                  Last 12h
                </SelectItem>
                <SelectItem value="24h">
                  Last 24h
                </SelectItem>
                <SelectItem value="week">
                  Last week
                </SelectItem>
                <SelectItem value="month">
                  Last month
                </SelectItem>
                <SelectItem value="3month">
                  Last 3 months
                </SelectItem>
                <SelectItem value="year">
                  Last year
                </SelectItem>
                <SelectItem value="custom">
                  Custom
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <DateRangePicker :is-calendar-enabled="isCalendarEnabled"> </DateRangePicker>
          <Button @click="handleFetch">Fetch</Button>

        </div>
        <div>
          <LineChart></LineChart>

        </div>

      </CardContent>
    </Card>

  </div>

</template>

<style scoped></style>
