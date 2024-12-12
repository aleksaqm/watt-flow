<script setup lang="ts">
</script>

<template>
  <div class="flex-1 bg-white rounded-md p-4 flex flex-col w-3/5" id="1">
    <span class="text-gray-900 text-lg text-center my-5 w-full">Meeting</span>
    <div class="flex flex-col p-10 items-center justify-center border border-gray-199 rounded-sm w-full">
      <div class="info-member">
        <span class="info-name text-gray-800">Meeting with:</span>
        <div class="info-value font-semibold">{{ "danilo" }}</div>
      </div>
      <div class="info-member">
        <span class="info-name text-gray-800">Date:</span>
        <div class="info-value font-semibold">{{ "12-12-2022" }}</div>
      </div>
      <div class="info-member">
        <span class="info-name text-gray-800">Time:</span>
        <div class="info-value font-semibold">{{ "12:20" }}</div>
      </div>
      <div class="info-member">
        <span class="info-name text-gray-800">Duration:</span>
        <div class="info-value font-semibold">{{ "30 minutes" }}</div>
      </div>

      <div class="info-member">
        <span class="info-name text-gray-800">Status:</span>
        <div class="info-value font-semibold">{{ "Pending" }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.info-member {
  width: 70%;
  display: flex;
  align-items: center;
  justify-content: start;
  color: #474747;
  padding: 20px;
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
