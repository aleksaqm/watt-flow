// Make sure BarChart.vue exists in the same directory, or update the path below if it's located elsewhere
export { default as BarChart } from "../chart-bar/BarChart.vue"

export interface BaseChartProps<T extends Record<string, any>> {
  /**
   * The data for the chart.
   */
  data: T[]
  /**
   * The key in `data` to be used as the index.
   * This key will be used to map the data to the x-axis in cartesian charts.
   */
  index: keyof T
  /**
   * The keys in `data` to be used as the categories.
   * These keys will be used to map the data to the y-axis in cartesian charts.
   */
  categories: (keyof T)[]
  /**
   * Controls the visibility of the X-axis.
   * @default true
   */
  showXAxis?: boolean
  /**
   * Controls the visibility of the Y-axis.
   * @default true
   */
  showYAxis?: boolean
  /**
   * Controls the visibility of tooltip.
   * @default true
   */
  showTooltip?: boolean
  /**
   * Controls the visibility of the legend.
   * @default true
   */
  showLegend?: boolean
  /**
   * Controls the visibility of the grid.
   * @default true
   */
  showGridLine?: boolean
  /**
   * The margin for the chart container.
   * @default { top: 0, bottom: 0, left: 0, right: 0 }
   */
  margin?: {
    top?: number
    bottom?: number
    left?: number
    right?: number
  }
  /**
   * Opacity of the bars when they are set to inactive (via legend).
   * @default 0.2
   */
  filterOpacity?: number
  /**
   * An array of color values to be used for the chart bars.
   * If not provided, the default color palette will be used.
   */
  colors?: string[]
  /**
   * Format the values on the y-axis and inside tooltip
   */
  yFormatter?: (value: number, i: number) => string
  /**
   * Format the labels on the x-axis
   */
  xFormatter?: (value: any, i: number) => string
}
