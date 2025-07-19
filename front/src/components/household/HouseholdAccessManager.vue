<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useDebounceFn } from '@vueuse/core';
import axios from 'axios';
import { Button } from '@/shad/components/ui/button'
import { Input } from '@/shad/components/ui/input'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/shad/components/ui/alert-dialog'

interface User {
  id: number;
  username: string;
  email: string;
}

const props = defineProps<{
  householdId: number;
}>();
const emit = defineEmits(['authorized']);

const usersWithAccess = ref<User[]>([]);
const isLoadingList = ref(true);

const searchInput = ref('');
const searchResults = ref<User[]>([]);
const isLoadingSearch = ref(false);

const errorMessage = ref<string | null>(null);



const fetchUsersWithAccess = async () => {
  isLoadingList.value = true;
  errorMessage.value = null;
  try {
    axios.get(`/api/household/access/${props.householdId}`).then((result) => {
        usersWithAccess.value = result.data.data || [];
    }).catch(error => { errorMessage.value = error.message;})
  } catch (error: any) {
    errorMessage.value = error.message;
  } finally {
    isLoadingList.value = false;
  }
};

const performSearch = useDebounceFn(async () => {
  if (searchInput.value.length < 2) {
    searchResults.value = [];
    return;
  }
  isLoadingSearch.value = true;
  errorMessage.value = null;
  try {
    const searchQuery = { Role: "Regular", Username: searchInput.value }
    const params = {
        sortBy: "username",
    }
    axios.post("/api/user/query", searchQuery, { params: params }).then((result) => {
        searchResults.value = result.data.users
    }).catch(error => { console.log(error); errorMessage.value = error.message; })
  } catch (error: any) {
    errorMessage.value = error.message;
  } finally {
    isLoadingSearch.value = false;
  }
}, 300);

watch(searchInput, performSearch);

const grantAccess = async (userId: number) => {
  errorMessage.value = null;
  try {
    const response = await axios.post(`/api/household/${props.householdId}/access`, 
      { userId: userId },
      {
        headers: { 'Content-Type': 'application/json' }
      }
    );
    
    searchInput.value = '';
    searchResults.value = [];
    await fetchUsersWithAccess();
    emit('authorized');
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || error.message;
  }
};

const revokeAccess = async (userId: number) => {
  errorMessage.value = null;
  try {
    await axios.delete(`/api/household/${props.householdId}/access/revoke/${userId}`);
    await fetchUsersWithAccess();
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || error.message;
  }
};

onMounted(() => {
  fetchUsersWithAccess();
});

</script>

<template>
  <div>
    <div v-if="errorMessage" class="p-3 mb-4 bg-red-100 text-red-700 border border-red-400 rounded">
      {{ errorMessage }}
    </div>

    <section class="mb-6">
      <h4 class="text-md font-semibold mb-3">Users with access</h4>
      
      <div v-if="isLoadingList" class="text-gray-500 text-sm">Loading...</div>
      <div v-else-if="usersWithAccess.length === 0" class="text-gray-500 text-sm">
        No one has access
      </div>
      <ul v-else class="space-y-2 max-h-48 overflow-y-auto pr-2">
        
        <li
          v-for="user in usersWithAccess"
          :key="user.id"
          class="flex items-center justify-between p-2 bg-gray-50 rounded-md border text-sm"
        >
          <div>
            <p class="font-medium">{{ user.username }}</p>
            <p class="text-xs text-gray-600">{{ user.email }}</p>
          </div>
          
          <AlertDialog>
            <AlertDialogTrigger>
                <Button
                    variant="destructive"
                    size="sm"
                >
                    Remove
                </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
            <AlertDialogHeader>
                <AlertDialogTitle>Remove {{ user.username }} access to your household?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction><Button @click="revokeAccess(user.id)">Remove</Button></AlertDialogAction>
            </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </li>
      </ul>
    </section>

    <section>
      <h4 class="text-md font-semibold mb-3">Authorize access</h4>
      <Input
        type="text"
        v-model="searchInput"
        placeholder="Search by username..."
      />

      <div v-if="isLoadingSearch" class="mt-3 text-gray-500 text-sm">Searching...</div>
      <ul v-if="searchResults.length > 0" class="mt-3 space-y-2 max-h-48 overflow-y-auto pr-2">
        <li
          v-for="user in searchResults"
          :key="user.id"
          class="flex items-center justify-between p-2 bg-blue-50 rounded-md border border-blue-200 text-sm"
        >
          <div>
            <p class="font-medium">{{ user.username }}</p>
            <p class="text-xs text-gray-600">{{ user.email }}</p>
          </div>
          
          <AlertDialog>
            <AlertDialogTrigger>
                <Button size="sm">
                    Authorize
                </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
            <AlertDialogHeader>
                <AlertDialogTitle>Authorize {{ user.username }} to access your household?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction><Button @click="grantAccess(user.id)">Authorize</Button></AlertDialogAction>
            </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </li>
      </ul>
      <div v-if="searchResults.length === 0 && searchInput.length > 1 && !isLoadingSearch" class="mt-3 text-gray-500 text-sm">
        No users found.
      </div>
    </section>
  </div>
</template>