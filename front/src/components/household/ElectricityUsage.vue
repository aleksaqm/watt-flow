<script setup lang="ts">
import { BarChart } from '@/shad/components/ui/chart-bar'
import Card from '@/shad/components/ui/card/Card.vue'
import CardContent from '@/shad/components/ui/card/CardContent.vue'
import CardHeader from '@/shad/components/ui/card/CardHeader.vue'
import CardTitle from '@/shad/components/ui/card/CardTitle.vue'
import Button from '@/shad/components/ui/button/Button.vue'
import { reactive, ref, computed, onMounted } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import axios from 'axios'

const props = defineProps({
  householdId: {
    type: String,
    required: true
  }
})

interface MonthlyData {
  month: string
  year: number
  consumption: number
}

interface DailyData {
  day: string
  consumption: number
}

type ChartDataItem = MonthlyData | DailyData

const currentView = ref<'monthly' | 'daily'>('monthly')

const dailyViewYear = ref(new Date().getFullYear())
const dailyViewMonth = ref(new Date().getMonth() + 1)

const selectedTimePeriod = ref("")
const showTimePeriodSelection = ref(false)
const isLoading = ref(false)

const fetchMonthlyConsumption = async (year: number, month: number) => {
  try {
    const response = await axios.get(`/api/household/${props.householdId}/consumption/monthly`, {
      params: { year, month }
    })
    return response.data.data
  } catch (error) {
    console.error('Error fetching monthly consumption:', error)
    return { year, month, consumption: 0, month_name: getMonthName(month) }
  }
}

const fetch12MonthsConsumption = async (endYear: number, endMonth: number) => {
  try {
    isLoading.value = true
    const response = await axios.get(`/api/household/${props.householdId}/consumption/12months`, {
      params: { endYear, endMonth }
    })
    return response.data.data
  } catch (error) {
    console.error('Error fetching 12 months consumption:', error)
    return []
  } finally {
    isLoading.value = false
  }
}

const fetchDailyConsumption = async (year: number, month: number) => {
  try {
    isLoading.value = true
    const response = await axios.get(`/api/household/${props.householdId}/consumption/daily`, {
      params: { year, month }
    })
    return response.data.data
  } catch (error) {
    console.error('Error fetching daily consumption:', error)
    return []
  } finally {
    isLoading.value = false
  }
}

const getMonthName = (month: number) => {
  const monthNames = ['', 'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December']
  return monthNames[month] || `Month ${month}`
}

const monthlyData = reactive<MonthlyData[]>([])

const leftmostYear = ref(new Date().getFullYear())
const rightmostYear = ref(new Date().getFullYear())

const dailyData = reactive<DailyData[]>([])

const loadDailyData = async (year: number, month: number) => {
  const data = await fetchDailyConsumption(year, month)
  
  dailyData.splice(0, dailyData.length)
  data.forEach((item: any) => {
    dailyData.push({
      day: item.day.toString(),
      consumption: item.consumption
    })
  })
}

const initializeData = async () => {
  const now = new Date()
  const currentYear = now.getFullYear()
  const currentMonth = now.getMonth() + 1
  
  const data = await fetch12MonthsConsumption(currentYear, currentMonth)
  
  monthlyData.splice(0, monthlyData.length) // Clear array
  data.forEach((item: any, index: number) => {
    let itemYear = currentYear
    const itemMonthNum = item.month
    
    if (itemMonthNum > currentMonth) {
      itemYear = currentYear - 1
    }
    
    monthlyData.push({
      month: getMonthShortName(item.month),
      year: itemYear,
      consumption: item.consumption
    })
  })
  
  if (monthlyData.length > 0) {
    const oldestMonthNum = getMonthNumber(monthlyData[0].month)
    const newestMonthNum = getMonthNumber(monthlyData[11].month)
    
    if (newestMonthNum >= currentMonth) {
      rightmostYear.value = currentYear
      leftmostYear.value = currentYear - 1
    } else {
      rightmostYear.value = currentYear
      leftmostYear.value = currentYear - 1
    }
  }
}

const getMonthShortName = (month: number) => {
  const shortNames = ['', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
    'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  return shortNames[month] || `${month}`
}

onMounted(() => {
  initializeData()
})

const chartData = computed(() => {
  if (currentView.value === 'monthly') {
    return monthlyData.map(item => ({
      month: `${item.month} ${item.year}`,
      originalMonth: item.month,
      year: item.year,
      consumption: item.consumption
    })) as any[]
  } else {
    return dailyData as ChartDataItem[]
  }
})

const chartIndex = computed(() => {
  return currentView.value === 'monthly' ? ('month' as keyof ChartDataItem) : ('day' as keyof ChartDataItem)
})

const chartCategories = computed(() => {
  return ['consumption' as keyof ChartDataItem]
})

const chartTitle = computed(() => {
  if (currentView.value === 'monthly') {
    return `Monthly Electricity Usage`
  } else {
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December']
    return `Daily Electricity Usage - ${monthNames[dailyViewMonth.value - 1]} ${dailyViewYear.value}`
  }
})

const goToPreviousPeriod = async () => {
  if (currentView.value === 'monthly') {
    if (monthlyData.length === 0) return
    
    const oldestMonth = monthlyData[0]
    const oldestMonthNum = getMonthNumber(oldestMonth.month)
    
    let newMonth = oldestMonthNum - 1
    let newYear = leftmostYear.value
    
    if (newMonth === 0) {
      newMonth = 12
      newYear--
    }
    
    const newData = await fetchMonthlyConsumption(newYear, newMonth)
    
    monthlyData.unshift({
      month: getMonthShortName(newData.month),
      year: newYear,
      consumption: newData.consumption
    })
    monthlyData.pop()
    
    leftmostYear.value = newYear
    
    monthlyData.forEach((item, index) => {
      const monthNum = getMonthNumber(item.month)
      if (monthNum >= getMonthNumber(monthlyData[0].month)) {
        item.year = leftmostYear.value
      } else {
        item.year = leftmostYear.value + 1
      }
    })
    
    const newestMonth = monthlyData[11]
    const newestMonthNum = getMonthNumber(newestMonth.month)
    if (newestMonthNum < newMonth) {
      rightmostYear.value = leftmostYear.value + 1
    } else {
      rightmostYear.value = leftmostYear.value
    }
    
  } else {
    let newMonth = dailyViewMonth.value - 1
    let newYear = dailyViewYear.value
    
    if (newMonth === 0) {
      newMonth = 12
      newYear--
    }
    
    dailyViewMonth.value = newMonth
    dailyViewYear.value = newYear
    
    await loadDailyData(newYear, newMonth)
  }
  
  console.log('Navigated to previous period:', currentView.value === 'monthly' ? 
    `${monthlyData[0]?.month} ${leftmostYear.value} - ${monthlyData[11]?.month} ${rightmostYear.value}` : 
    `${dailyViewMonth.value}/${dailyViewYear.value}`)
}

const goToNextPeriod = async () => {
  if (currentView.value === 'monthly') {
    if (monthlyData.length === 0) return
    
    const newestMonth = monthlyData[11]
    const newestMonthNum = getMonthNumber(newestMonth.month)
    
    let newMonth = newestMonthNum + 1
    let newYear = rightmostYear.value
    
    if (newMonth === 13) {
      newMonth = 1
      newYear++
    }
    
    const newData = await fetchMonthlyConsumption(newYear, newMonth)
    
    monthlyData.shift()
    monthlyData.push({
      month: getMonthShortName(newData.month),
      year: newYear,
      consumption: newData.consumption
    })
    
    rightmostYear.value = newYear
    
    monthlyData.forEach((item, index) => {
      const monthNum = getMonthNumber(item.month)
      if (monthNum <= getMonthNumber(monthlyData[11].month)) {
        item.year = rightmostYear.value
      } else {
        item.year = rightmostYear.value - 1
      }
    })
    
    const oldestMonth = monthlyData[0]
    const oldestMonthNum = getMonthNumber(oldestMonth.month)
    if (oldestMonthNum > newMonth) {
      leftmostYear.value = rightmostYear.value - 1
    } else {
      leftmostYear.value = rightmostYear.value
    }
    
  } else {
    let newMonth = dailyViewMonth.value + 1
    let newYear = dailyViewYear.value
    
    if (newMonth === 13) {
      newMonth = 1
      newYear++
    }
    
    dailyViewMonth.value = newMonth
    dailyViewYear.value = newYear
    
    await loadDailyData(newYear, newMonth)
  }
  
  console.log('Navigated to next period:', currentView.value === 'monthly' ? 
    `${monthlyData[0]?.month} ${leftmostYear.value} - ${monthlyData[11]?.month} ${rightmostYear.value}` : 
    `${dailyViewMonth.value}/${dailyViewYear.value}`)
}

const getMonthNumber = (shortName: string) => {
  const monthMap: { [key: string]: number } = {
    'Jan': 1, 'Feb': 2, 'Mar': 3, 'Apr': 4, 'May': 5, 'Jun': 6,
    'Jul': 7, 'Aug': 8, 'Sep': 9, 'Oct': 10, 'Nov': 11, 'Dec': 12
  }
  return monthMap[shortName] || 1
}

const handleBarClick = async (data: any, category: string, index: number) => {
  if (currentView.value === 'monthly') {
    const clickedMonthName = data.originalMonth || data.month.split(' ')[0] // Get original month name
    const clickedMonthNum = getMonthNumber(clickedMonthName)
    const clickedYear = data.year
    
    dailyViewYear.value = clickedYear
    dailyViewMonth.value = clickedMonthNum
    currentView.value = 'daily'
    
    await loadDailyData(clickedYear, clickedMonthNum)
    
    console.log(`Switching to daily view for: ${clickedMonthName} ${clickedYear}`)
  }
}

const goBackToMonthly = () => {
  currentView.value = 'monthly'
}

const yFormatter = (value: number) => {
  return `${value.toFixed(2)} kWh`
}

</script>

<template>
  <div class="flex flex-col gap-3 items-center justify-center w-full mb-10">
    <Card class="w-5/6">
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle>
            <span class="text-gray-600 text-xl">{{ chartTitle }}</span>
          </CardTitle>
          
          <div class="flex items-center gap-2">
            <Button 
              variant="outline" 
              size="sm" 
              @click="goToPreviousPeriod"
            >
              <ChevronLeft class="h-4 w-4" />
            </Button>
            
            <Button 
              variant="outline" 
              size="sm" 
              @click="goToNextPeriod"
            >
              <ChevronRight class="h-4 w-4" />
            </Button>
            
            <!-- Back to monthly view button -->
            <Button 
              v-if="currentView === 'daily'"
              variant="outline" 
              size="sm" 
              @click="goBackToMonthly"
            >
              Back to Monthly
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div class="text-xs text-gray-500 mb-4">
          kWh - Kilowatt hours
        </div>
        <div class="w-full">
          <div v-if="isLoading" class="flex items-center justify-center h-64">
            <div class="text-gray-500">Loading consumption data...</div>
          </div>
          
          <div v-else class="w-full min-h-[400px]">
            <BarChart
              :data="chartData"
              :index="chartIndex"
              :categories="chartCategories"
              :y-formatter="yFormatter"
              :colors="['#3b82f6']"
              :bar-spacing="0.4"
              @bar-click="handleBarClick"
              class="w-full"
            />
          </div>
          
          <div class="text-center text-sm text-gray-500 mt-8"> 
            {{ currentView === 'monthly' ? 'Click on a month to see daily consumption' : 'Daily electricity consumption' }}
          </div>
        </div>
      </CardContent>
    </Card>
  </div>

</template>

<style scoped>
/* Additional styles if needed */
</style>
