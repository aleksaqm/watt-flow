<script setup lang="ts">
import DatePicker from '@/components/schedule/DatePicker.vue';
import Meeting from '@/components/schedule/Meeting.vue';

</script>



<template>
  <main>
    <div class="flex flex-col items-center justify-center ">
      <span class="text-xl m-10">My meetings</span>
      <DatePicker></DatePicker>
      <Meeting></Meeting>
    </div>
  </main>
</template>

<style></style>
