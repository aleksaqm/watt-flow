<script setup lang="ts">
import * as z from 'zod'
import { toTypedSchema } from '@vee-validate/zod'
import { Field, useForm, useField } from 'vee-validate'
import { ref, computed, nextTick, onMounted } from 'vue'
import { useToast } from '../../shad/components/ui/toast/use-toast'

import Toaster from '@/shad/components/ui/toast/Toaster.vue';
import { Button } from '@/shad/components/ui/button'
import { Input } from '@/shad/components/ui/input'
import { Label } from '@/shad/components/ui/label'
import { ZodIssueCode } from "zod";
import axios from 'axios'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shad/components/ui/popover'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shad/components/ui/command'
import "leaflet/dist/leaflet.css";
import { LMap, LTileLayer, LMarker, LPopup } from "@vue-leaflet/vue-leaflet";
import 'leaflet/dist/leaflet.css';
import { useUserStore } from '@/stores/user';


const formSchema = toTypedSchema(
  z.object({
    street: z.string(),
    streetNumber: z.string(),
    city: z.string().nonempty({ message: "City is required" }),
    numberOfFloors: z.number().min(1, { message: "At least one floor is required" }),
    lon: z.number(),
    lat: z.number(),
    households: z.array(
      z.object({
        floor: z
          .string()
          .transform((val) => parseInt(val, 10))
          .refine(
            (val) => !isNaN(val) && val > 0,
            { message: "Floor must be a positive integer", path: ["floor"] }
          ),
        suite: z.string().refine((val) => val.trim() !== "", { message: "Suite is required" }),
        area: z.number().positive({ message: "Area must be a positive number" }),
        identifier: z.string().refine((val) => val.trim() !== "", { message: "Cadastral number is required" }),
      })
    ),
    propertyImages: z
      .array(
        z.instanceof(File).refine(
          (file) => file.type.startsWith("image/"),
          { message: "Only image files are allowed" }
        )
      )
      .min(1, { message: "At least one property image is required" })
      .max(10, { message: "You can upload up to 10 property images" }),
    documents: z
      .array(
        z.instanceof(File).refine(
          (file) => file.type === "application/pdf",
          { message: "Only PDF files are allowed" }
        )
      )
      .min(1, { message: "At least one document is required" })
      .max(5, { message: "You can upload up to 5 documents" }),
    }).superRefine((data, ctx) => {
    data.households.forEach((household, index) => {
      if (household.floor > data.numberOfFloors) {
        ctx.addIssue({
          code: ZodIssueCode.custom,
          path: ["households", index, "floor"],
          message: "Floor must be less than or equal to the number of floors",
        });
      }
    });

    const identifiers = data.households.map((household) => household.identifier);
    const duplicates = identifiers.filter((id, index) => identifiers.indexOf(id) !== index);

    if (duplicates.length > 0) {
      duplicates.forEach((duplicateId) => {
        const duplicateIndex = identifiers.indexOf(duplicateId);
        ctx.addIssue({
          code: ZodIssueCode.custom,
          path: ["households", duplicateIndex, "identifier"],
          message: "Identifier must be unique",
        });
      });
    }

    const suites = data.households.map((household) => household.suite);
    const suiteDuplicates = suites.filter((suite, index) => suites.indexOf(suite) !== index);
    if (suiteDuplicates.length > 0) {
      suiteDuplicates.forEach((duplicateSuite) => {
        const duplicateIndex = suites.indexOf(duplicateSuite);
        ctx.addIssue({
          code: ZodIssueCode.custom,
          path: ["households", duplicateIndex, "suite"],
          message: "Suite must be unique",
        });
      });
    }
  })
);

const { handleSubmit, errors, setFieldValue } = useForm({
  validationSchema: formSchema,
});

const { toast } = useToast()

const street = ref('');
const city = ref('');
const streetNumber = ref('');
const numberOfFloors = ref<number>();
const propertyImages = ref<File[]>([]);
const documents = ref<File[]>([]);
const lon = ref(0.0)
const lat = ref(0.0)

const center = ref<[number, number]>([45.254242348610845, 19.843728661762427]);

const zoom = ref(6);
const markerPosition = ref<[number, number]>([45.254242348610845, 19.843728661762427]);

const householdEntries = ref<{ floor: string; suite: string; area: number; identifier: string }[]>([]);

const addHouseholdEntry = () => {
  householdEntries.value.push({
    floor: '',
    suite: '',
    area: 0,
    identifier: '',
  });
};

const cityLabel = computed(() => {
  return city.value ? city.value : "Select city..."
});

const removeHouseholdEntry = (index: number) => {
  householdEntries.value.splice(index, 1);
};

const selectCity = async (selectedValue: string) => {
  city.value = city.value === selectedValue ? '' : selectedValue;
  console.log(city.value);
  setFieldValue('city', city.value)
  if (city.value) {
    const geocodeResult = await geocodeCity(city.value);
    console.log(geocodeResult);
    if (geocodeResult) {
      const { lat, lon } = geocodeResult;
      markerPosition.value = [lat, lon];
      zoom.value = 12;
      center.value = [lat, lon]
    }
  }

};

const geocodeCity = async (city: string) => {
  const url = `https://nominatim.openstreetmap.org/search?city=${encodeURIComponent(city)}&format=json&limit=1`;

  try {
    const response = await fetch(url);
    const data = await response.json();

    if (data && data.length > 0) {
      const { lat, lon } = data[0];
      return { lat, lon };
    } else {
      alert("City not found!");
      return null;
    }
  } catch (error) {
    console.error("Error fetching geocoding data:", error);
    alert("Error fetching city data.");
    return null;
  }
};

const onFileChange = (event: any, type: 'propertyImages' | 'documents') => {
  const target = event.target as HTMLInputElement;
  const files = target.files ? Array.from(target.files) : [];

  if (type === 'propertyImages') {
    propertyImages.value = files;
    setFieldValue('propertyImages', propertyImages.value)
  } else {
    documents.value = files;
    setFieldValue('documents', documents.value)
  }
};



const onMapClick = async (event: { latlng: { lat: number, lng: number } }) => {
  console.log(errors)
  markerPosition.value = [event.latlng.lat, event.latlng.lng];
  center.value = [event.latlng.lat, event.latlng.lng];
  console.log('Latitude:', event.latlng.lat);
  console.log('Longitude:', event.latlng.lng);

  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?lat=${event.latlng.lat}&lon=${event.latlng.lng}&format=json&accept-language=sr-Latn`
    );

    const data = await response.json();
    const fetchedStreet = data.address.road;
    const fetchedStreetNumber = data.address.house_number;
    setFieldValue('street', street.value);
    setFieldValue('streetNumber', streetNumber.value);
    setFieldValue('lon', event.latlng.lng)
    setFieldValue('lat', event.latlng.lat)
    street.value = fetchedStreet;
    streetNumber.value = fetchedStreetNumber;
    lon.value = event.latlng.lng;
    lat.value = event.latlng.lat;
    nextTick(() => {
      console.log('Address:', street.value, streetNumber.value);
    });
  } catch (error) {
    console.error('Error fetching address:', error);
  }
};

const convertToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => {
      if (reader.result) {
        resolve(reader.result as string)
      } else {
        reject('Failed to convert file to base64')
      }
    }
    reader.onerror = () => reject('Failed to read file')
    reader.readAsDataURL(file)
  })
}


const submitForm = async () => {
  try {
    if (householdEntries.value.length == 0){
        toast({
        title: 'Property Registration Error',
        description: "You can register property without households.",
        variant: 'destructive',
      });
    }

    const imagesBase64 = await Promise.all(
      propertyImages.value.map((file) => convertToBase64(file))
    );

    const documentsBase64 = await Promise.all(
      documents.value.map((file) => convertToBase64(file))
    );
    const userStore = useUserStore()

    const data = {
      floors: numberOfFloors.value,
      status: 0,
      owner_id: userStore.id,
      images: imagesBase64,
      documents: documentsBase64,
      household: householdEntries.value.map(entry => ({
        floor: parseInt(entry.floor, 10),
        suite: entry.suite,
        status: 0,
        sq_footage: entry.area,
        cadastral_number: entry.identifier,
        device_status: {
          device_id: crypto.randomUUID(),
          is_active: false
        },
      })),

      address: {
        city: city.value,
        street: street.value,
        number: streetNumber.value,
        latitude: lat.value,
        longitude: lon.value,
      }
    };

    console.log(data);
    const response = await axios.post('/api/property', data);
    console.log(response);

    toast({
      title: 'Property Registration Successful',
      description: 'Your request is now under review.',
      variant: 'default',
    });
  } catch (error) {
    let errorMessage = 'An unexpected error occurred';

    if (axios.isAxiosError(error)) {
      const serverMessage = error.response?.data?.error || '';

      if (serverMessage.includes('duplicate key value') || serverMessage.includes('unique constraint')) {
        errorMessage = 'Cadastral number already exists in database. Please check it again.';
      } else if (serverMessage) {
        errorMessage = serverMessage;
      } else {
        errorMessage = `Server error: ${error.response?.status}`;
      }
    } else if (error instanceof Error) {
      errorMessage = `Request setup error: ${error.message}`;
    } else {
      errorMessage = 'An unexpected error occurred';
    }
    toast({
      title: 'Property Registration Error',
      description: errorMessage,
      variant: 'destructive',
    });
  }
};

const cities = ref<string[]>([]);

const fetchCities = async () => {
  try{
    const response = await axios.get('/api/cities')
    cities.value = response.data?.data || [];
    console.log(cities)
  }catch (error){
    console.log(error)
  }
}


onMounted(fetchCities)

const onSubmit = handleSubmit(submitForm);
</script>

<template>
  <div class="w-2/4 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <div class="flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-lg">Property Registration</span>
      <form class="w-full space-y-6" @submit="onSubmit">

        <div class="w-full h-64 sm:h-96">
          <l-map ref="map" v-model:zoom="zoom" :center="center" @click="onMapClick" :use-global-leaflet="false"
            class="w-full h-full">
            <l-tile-layer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" layer-type="base"
              name="OpenStreetMap"></l-tile-layer>
            <l-marker :lat-lng="markerPosition">
              <l-popup>Marker here</l-popup>

            </l-marker>
          </l-map>
        </div>

        <div style="display: flex; gap: 1rem;">
          <FormField name="street" v-slot="{ field }">
            <FormItem style="flex: 1; position: relative;">
              <FormLabel>Street</FormLabel>
              <FormControl>
                <Input type="text" disabled="true" v-bind="field" placeholder="Select address on map" v-model="street" />
              </FormControl>
              <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.street">{{ errors.street }}</FormMessage>
            </FormItem>
          </FormField>
          <FormField name="streetNumber" v-slot="{ field }">
            <FormItem style="flex: 1;  position: relative;">
              <FormLabel>Street number</FormLabel>
              <FormControl>
                <Input type="text" v-bind="field" disabled="true" v-model="streetNumber" placeholder="Select address on map" />
              </FormControl>
              <FormMessage class="absolute -bottom-5 left-0 text-xs"  v-if="errors.streetNumber">{{ errors.streetNumber }}</FormMessage>
            </FormItem>
          </FormField>
        </div>

        <FormField name="city">
          <FormItem style="flex: 1; position: relative;">
            <FormLabel>City</FormLabel>
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" class="w-full justify-between">
                  {{ cityLabel }}
                </Button>
              </PopoverTrigger>
              <PopoverContent class="w-full p-0">
                <Command>
                  <CommandInput placeholder="Search city..." />
                  <CommandList>
                    <CommandEmpty>No city found.</CommandEmpty>
                    <CommandGroup>
                      <CommandItem 
                        v-for="cityItem in cities" 
                        :key="cityItem" 
                        :value="cityItem" 
                        @select="selectCity(cityItem)">
                        {{ cityItem }}
                      </CommandItem>
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
            <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.city">{{ errors.city }}</FormMessage>
          </FormItem>
        </FormField>


        <FormField name="numberOfFloors" v-slot="{ field }">
          <FormItem style="flex: 1; position: relative;">
            <FormLabel>Number of Floors</FormLabel>
            <FormControl>
              <Input type="number" v-bind="field" v-model="numberOfFloors" placeholder="Enter number of floors" />
            </FormControl>
            <FormMessage  class="absolute -bottom-5 left-0 text-xs"  v-if="errors.numberOfFloors">{{ errors.numberOfFloors }}</FormMessage>
          </FormItem>
        </FormField>

        <FormField  class="pt-2" name="propertyImages">
          <FormItem  style="position: relative;">
            <FormLabel>Property Images: </FormLabel>
            <FormControl>
              <input type="file" @change="e => onFileChange(e, 'propertyImages')" accept="image/*" multiple />
            </FormControl>
            <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.propertyImages">{{ errors.propertyImages }}</FormMessage>
          </FormItem>
          
        </FormField>

        <FormField name="documents">
          <FormItem style="position: relative;">
            <FormLabel>Documents (PDFs): </FormLabel>
            <FormControl>
              <input type="file" @change="e => onFileChange(e, 'documents')" accept="application/pdf" multiple />
            </FormControl>
            
            <FormMessage class="absolute -bottom-5 left-0 text-xs"  v-if="errors ">{{ errors  }}</FormMessage>
          </FormItem>
        </FormField>

        
        <div v-for="(household, index) in householdEntries" :key="index" class="space-y-4">
          <div class="flex gap-4 items-center">
            <FormField :name="`households.${index}.floor`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Floor</FormLabel>
                <FormControl>
                  <Input type="text" v-bind="field" v-model="household.floor" placeholder="Enter floor" />
                </FormControl>
                <FormMessage  class="absolute -bottom-5 left-0 text-xs" >{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.suite`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Suite</FormLabel>
                <FormControl>
                  <Input type="text" v-bind="field" v-model="household.suite" placeholder="Enter suite" />
                </FormControl>
                <FormMessage  class="absolute -bottom-5 left-0 text-xs" >{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.area`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Area (m²)</FormLabel>
                <FormControl>
                  <Input type="number" v-bind="field" v-model="household.area" placeholder="Enter area in m²" />
                </FormControl>
                <FormMessage  class="absolute -bottom-5 left-0 text-xs" >{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.identifier`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Cadastral num</FormLabel>
                <FormControl>
                  <Input type="text" v-bind="field" v-model="household.identifier"
                    placeholder="Enter cadastral number" />
                </FormControl>
                <FormMessage  class="absolute -bottom-5 left-0 text-xs" >{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <Button type="button" @click="removeHouseholdEntry(index)" class="bg-red-500 text-white">
              Remove Household
            </Button>
          </div>
        </div>
        
        <div class="flex items-center gap-5 text-red-500">
          <Button type="button" @click="addHouseholdEntry" class="bg-indigo-500 text-white hover:bg-gray-600">Add
          Household</Button>
          <span v-if="errors.households">{{ errors.households }}</span>
        </div>
        
        
        <Button type="submit" class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2">
          Submit Request
        </Button>
      </form>
    </div>
  </div>
  <Toaster />
</template>

<style scoped></style>
