<script setup lang="ts">
import { Button } from '@/shad/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import { useUserStore } from '@/stores/user';
import { toTypedSchema } from '@vee-validate/zod';
import axios from 'axios';
import { useForm } from 'vee-validate';
import { ref } from 'vue';
import * as z from 'zod'
import { useToast } from '../../shad/components/ui/toast/use-toast'
import Toaster from '@/shad/components/ui/toast/Toaster.vue';


const props = defineProps<{ household: number }>();
const images = ref<File[]>([]);
const imagePreviews = ref<string[]>([]); // Added for image previews
const documents = ref<File[]>([]);
const documentPreviews = ref<string[]>([]); // Added for document previews
const { toast } = useToast();
const emit = defineEmits(['requestSent']);

const convertToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      if (reader.result) {
        resolve(reader.result as string);
      } else {
        reject('Failed to convert file to base64');
      }
    };
    reader.onerror = () => reject('Failed to read file');
    reader.readAsDataURL(file);
  });
};

const formSchema = toTypedSchema(
  z.object({
    images: z
      .array(
        z.instanceof(File).refine((file) => file.type.startsWith('image/'), {
          message: 'Only image files are allowed',
        })
      )
      .min(1, { message: 'At least one property image is required' })
      .max(10, { message: 'You can upload up to 10 property images' }),
    documents: z
      .array(
        z.instanceof(File).refine((file) => file.type === 'application/pdf', {
          message: 'Only PDF files are allowed',
        })
      )
      .min(1, { message: 'At least one document is required' })
      .max(5, { message: 'You can upload up to 5 documents' }),
  })
);

const { handleSubmit, errors, setFieldValue } = useForm({
  validationSchema: formSchema,
});

const onFileChange = (event: any, type: 'images' | 'documents') => {
  const target = event.target as HTMLInputElement;
  const files = target.files ? Array.from(target.files) : [];

  if (type === 'images') {
    images.value = files;
    setFieldValue('images', images.value);

    // Generate image previews
    imagePreviews.value = files.map((file) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      return URL.createObjectURL(file); // Temporary URL
    });
  } else {
    documents.value = files;
    setFieldValue('documents', documents.value);

    // Display document names
    documentPreviews.value = files.map((file) => file.name);
  }
};

const submitForm = async () => {
  try {
    const imagesBase64 = await Promise.all(
      images.value.map((file) => convertToBase64(file))
    );

    const documentsBase64 = await Promise.all(
      documents.value.map((file) => convertToBase64(file))
    );
    const userStore = useUserStore();

    const data = {
      images: imagesBase64,
      documents: documentsBase64,
      owner_id: userStore.id,
      household_id: props.household,
    };

    console.log(data);
    const response = await axios.post('/api/household/owner', data);
    console.log(response);
    emit('requestSent');
    toast({
      title: 'Request sent successfully!',
      description: 'Your request is now under review.',
      variant: 'default',
    });
  } catch (error) {
    let errorMessage = 'An unexpected error occurred';

    if (axios.isAxiosError(error)) {
      errorMessage = error.response?.data?.error || errorMessage;
    }
    toast({
      title: 'Request Error',
      description: errorMessage,
      variant: 'destructive',
    });
  }
};

const onSubmit = handleSubmit(submitForm);
</script>

<template>
  <form class="w-full space-y-6" @submit="onSubmit">
    <FormField class="pt-2" name="images">
      <FormItem style="position: relative;">
        <FormLabel>Property Images: </FormLabel>
        <FormControl>
          <input
            type="file"
            @change="e => onFileChange(e, 'images')"
            accept="image/*"
            multiple
          />
        </FormControl>
        <FormMessage
          class="absolute -bottom-5 left-0 text-xs"
          v-if="errors.images"
        >
          {{ errors.images }}
        </FormMessage>
        <div class="flex flex-wrap gap-4 mt-2">
          <img
            v-for="(preview, index) in imagePreviews"
            :key="index"
            :src="preview"
            alt="Image Preview"
            class="w-24 h-24 object-cover rounded"
          />
        </div>
      </FormItem>
    </FormField>

    <FormField name="documents">
      <FormItem style="position: relative;">
        <FormLabel>Documents (PDFs): </FormLabel>
        <FormControl>
          <input
            type="file"
            @change="e => onFileChange(e, 'documents')"
            accept="application/pdf"
            multiple
          />
        </FormControl>
        <FormMessage
          class="absolute -bottom-5 left-0 text-xs"
          v-if="errors.documents"
        >
          {{ errors.documents }}
        </FormMessage>
        <ul class="list-disc list-inside mt-2">
          <li v-for="(name, index) in documentPreviews" :key="index">
            {{ name }}
          </li>
        </ul>
      </FormItem>
    </FormField>

    <Button
      type="submit"
      class="w-full bg-gray-800 text-white hover:bg-gray-600 rounded-full py-2"
    >
      Submit
    </Button>
  </form>
  <Toaster />
</template>