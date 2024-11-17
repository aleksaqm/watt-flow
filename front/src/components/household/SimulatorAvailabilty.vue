<script setup lang="ts">
import Card from '@/shad/components/ui/card/Card.vue';
import CardContent from '@/shad/components/ui/card/CardContent.vue';
import CardHeader from '@/shad/components/ui/card/CardHeader.vue';
import CardTitle from '@/shad/components/ui/card/CardTitle.vue';
import { AreaChart } from '@/shad/components/ui/chart-area'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shad/components/ui/select'

import {
  CalendarDate,
} from '@internationalized/date'
import DateRangePicker from './DateRangePicker.vue';
import { reactive, ref, watch } from 'vue';
import Button from '@/shad/components/ui/button/Button.vue';
import axios from 'axios';
import { toDate } from 'radix-vue/date';


const now = new Date()
const end = new CalendarDate(now.getFullYear(), now.getMonth(), now.getDay())

const selectedDates = ref({
  end: end,
  start: end.subtract({ days: 10 }),
})



const isCalendarEnabled = ref(false)
const selectedTimePeriod = ref("")
let selectedGroupPeriod = ""

watch(selectedTimePeriod, (new_value) => {
  isCalendarEnabled.value = new_value === "custom"
})
const GroupPeriodMap: { [key: string]: string } = {
  "3h": "1h",
  "6h": "1h",
  "12h": "1h",
  "24h": "1h",
  "7d": "1d",
  "30d": "1d",
  "90d": "1d",
  "365d": "7d",
}

const MinutesInGroupPeriod: { [key: string]: number } = {
  "1h": 60,
  "1d": 1440,
  "7d": 10080,
}


const PrecisionMap: { [key: string]: string } = {
  "3h": "m",
  "6h": "m",
  "12h": "m",
  "24h": "m",
  "7d": "m",
  "30d": "m",
  "90d": "m",
  "365d": "m"
}
interface ChartValue {
  time: string,
  value: number
}

const chartData = reactive<{
  data: ChartValue[]
  config: any[]

}>({
  data: [],
  config: [],
})

interface FluxQuery {
  TimePeriod: string,
  GroupPeriod: string,
  DeviceId: string,
  Precision: string,
  StartDate?: Date,
  EndDate?: Date
}

const handleFetch = () => {
  if (selectedTimePeriod.value == "") {
    return
  }
  let query: FluxQuery | null = null
  if (selectedTimePeriod.value == "custom") {
    const startDate = selectedDates.value.start.toDate("Europe/Sarajevo")
    const endDate = selectedDates.value.end.toDate("Europe/Sarajevo")
    const difference = (endDate.getTime() - startDate.getTime()) / 3600000
    console.log("diff", difference)
    if (difference <= 24) {
      selectedGroupPeriod = "1h"
    } else if (difference <= 720) {
      selectedGroupPeriod = "1d"
    } else {
      selectedGroupPeriod = "7d"
    }

    query = {
      TimePeriod: selectedTimePeriod.value,
      GroupPeriod: selectedGroupPeriod,
      Precision: "m",
      DeviceId: "be781b42-c3b0-475b-bdc5-cb467d0f4fa1",
      StartDate: startDate,
      EndDate: endDate,
    }

  } else {
    selectedGroupPeriod = GroupPeriodMap[selectedTimePeriod.value]

    query = {
      TimePeriod: selectedTimePeriod.value,
      GroupPeriod: GroupPeriodMap[selectedTimePeriod.value],
      Precision: PrecisionMap[selectedTimePeriod.value],
      DeviceId: "be781b42-c3b0-475b-bdc5-cb467d0f4fa1",
    }
  }

  axios.post('/api/device-status/query-status', query).then(
    (result) => {
      formatData(result.data.data.Rows)
    }
  ).catch((error) => console.log(error))

}

const formatData = (data: any[]) => {
  const length = data.length
  const standardUnit = MinutesInGroupPeriod[selectedGroupPeriod]

  const lastTick = new Date(data[length - 1].TimeField).getTime()
  const secondTick = new Date(data[length - 2].TimeField).getTime()

  const firstUnit = (lastTick - secondTick) / 60000




  let lastData = data[length - 1].Value
  let currentValue = 0
  let unit = 0
  let remainder = 0

  chartData.data = Array.from({ length: length - 1 })

  for (let i = length - 2; i >= 0; i--) {
    if (i == length - 2) {
      unit = firstUnit
    } else {
      unit = standardUnit
    }
    remainder += lastData
    lastData = data[i].Value
    if (remainder > unit) {
      remainder -= unit
      currentValue = unit
    } else {
      currentValue = remainder
      remainder = 0
    }
    chartData.data[i] =
    {
      "time": xFormatter(new Date(data[i].TimeField)),
      "value": (currentValue / unit) * 100,
    }
  }
}

const xFormatter = (date: Date) => {
  switch (selectedTimePeriod.value) {
    case "3h": case "6h": case "12h": case "24h":
      return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        hourCycle: "h24"
      })
    case "7d": case "30d": case "90d":
      return date.toLocaleDateString("en-US", {
        day: "numeric",
        month: "short",
      })
    case "365d": case "custom":
      return date.toLocaleDateString("en-US", {
        day: "numeric",
        month: "short",
      })
    default:
      return ""
  }
}




</script>
<template>
  <div class="flex flex-col gap-3 items-center justify-center w-full mb-10">
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
                <SelectItem value="7d">
                  Last week
                </SelectItem>
                <SelectItem value="30d">
                  Last month
                </SelectItem>
                <SelectItem value="90d">
                  Last 3 months
                </SelectItem>
                <SelectItem value="365d">
                  Last year
                </SelectItem>
                <SelectItem value="custom">
                  Custom
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <DateRangePicker :is-calendar-enabled="isCalendarEnabled" v-model:value="selectedDates"> </DateRangePicker>
          <Button @click="handleFetch">Fetch</Button>

        </div>
        <div class="p-10">
          <AreaChart :show-legend="false" :data="chartData.data" index="time" :categories="['value']">
          </AreaChart>

        </div>

      </CardContent>
    </Card>

  </div>

</template>

<style scoped></style>
