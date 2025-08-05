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
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/shad/components/ui/alert-dialog'
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
import router from '@/router'


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
  validateOnMount: false,
  validationSchema: formSchema,
});

const { toast } = useToast()

const street = ref('');
const city = ref('');
const streetNumber = ref('');
const numberOfFloors = ref<number>();
const propertyImages = ref<File[]>([]);
const propertyImagePreviews = ref<string[]>([]);
const documents = ref<File[]>([]);
const documentPreviews = ref<string[]>([]);
const lon = ref(0.0)
const lat = ref(0.0)

const center = ref<[number, number]>([45.254242348610845, 19.843728661762427]);

const zoom = ref(6);
const markerPosition = ref<[number, number]>([45.254242348610845, 19.843728661762427]);

// Confirmation dialog state
const showConfirmDialog = ref(false);

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

const selectCity = async (selectedValue: string, shouldGeocode: boolean = true) => {
  city.value = selectedValue;
  setFieldValue('city', city.value)
  if (city.value && shouldGeocode) {
    const geocodeResult = await geocodeCity(city.value);
    setFieldValue("street", "")
    setFieldValue("streetNumber", "")
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

const onFileChange = (event: Event, type: 'propertyImages' | 'documents') => {
  const target = event.target as HTMLInputElement;
  if (!target.files) return;

  const newFiles = Array.from(target.files);

  if (type === 'propertyImages') {
    propertyImages.value.push(...newFiles);
    propertyImagePreviews.value = propertyImages.value.map(file => URL.createObjectURL(file));
    setFieldValue('propertyImages', propertyImages.value);

  } else if (type === 'documents') {
    documents.value.push(...newFiles);
    documentPreviews.value = documents.value.map(file => file.name);
    
    setFieldValue('documents', documents.value);
  }

  target.value = '';
};

const removeFile = (index: number, type: 'propertyImages' | 'documents') => {
  if (type === 'propertyImages') {
    propertyImages.value.splice(index, 1);
    propertyImagePreviews.value.splice(index, 1);
    
    setFieldValue('propertyImages', propertyImages.value);

  } else if (type === 'documents') {
    documents.value.splice(index, 1);
    documentPreviews.value.splice(index, 1);
    
    setFieldValue('documents', documents.value);
  }
};

const openFile = (file: File) => {
  const url = URL.createObjectURL(file);
  window.open(url, '_blank');
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
    setFieldValue('street', fetchedStreet);
    setFieldValue('streetNumber', fetchedStreetNumber);
    setFieldValue('lon', event.latlng.lng)
    setFieldValue('lat', event.latlng.lat)
    street.value = fetchedStreet;
    streetNumber.value = fetchedStreetNumber;
    lon.value = event.latlng.lng;
    lat.value = event.latlng.lat;
    console.log("DATA: ", data)
    const potentialCitiesFromApi = [
      data.address.city,
      data.address.town,
      data.address.village,
      data.address.municipality,
    ].filter(Boolean);

    let matchedCity: string | undefined = undefined;

    const sortedCities = [...cities.value].sort((a, b) => b.length - a.length);

    for (const potentialCity of potentialCitiesFromApi) {
      const foundInList = sortedCities.find(
        (c) => potentialCity.toLowerCase().includes(c.toLowerCase())
      );
      
      if (foundInList) {
        matchedCity = foundInList;
        break; 
      }
    }
    if (matchedCity) {
      selectCity(matchedCity, false);
    } else {
      console.warn(`City "${potentialCitiesFromApi.join(', ')}" from API not found in the predefined city list.`);
      selectCity(''); 
    }

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
    if (householdEntries.value.length == 0) {
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

    const response = await axios.post('/api/property', data);

    toast({
      title: 'Property Registration Successful',
      description: 'Your property request has been submitted successfully and will be reviewed by administrators.',
      variant: 'default',
    });

    // Small delay to ensure toast is shown and state is reset before navigation
    setTimeout(() => {
      router.replace("/my-property-request");
    }, 100);
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
  try {
    const response = await axios.get('/api/cities')
    cities.value = response.data?.data || [];
  } catch (error) {
    console.log(error)
  }
}


onMounted(fetchCities)

const showConfirmationDialog = () => {
  showConfirmDialog.value = true;
};

const confirmSubmit = async () => {
  showConfirmDialog.value = false;
  await submitForm();
};

const onSubmit = handleSubmit(showConfirmationDialog);
</script>

<template>
  <div class="w-2/4 p-7 flex flex-col justify-center items-center bg-white shadow-lg">
    <div class="flex flex-col justify-center items-center gap-5 w-full">
      <span class="text-gray-800 text-lg">Property Registration</span>
      <form class="w-full space-y-6" @submit="onSubmit">

        <div class="w-full h-64 sm:h-96 ">
          <l-map ref="map" v-model:zoom="zoom" :center="center" @click="onMapClick" :use-global-leaflet="false"
            class="w-full h-full z-[1]">
            <l-tile-layer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" layer-type="base"
              name="OpenStreetMap"></l-tile-layer>
            <l-marker :lat-lng="markerPosition">
              <l-popup> {{ street }} {{ streetNumber }}</l-popup>

            </l-marker>
          </l-map>
        </div>

        <div style="display: flex; gap: 1rem;">
          <FormField name="street" v-slot="{ field }">
            <FormItem style="flex: 1; position: relative;">
              <FormLabel>Street</FormLabel>
              <FormControl>
                <Input type="text" disabled="true" v-bind="field" placeholder="Select address on map"
                  v-model="street" />
              </FormControl>
              <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.street">{{ errors.street }}
              </FormMessage>
            </FormItem>
          </FormField>
          <FormField name="streetNumber" v-slot="{ field }">
            <FormItem style="flex: 1;  position: relative;">
              <FormLabel>Street number</FormLabel>
              <FormControl>
                <Input type="text" v-bind="field" disabled="true" v-model="streetNumber"
                  placeholder="Select address on map" />
              </FormControl>
              <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.streetNumber">{{ errors.streetNumber
                }}</FormMessage>
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
              <PopoverContent class="w-full p-0 z-[1000]">
                <Command>
                  <CommandInput placeholder="Search city..." />
                  <CommandList>
                    <CommandEmpty>No city found.</CommandEmpty>
                    <CommandGroup>
                      <CommandItem v-for="cityItem in cities" :key="cityItem" :value="cityItem"
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
            <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.numberOfFloors">{{ errors.numberOfFloors
              }}</FormMessage>
          </FormItem>
        </FormField>

        <FormField class="pt-2" name="propertyImages">
          <FormItem style="position: relative;">
            <FormLabel>Property Images: </FormLabel>
            <FormControl>
              <input type="file" @change="e => onFileChange(e, 'propertyImages')" accept="image/*" multiple />
            </FormControl>
            <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors.propertyImages">{{ errors.propertyImages
              }}</FormMessage>
              <div v-if="propertyImagePreviews.length > 0" class="flex flex-wrap gap-4 mt-4">
                <div
                  v-for="(previewUrl, index) in propertyImagePreviews"
                  :key="index"
                  class="relative group"
                >
                  <img
                    :src="previewUrl"
                    alt="Image Preview"
                    class="w-24 h-24 object-cover rounded-md"
                  />
                  <button
                    type="button"
                    @click="removeFile(index, 'propertyImages')"
                    class="absolute -top-2 -right-2 bg-white rounded-full p-0.5 text-red-500 hover:text-red-700 hover:scale-110 transition-transform"
                    aria-label="Remove image"
                  >
                    X
                  </button>
                  <div 
                    class="absolute inset-2 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer rounded-md"
                    @click="openFile(propertyImages[index])"
                  >
                    <span class="text-white text-xs">Open</span>
                  </div>
                </div>
              </div>
          </FormItem>

        </FormField>

        <FormField name="documents">
          <FormItem style="position: relative;">
            <FormLabel>Documents (PDFs): </FormLabel>
            <FormControl>
              <input type="file" @change="e => onFileChange(e, 'documents')" accept="application/pdf" multiple />
            </FormControl>

            <FormMessage class="absolute -bottom-5 left-0 text-xs" v-if="errors">{{ errors }}</FormMessage>

            <div v-if="documentPreviews.length > 0" class="mt-4 space-y-2">
              <p class="text-sm font-medium">Choosen documents:</p>
              <ul class="list-disc list-inside">
                <li 
                  v-for="(docName, index) in documentPreviews" 
                  :key="index" 
                  class="text-sm text-gray-700"
                >
                  <button 
                    type="button" 
                    @click="openFile(documents[index])"
                    class="text-blue-600 hover:underline"
                  >
                    {{ docName }}
                  </button>
                  <button
                    type="button"
                    @click="removeFile(index, 'documents')"
                    class="ml-4 text-gray-400 hover:text-red-600"
                    aria-label="Remove document"
                  >
                    x
                  </button>
                </li>
              </ul>
            </div>
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
                <FormMessage class="mt-2 text-xs">{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.suite`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Suite</FormLabel>
                <FormControl>
                  <Input type="text" v-bind="field" v-model="household.suite" placeholder="Enter suite" />
                </FormControl>
                <FormMessage class="mt-2 text-xs">{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.area`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Area (m²)</FormLabel>
                <FormControl>
                  <Input type="number" v-bind="field" v-model="household.area" placeholder="Enter area in m²" />
                </FormControl>
                <FormMessage class="mt-2 text-xs">{{ errors[0] }}</FormMessage>
              </FormItem>
            </FormField>

            <FormField :name="`households.${index}.identifier`" v-slot="{ field, errors }">
              <FormItem style=" position: relative;">
                <FormLabel>Cadastral num</FormLabel>
                <FormControl>
                  <Input type="text" v-bind="field" v-model="household.identifier"
                    placeholder="Enter cadastral number" />
                </FormControl>
                <FormMessage class="mt-2 text-xs">{{ errors[0] }}</FormMessage>
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
  
  <!-- Submit Confirmation Dialog -->
  <AlertDialog :open="showConfirmDialog" @update:open="showConfirmDialog = $event">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Are you sure?</AlertDialogTitle>
        <AlertDialogDescription>
          Are you sure you want to submit this property request?
          <br><br>
          <strong>Property Details:</strong><br>
          • Address: {{ street }} {{ streetNumber }}, {{ city }}<br>
          • Number of Floors: {{ numberOfFloors }}<br>
          • Households: {{ householdEntries.length }}<br>
          • Property Images: {{ propertyImages.length }}<br>
          • Documents: {{ documents.length }}<br>
          <br>
          Once submitted, your request will be reviewed by administrators.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>No, let me review</AlertDialogCancel>
        <AlertDialogAction 
          @click="confirmSubmit"
          class="bg-indigo-500 hover:bg-gray-600"
        >
          Yes, submit request
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
  
  <Toaster />
</template>

<style scoped></style>
