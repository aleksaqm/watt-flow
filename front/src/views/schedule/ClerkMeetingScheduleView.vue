<script setup lang="ts">
import DatePicker from '@/components/schedule/DatePicker.vue';
import Meeting from '@/components/schedule/Meeting.vue';
import { useUserStore } from '@/stores/user';
import axios from 'axios';
import { ref } from 'vue';

const meeting = ref<{ User?: string; Date?: string; Time?: string; Duration: string; }>({ Duration: "0" })
const updateMeeting = async (meetingId: number) => {
  axios.get("/api/meeting/" + meetingId).then((result) => {
    const date = new Date(result.data.data.start_time)
    console.log(date)
    meeting.value = {
      User: result.data.data.username,
      Date: date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDate(),
      Time: formattedTime(date.getHours(), date.getMinutes()),
      Duration: String(result.data.data.duration)
    }
  }).catch(error => console.log(error))

  console.log(meeting.value)

}

const formattedTime = (hour: number, minute: number) =>
  `${String(hour).padStart(2, "0")}:${String(minute).padStart(2, "0")}`

const userStore = useUserStore();
</script>



<template>
  <main>
    <div class="flex flex-col items-center justify-center ">
      <span class="text-xl m-10">My meetings</span>
      <DatePicker :user-id="userStore.id" :username="null" @meeting-id="updateMeeting"></DatePicker>
      <Meeting v-if="meeting.User" :meeting="meeting"></Meeting>
    </div>
  </main>
</template>

<style></style>
