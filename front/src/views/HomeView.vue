<script setup lang="ts">
import Button from '../shad/components/ui/button/Button.vue';
</script>

<template>
  <div class="min-h-screen flex flex-col items-center bg-gradient-to-r from-[#1d1e26] to-[#4d596a] text-white px-6 py-8">
    <!-- Header Section -->
    <header class="text-center mt-6">
      <h1 class="text-4xl font-bold mb-4">Welcome to Watt-Flow</h1>
      <p class="text-lg text-gray-300 max-w-2xl mx-auto">
        Track and optimize your household's energy consumption with ease. 
        Stay informed and take control of your electricity usage.
      </p>
    </header>

    <!-- Decorative Section -->
    <section class="mt-8 bg-white rounded-lg shadow-lg p-8 text-gray-800 w-full max-w-3xl">
      <h2 class="text-2xl font-semibold text-center mb-6">Why Choose Watt-Flow?</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div class="flex items-center space-x-4">
          <div class="p-3 rounded-full bg-[#1d1e26] text-white">
            ⚡
          </div>
          <p class="text-lg">Monitor real-time energy usage</p>
        </div>
        <div class="flex items-center space-x-4">
          <div class="p-3 rounded-full bg-[#1d1e26] text-white">
            💡
          </div>
          <p class="text-lg">Identify consumption patterns</p>
        </div>
        <div class="flex items-center space-x-4">
          <div class="p-3 rounded-full bg-[#1d1e26] text-white">
            📊
          </div>
          <p class="text-lg">Generate detailed usage reports</p>
        </div>
        <div class="flex items-center space-x-4">
          <div class="p-3 rounded-full bg-[#1d1e26] text-white">
            🌍
          </div>
          <p class="text-lg">Contribute to a sustainable future</p>
        </div>
      </div>
    </section>

  </div>
</template>
