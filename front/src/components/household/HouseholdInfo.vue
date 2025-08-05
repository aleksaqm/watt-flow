<script setup lang="ts">
import Card from '@/shad/components/ui/card/Card.vue';
import CardContent from '@/shad/components/ui/card/CardContent.vue';
import CardHeader from '@/shad/components/ui/card/CardHeader.vue';
import CardTitle from '@/shad/components/ui/card/CardTitle.vue';
import { useRoute } from 'vue-router';
import type { HouseholdFull } from './household';
import { ref } from 'vue';
import { computed } from 'vue';
import { Button } from '@/shad/components/ui/button';
import "leaflet/dist/leaflet.css";
import { LMap, LTileLayer, LMarker, LPopup } from "@vue-leaflet/vue-leaflet";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
  DialogClose
} from '@/shad/components/ui/dialog';
import HouseholdAccessManager from '@/components/household/HouseholdAccessManager.vue';
import { getUsernameFromToken } from '@/utils/jwtDecoder';



const props = defineProps<{
  household: HouseholdFull
}>()

const center = computed<[number, number]>(() => {
  return [props.household.latitude, props.household.longitude]
})

const zoom = ref(15);

const dialogKey = ref(0);


const handleAccessChange = () => {
};

const isOwner = computed(() => {
  return getUsernameFromToken() === props.household.owner_name;
})

</script>
<template>
  <div class="flex flex-col my-10 gap-3 items-center justify-center">
    <Card>
      <CardHeader>
        <CardTitle>
          Household:
          <span class="mx-4 text-indigo-500">{{ props.household.cadastral_number }}</span>
        </CardTitle>
      </CardHeader>
    </Card>

    <div class="flex justify-center gap-20 my-20 items-center max-w-full w-full flex-wrap">

      <Card class="w-2/5 min-w-fit">
        <CardHeader>
          <CardTitle>
            <span class="text-gray-600 text-xl">Information</span>
          </CardTitle>
          <CardContent class="mt-5 flex flex-col gap-2">
            <div class="info-member">
              <span class="info-name text-gray-800">Owner:</span>
              <div class="info-value font-semibold">{{ props.household.owner_name == "" ? "unowned" :
                props.household.owner_name }}</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Address:</span>
              <div class="info-value font-semibold">{{ props.household.street }}</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">City:</span>
              <div class="info-value font-semibold">{{ props.household.city }}</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Status:</span>
              <div class="info-value font-semibold">{{ props.household.status }}</div>
            </div>
            <div class="info-member">
              <span class="info-name text-gray-800">Floor:</span>
              <div class="info-value font-semibold">{{ props.household.floor }}</div>
            </div>
            <div class="info-member">
              <span class="info-name text-gray-800">Suite:</span>
              <div class="info-value font-semibold">{{ props.household.suite }}</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Power meter:</span>
              <div class="info-value font-semibold text-xs">{{ props.household.device_address }}</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Square footage:</span>
              <div class="info-value font-semibold">{{ props.household.sq_footage }} m<sup>2</sup></div>
            </div>
          </CardContent>
        </CardHeader>

        <div class="flex justify-center pb-8" v-if="isOwner">
          <Dialog :key="dialogKey" >
            <DialogTrigger>
              <Button variant="outline" class="border-2">Authorize access</Button>
            </DialogTrigger>
            <DialogContent >
              <DialogHeader>
                <DialogTitle>Give access to this household</DialogTitle>
              </DialogHeader>
              <HouseholdAccessManager 
                :householdId="household.id" 
                @authorized="handleAccessChange" 
              />
            </DialogContent>
          </Dialog>
        </div>

      </Card>

      <Card class="min-w-fit w-1/2">
        <CardHeader>
          <CardTitle>
            <span class="text-gray-600 text-xl">Location</span>
          </CardTitle>
          <CardContent class="mt-5 flex flex-col gap-2">

            <div class="w-full h-64 sm:h-96 z-[1]">
              <l-map ref="map" v-model:zoom="zoom" :center="center" :use-global-leaflet="false" class="w-full h-full" >
                <l-tile-layer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" layer-type="base"
                  name="OpenStreetMap"></l-tile-layer>
                <l-marker :lat-lng="center">
                  <l-popup>{{ props.household.street + ", " + props.household.number }} </l-popup>

                </l-marker>
              </l-map>
            </div>
          </CardContent>
        </CardHeader>

      </Card>
    </div>

  </div>


</template>

<style scoped>
.info-member {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: start;
  gap: 10px;
  color: #474747;
}

.info-name {
  width: 30%;
}

.info-value {
  width: 70%;
  border: 1px solid #adadad;
  padding: 4px;
  padding-left: 20px;
  padding-right: 20px;
  border-radius: 7px;
}
</style>
