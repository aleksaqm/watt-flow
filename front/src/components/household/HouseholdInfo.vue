<script setup lang="ts">
import Card from '@/shad/components/ui/card/Card.vue';
import CardContent from '@/shad/components/ui/card/CardContent.vue';
import CardHeader from '@/shad/components/ui/card/CardHeader.vue';
import CardTitle from '@/shad/components/ui/card/CardTitle.vue';

</script>
<template>
  <div class="flex flex-col my-10 gap-3 items-center justify-center">
    <Card>
      <CardHeader>
        <CardTitle>
          Household:
          <span class="mx-4 text-indigo-500">1232131</span>
        </CardTitle>
      </CardHeader>
    </Card>

    <div class="flex justify-center gap-20 my-20 items-center max-w-full w-full flex-wrap">

      <Card class="w-2/5 min-w-fit">
        <CardHeader>
          <CardTitle>
            <span class="text-gray-600 text-xl">Information</span>
          </CardTitle>
          <CardContent class="mt-5 flex flex-col gap-2">
            <div class="info-member">
              <span class="info-name text-gray-800">Owner:</span>
              <div class="info-value font-semibold">none</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Address:</span>
              <div class="info-value font-semibold">Bulevar oslobodjenja 12</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">City:</span>
              <div class="info-value font-semibold">Novi Sad</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Status:</span>
              <div class="info-value font-semibold">Unowned</div>
            </div>
            <div class="info-member">
              <span class="info-name text-gray-800">Floor:</span>
              <div class="info-value font-semibold">2</div>
            </div>
            <div class="info-member">
              <span class="info-name text-gray-800">Suite:</span>
              <div class="info-value font-semibold">21</div>
            </div>

            <div class="info-member">
              <span class="info-name text-gray-800">Square footage:</span>
              <div class="info-value font-semibold">89 m<sup>2</sup></div>
            </div>
          </CardContent>
        </CardHeader>

      </Card>

      <Card class="w-2/5 min-w-fit h-96">
        <CardHeader>
          <CardTitle>
            <span class="text-gray-600 text-xl">Location</span>
          </CardTitle>
          <CardContent class="mt-5 flex flex-col gap-2">
          </CardContent>
        </CardHeader>

      </Card>
    </div>

  </div>

</template>

<style scoped>
.info-member {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: start;
  gap: 10px;
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
