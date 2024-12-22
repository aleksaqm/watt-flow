<script setup lang="ts">
import type { CarouselApi } from '@/shad/components/ui/carousel'
import { Card, CardContent } from '@/shad/components/ui/card'
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from '@/shad/components/ui/carousel'
import { computed, onBeforeMount, onMounted, ref } from 'vue'
import { watchOnce } from '@vueuse/core';
import Button from '@/shad/components/ui/button/Button.vue';
import axios from 'axios';

const props = defineProps<{ images: string[]; documents: string[] }>()
const api = ref<CarouselApi>()
const totalCount = ref(0)
const current = ref(0)

const baseURL = "http://localhost:8080/"

const imageSources = computed(() => props.images.map(image => baseURL + image))
const documentSources = computed(() => props.documents.map(document => baseURL + document))
const isLoaded = ref<boolean>(false)

const images = ref<string[]>([])
const documents = ref<string[]>([])


function setApi(val: CarouselApi) {
  api.value = val
}

watchOnce(api, (api) => {
  if (!api) return

  totalCount.value = api.scrollSnapList().length
  current.value = api.selectedScrollSnap() + 1

  api.on('select', () => {
    current.value = api.selectedScrollSnap() + 1
  })
})

async function getDocs() {
  try {
    for (let index = 0; index < imageSources.value.length; index++) {
      const response = await axios.get(imageSources.value[index], { responseType: 'blob' });
      const blob = new Blob([response.data], { type: response.headers['content-type'] });
      images.value.push(URL.createObjectURL(blob));
    }
    for (let index = 0; index < documentSources.value.length; index++) {
      const response2 = await axios.get(documentSources.value[index], { responseType: 'blob' });
      const blob = new Blob([response2.data], { type: response2.headers['content-type'] });
      documents.value.push(URL.createObjectURL(blob))
    }
  } catch (error) {
    console.error("Fail")
  }
}

const openFile = (url: string) => {
  window.open(url, '_blank')
}

onBeforeMount(async () => {
  await getDocs()
  isLoaded.value = true
})
</script>

<template>
  <div class="flex flex-col items-center w-full sm:w-auto">
    <!-- Carousel -->
    <Carousel v-if="isLoaded" class="relative w-full max-w-xs" @init-api="setApi">
      <CarouselContent>
        <!-- Dynamically render each image -->
        <CarouselItem v-for="(src, index) in images" :key="index">
          <div class="p-1">
            <Card>
              <CardContent class="flex aspect-square items-center justify-center p-6">
                <img :src="src" alt="Carousel Image" class="object-cover w-full h-full rounded-md"
                  @click="openFile(src)" />
              </CardContent>
            </Card>
          </div>
        </CarouselItem>
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>

    <!-- Indicator below carousel -->
    <div class="py-2 text-center text-sm text-muted-foreground">
      Image {{ current }} of {{ totalCount }}
    </div>

    <!-- PDF Documents below the carousel -->
    <div v-if="isLoaded" class="mt-4 w-full max-w-xs">
      <div v-for="(doc, index) in documents" :key="index" class="mb-2">
        <Card>
          <CardContent class="flex items-center justify-between p-4">
            <span class="text-sm text-muted-foreground">Document {{ index + 1 }}</span>
            <Button variant="link" class="text-blue-500" @click="openFile(doc)">
              Open
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
