<script setup lang="ts">
import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { useToast } from '@/shad/components/ui/toast/use-toast'
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
import Checkbox from '@/shad/components/ui/checkbox/Checkbox.vue';
import ConsumptionTooltip from './ConsumptionTooltip.vue';


const { toast } = useToast()


const props = defineProps({
  deviceId: {
    type: String,
    default: ""
  }
})

const now = new Date()
const end = new CalendarDate(now.getFullYear(), now.getMonth(), now.getDay())

const selectedDates = ref({
  end: end,
  start: end.subtract({ days: 10 }),
})



const isCalendarEnabled = ref(false)
const selectedTimePeriod = ref("")
let selectedGroupPeriod = ""
const isRealtimeEnabled = ref(false)
const isRealtimeSelected = ref(false)

const handleChange = (value: boolean) => {
  isRealtimeSelected.value = value
}

watch(selectedTimePeriod, (new_value) => {
  isCalendarEnabled.value = new_value === "custom"
  isRealtimeEnabled.value = new_value === "3h"
  if (!isRealtimeEnabled.value) isRealtimeSelected.value = false
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
  value: number,
}

const chartData = reactive<{
  data: ChartValue[]
  config: any[]
  unit: number
  totalConsumption: number
  averageConsumption: number

}>({
  data: [],
  config: [],
  unit: 0,
  totalConsumption: 0,
  averageConsumption: 0,
})

interface FluxQuery {
  TimePeriod: string,
  GroupPeriod: string,
  DeviceId: string,
  Precision: string,
  StartDate?: Date,
  EndDate?: Date
  Realtime: boolean
}
let lastConsumptionValue = -1
let refreshJob = -1

let ws: WebSocket | null = null
const serverUrl = `ws://localhost:9000/ws?deviceId=${String(props.deviceId)}&connType=consumption`
const isConnectedToWs = ref(false)

const connectToWebSocket = async () => {
  const authToken = localStorage.getItem("authToken")
  if (authToken != null) {
    console.log(`Connecting to WebSocket: ${serverUrl}`)
    console.log(`Using auth token: ${authToken.substring(0, 20)}...`) // Log partial token for debugging
    
    try {
      ws = new WebSocket(serverUrl, [authToken, "token"])
      
      ws.addEventListener("open", (event) => { 
        console.log("Connected to consumption WebSocket!"); 
        isConnectedToWs.value = true 
      })
      
      ws.addEventListener("message", (event) => {
        console.log("Received WebSocket message:", event.data)
        handleConsumptionUpdate(event)
      })
      
      ws.addEventListener("error", (event) => {
        console.error("WebSocket error:", event)
        console.error("WebSocket state:", ws?.readyState)
      })
      
      ws.addEventListener("close", (event) => {
        console.log("WebSocket connection closed:", event.code, event.reason)
        if (event.code === 1006) {
          console.log("Connection closed abnormally - likely authentication failure")
        }
        isConnectedToWs.value = false
      })
    } catch (error) {
      console.error("Failed to create WebSocket connection:", error)
    }
  } else {
    console.log("Auth token not found in localStorage!")
    return
  }
}

const updateRealtimeChart = () => {
  if (lastConsumptionValue != -1) {
    chartData.data.push(

      {
        "time": xFormatter(new Date()),
        "value": lastConsumptionValue,
      }
    )
    chartData.data.shift()
    console.log("Updated chart from last value!")
  }

}

const handleConsumptionUpdate = (event: any) => {
  const data = JSON.parse(event.data)
  chartData.data.push(
    {
      "time": xFormatter(new Date()),
      "value": data.Consumption || 0,
    }
  )
  chartData.data.shift()
  lastConsumptionValue = data.Consumption || 0
  console.log("Received consumption change from server!")
}

const handleFetch = () => {
  console.log(props.deviceId)
  let query: FluxQuery | null = null
  if (isRealtimeSelected.value) {
    if (ws?.readyState !== WebSocket.OPEN) {
      connectToWebSocket()
      refreshJob = setInterval(updateRealtimeChart, 60000)
      console.log(refreshJob)
    } else {
      return
    }
    query = {
      TimePeriod: "3h",
      GroupPeriod: "1m",
      Precision: "m",
      DeviceId: String(props.deviceId),
      Realtime: true
    }
  } else {
    if (ws != null)
      ws.close()
    lastConsumptionValue = -1
    if (refreshJob != -1) {
      clearInterval(refreshJob)
      console.log("Cleared job")
    }
  }

  if (query == null) {
    switch (selectedTimePeriod.value) {
      case "":
        toast({
          title: 'Fetch Failed',
          description: "Please select time period!",
          variant: 'default',
        })
        return
      case "custom":
        const startDate = selectedDates.value.start.toDate("Europe/Sarajevo")
        const endDate = selectedDates.value.end.toDate("Europe/Sarajevo")
        const difference = (endDate.getTime() - startDate.getTime()) / 3600000
        if (difference >= 8760) {
          toast({
            title: 'Invalid time period',
            description: "Maximum time period is one year!",
            variant: 'default',
          })
          return

        }

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
          DeviceId: String(props.deviceId),
          StartDate: startDate,
          EndDate: endDate,
          Realtime: false
        }
        break

      default:
        selectedGroupPeriod = GroupPeriodMap[selectedTimePeriod.value]

        query = {
          TimePeriod: selectedTimePeriod.value,
          GroupPeriod: GroupPeriodMap[selectedTimePeriod.value],
          Precision: PrecisionMap[selectedTimePeriod.value],
          DeviceId: String(props.deviceId),
          Realtime: false
        }
    }
  }
  // Changed endpoint to consumption
  axios.post('/api/device-consumption/query-consumption', query).then(
    (result) => {
      console.log("API Response:", result.data)
      if (isRealtimeSelected.value) {
        formatRealtimeData(result.data.data.Rows)
      } else {
        formatData(result.data.data.Rows)
      }
    }
  ).catch((error) => {
    console.error("API Error:", error)
    toast({
      title: 'Query Failed',
      description: "Failed to fetch consumption data",
      variant: 'destructive',
    })
  })

}

const formatRealtimeData = (data: any[]) => {
  chartData.data = []
  if (!data || data.length === 0) {
    console.log("No realtime data received")
    return
  }
  
  const length = data.length
  for (let i = 0; i < length; i++) {
    chartData.data.push({
      "time": xFormatter(new Date(data[i].TimeField)),
      "value": data[i].Value || 0,
    })
  }
  lastConsumptionValue = data[data.length - 1]?.Value || 0
  chartData.unit = 0 // Not used for consumption
}

const formatData = (data: any[]) => {
  const length = data.length
  let totalConsumption = 0

  chartData.data = []
  chartData.unit = 0 // Not used for consumption

  for (let i = 0; i < length; i++) {
    const consumptionValue = data[i].Value
    totalConsumption += consumptionValue
    
    chartData.data.push({
      "time": xFormatter(new Date(data[i].TimeField)),
      "value": consumptionValue, // actual consumption value in kWh
    })
  }

  chartData.totalConsumption = totalConsumption
  chartData.averageConsumption = length > 0 ? totalConsumption / length : 0
}

const xFormatter = (date: Date) => {
  switch (selectedTimePeriod.value) {
    case "3h":
      if (isRealtimeSelected.value)
        return date.toLocaleTimeString("en-US", {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hourCycle: "h24"
        })

      return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        hourCycle: "h24"
      })

    case "6h": case "12h": case "24h":
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
    <Card class="w-5/6">
      <CardHeader>
        <CardTitle>
          <span class="text-gray-600 text-xl">Electricity consumption</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div class="flex justify-around items-center gap-10">
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
          <Checkbox id="realtime" :checked="isRealtimeSelected" @update:checked="handleChange"
            :disabled="!isRealtimeEnabled" />
          <label for="realtime"
            class="text-sm text-gray-600 font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Realtime</label>
          <DateRangePicker :is-calendar-enabled="isCalendarEnabled" v-model:value="selectedDates"> </DateRangePicker>
          <Button @click="handleFetch">Fetch</Button>

        </div>
        <div class="p-10">
          <AreaChart :show-legend="false" :data="chartData.data" index="time" :categories="['value']"
            :custom-tooltip="ConsumptionTooltip" :is-realtime="isRealtimeSelected" :unit="chartData.unit">
          </AreaChart>
          <div class="flex flex-row gap-20 justify-center text-sm p-5 text-gray-600" v-if="!isRealtimeSelected">
            <div class="flex flex-row items-center justify-center">
              <span>Total consumption for selected period: </span>
              <span class='text-indigo-500 ml-5'>{{ Math.round(chartData.totalConsumption * 100) / 100 }} kWh</span>
            </div>

            <div class="flex flex-row items-center justify-center">
              <span>Average consumption: </span>
              <span class='text-indigo-500 ml-5'>{{ Math.round(chartData.averageConsumption * 100) / 100 }} kWh</span>
            </div>
          </div>


        </div>

      </CardContent>
    </Card>

  </div>
  <Toaster />

</template>

<style scoped></style>
