<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/shad/components/ui/card'

defineProps<{
  title?: string
  data: {
    name: string
    color: string
    value: any  // consumption value in kWh
  }[]
  isRealtime: boolean
  unit: number
}>()
</script>

<template>
  <Card class="text-sm">
    <CardHeader v-if="title" class="p-3 border-b">
      <CardTitle>
        {{ title }}
      </CardTitle>
    </CardHeader>
    <CardContent class="p-3 min-w-[180px] flex flex-col gap-1">
      <div class="flex justify-between">
        <div class="flex items-center">
          <span class="w-2.5 h-2.5 mr-2">
            <svg width="100%" height="100%" viewBox="0 0 30 30">
              <path d=" M 15 15 m -14, 0 a 14,14 0 1,1 28,0 a 14,14 0 1,1 -28,0" :stroke="data[0].color"
                :fill="data[0].color" stroke-width="1" />
            </svg>
          </span>
          <span>Consumption</span>
        </div>
        <span class="font-semibold ml-4">{{ Math.round(data[0].value * 100) / 100 }} kWh</span>
      </div>
    </CardContent>
  </Card>
</template>
