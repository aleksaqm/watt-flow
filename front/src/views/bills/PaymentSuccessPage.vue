<script setup lang="ts">
import { useRoute } from 'vue-router';
import { CheckCircleIcon } from '@heroicons/vue/24/solid';
import Button from '@/shad/components/ui/button/Button.vue';
const route = useRoute();
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
    <div class="w-full max-w-lg mx-auto bg-white p-10 rounded-2xl shadow-sm text-center">
      <CheckCircleIcon class="h-20 w-20 text-green-500 mx-auto mb-4" />
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Payment Successfull!</h1>

      <p class="text-sm text-gray-500">
            The payment confirmation with the payment slip has been sent to your email address.
      </p>
      <div class="mt-8">
        <Button>
            <router-link to="/bills/owner" class="px-6 py-3">
                Show my bills
            </router-link>
        </Button>
        
      </div>
    </div>
  </div>
</template>