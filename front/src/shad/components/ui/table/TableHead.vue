<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import {ArrowDownIcon, ArrowUpIcon} from '@heroicons/vue/16/solid'


const props = defineProps<{
  class?: HTMLAttributes['class']
  orientation?: string
}>()
</script>

<template>
  <th :class="cn('h-12 px-4 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0', props.class)">
    <div class="flex">
      <slot />
      <span v-if="orientation == 'asc'"> <ArrowDownIcon class="w-5"/></span>
      <span v-else-if="orientation == 'desc'"><ArrowUpIcon class="w-5"/></span>
    </div>
    
  </th>
</template>
