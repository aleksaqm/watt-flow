<script setup lang="ts">

import RegisterLink from '@/components/RegisterLink.vue';
import SuperAdminChangePasswordForm from '@/components/SuperAdminChangePasswordForm.vue';
import SuperAdminWelcome from '@/components/SuperAdminWelcome.vue';
import axios from 'axios';
import { onBeforeMount, ref } from 'vue';
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(true);


onBeforeMount(async () => {
    try{
        const response = await axios.get("/api/admin/active")
        // console.log(response.data)
        if (response.data['data']){
            router.push({name: 'login'})
        }
        loading.value = false;
    }
    catch(error){
        router.push({name: 'login'})
    }
})

</script>

<template>
    <!-- <main v-if="!loading" class="flex-1 flex justify-center items-center">
        <h1>Super Admin, Go change your password</h1>
        <SuperAdminChangePasswordForm />
    </main> -->
    <main v-if="!loading" class="flex-1 flex items-center justify-center">
        <div class="flex w-full max-w-4xl shadow-lg">
            <SuperAdminWelcome />
            <SuperAdminChangePasswordForm />
        </div>
        <!-- <h1>Super Admin, Go change your password</h1> -->
         
    </main>
</template>

<style></style>