<script setup lang="ts" generic="T extends Record<string, any>">
import type { BaseChartProps } from '.'
import { ChartLegend, defaultColors } from '@/shad/components/ui/chart'
import { cn } from '@/lib/utils'
import { type BulletLegendItemInterface } from '@unovis/ts'
import { computed, ref } from 'vue'

const props = withDefaults(defineProps<BaseChartProps<T> & {
  /**
   * Controls the spacing between bars.
   * @default 0.1
   */
  barSpacing?: number
  /**
   * Controls whether bars should be rounded.
   * @default false
   */
  roundedBars?: boolean
}>(), {
  filterOpacity: 0.2,
  margin: () => ({ top: 20, bottom: 60, left: 70, right: 30 }),
  showXAxis: true,
  showYAxis: true,
  showTooltip: true,
  showLegend: true,
  showGridLine: true,
  barSpacing: 0.05,
  roundedBars: false,
})

const emits = defineEmits<{
  legendItemClick: [d: BulletLegendItemInterface, i: number]
  barClick: [data: T, category: string, index: number]
}>()

type KeyOfT = Extract<keyof T, string>
type Data = typeof props.data[number]

const index = computed(() => props.index as KeyOfT)
const colors = computed(() => props.colors?.length ? props.colors : defaultColors(props.categories.length))

const legendItems = ref<BulletLegendItemInterface[]>(props.categories.map((category, i) => ({
  name: category as string,
  color: colors.value[i],
  inactive: false,
})))

const hoveredBar = ref<{ dataIndex: number, categoryIndex: number } | null>(null)
const tooltipData = ref<{ x: number, y: number, data: any } | null>(null)

const chartWidth = 800
const chartHeight = 450
const marginTop = props.margin.top || 20
const marginBottom = props.margin.bottom || 60
const marginLeft = props.margin.left || 70
const marginRight = props.margin.right || 30
const innerWidth = chartWidth - marginLeft - marginRight
const innerHeight = chartHeight - marginTop - marginBottom

const maxValue = computed(() => {
  let max = 0
  props.data.forEach(d => {
    props.categories.forEach(category => {
      const value = d[category] as number
      if (value > max) max = value
    })
  })
  return max * 1.1 // Add 10% padding
})

const xScale = computed(() => {
  const totalBarWidth = innerWidth / props.data.length
  return (index: number) => index * totalBarWidth + totalBarWidth / 2
})

const yScale = computed(() => {
  return (value: number) => innerHeight - (value / maxValue.value) * innerHeight
})

const barWidth = computed(() => {
  const totalBarWidth = innerWidth / props.data.length
  const availableWidth = totalBarWidth * (1 - props.barSpacing)
  return availableWidth / props.categories.length
})

function handleLegendItemClick(d: BulletLegendItemInterface, i: number) {
  emits('legendItemClick', d, i)
}

function handleBarClick(data: T, categoryIndex: number, dataIndex: number) {
  const category = props.categories[categoryIndex] as string
  emits('barClick', data, category, dataIndex)
}

function handleBarHover(event: MouseEvent, data: T, categoryIndex: number, dataIndex: number) {
  const rect = (event.target as SVGElement).getBoundingClientRect()
  const containerRect = (event.target as SVGElement).closest('.relative')?.getBoundingClientRect()
  
  if (containerRect) {
    tooltipData.value = {
      x: rect.left + rect.width / 2 - containerRect.left,
      y: rect.top - containerRect.top - 10,
      data: {
        item: data,
        category: props.categories[categoryIndex] as string,
        value: data[props.categories[categoryIndex]] as number
      }
    }
  }
  hoveredBar.value = { dataIndex, categoryIndex }
}

function handleBarLeave() {
  hoveredBar.value = null
  tooltipData.value = null
}

const gridLines = computed(() => {
  const lines = []
  const step = maxValue.value / 5
  for (let i = 0; i <= 5; i++) {
    const value = i * step
    const y = yScale.value(value)
    lines.push({ value, y })
  }
  return lines
})
</script>

<template>
  <div :class="cn('w-full h-[400px] flex flex-col items-end', $attrs.class ?? '')">
    <ChartLegend v-if="showLegend" v-model:items="legendItems" @legend-item-click="handleLegendItemClick" />

    <div class="relative w-full flex-1">
      <svg :width="chartWidth" :height="chartHeight" class="overflow-visible">
        <!-- Grid lines -->
        <g v-if="showGridLine">
          <line
            v-for="line in gridLines"
            :key="line.value"
            :x1="marginLeft"
            :y1="line.y + marginTop"
            :x2="marginLeft + innerWidth"
            :y2="line.y + marginTop"
            stroke="hsl(var(--border))"
            stroke-width="1"
            stroke-dasharray="2,2"
            opacity="0.3"
          />
        </g>

        <g v-if="showYAxis">
          <line
            :x1="marginLeft"
            :y1="marginTop"
            :x2="marginLeft"
            :y2="marginTop + innerHeight"
            stroke="hsl(var(--border))"
            stroke-width="1"
          />
          <!-- Y-axis labels -->
          <text
            v-for="line in gridLines"
            :key="line.value"
            :x="marginLeft - 10"
            :y="line.y + marginTop + 4"
            text-anchor="end"
            fill="hsl(var(--foreground))"
            font-size="12"
          >
            {{ yFormatter ? yFormatter(line.value, 0) : Math.round(line.value) }}
          </text>
        </g>

        <!-- X-axis -->
        <g v-if="showXAxis">
          <line
            :x1="marginLeft"
            :y1="marginTop + innerHeight"
            :x2="marginLeft + innerWidth"
            :y2="marginTop + innerHeight"
            stroke="hsl(var(--border))"
            stroke-width="1"
          />
          <!-- X-axis labels -->
          <text
            v-for="(dataItem, dataIndex) in data"
            :key="dataIndex"
            :x="marginLeft + xScale(dataIndex)"
            :y="marginTop + innerHeight + 20"
            text-anchor="middle"
            fill="hsl(var(--foreground))"
            font-size="12"
          >
            {{ xFormatter ? xFormatter(dataItem[index], dataIndex) : dataItem[index] }}
          </text>
        </g>

        <!-- Bars -->
        <g>
          <template v-for="(dataItem, dataIndex) in data" :key="dataIndex">
            <rect
              v-for="(category, categoryIndex) in categories"
              :key="`${dataIndex}-${categoryIndex}`"
              :x="marginLeft + xScale(dataIndex) - (categories.length * barWidth / 2) + (categoryIndex * barWidth)"
              :y="marginTop + yScale(dataItem[category] as number)"
              :width="barWidth"
              :height="innerHeight - yScale(dataItem[category] as number)"
              :fill="colors[categoryIndex]"
              :opacity="legendItems.find(item => item.name === category)?.inactive ? filterOpacity : 
                (hoveredBar && hoveredBar.dataIndex === dataIndex && hoveredBar.categoryIndex === categoryIndex) ? 0.8 : 1"
              :rx="roundedBars ? 2 : 0"
              :ry="roundedBars ? 2 : 0"
              class="cursor-pointer transition-opacity duration-200"
              @click="handleBarClick(dataItem, categoryIndex, dataIndex)"
              @mouseenter="handleBarHover($event, dataItem, categoryIndex, dataIndex)"
              @mouseleave="handleBarLeave"
            />
          </template>
        </g>
      </svg>

      <!-- Tooltip -->
      <div
        v-if="showTooltip && tooltipData"
        class="absolute pointer-events-none z-10 rounded-lg border bg-background p-2 shadow-lg whitespace-nowrap"
        :style="{
          left: tooltipData.x + 'px',
          top: tooltipData.y + 'px',
          transform: 'translate(-50%, -100%)'
        }"
      >
        <div class="text-sm">
          <div class="font-medium">{{ tooltipData.data.category }}</div>
          <div class="text-muted-foreground">{{ Number(tooltipData.data.value).toFixed(2) }} kWh</div>
        </div>
      </div>
    </div>
  </div>
</template>
