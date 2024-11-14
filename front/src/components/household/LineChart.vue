<script setup lang="ts">
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
)
const chartData = {
  labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
  datasets: [
    {
      label: 'Availability percentage',
      backgroundColor: '#ffffff',
      borderColor: '#6366F1',
      data: [40, 39, 10, 40, 39, 80, 40],
      tension: 0.1,
      fill: false,
    }
  ]
}

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
}
</script>

<template>
  <div class="m-10">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
