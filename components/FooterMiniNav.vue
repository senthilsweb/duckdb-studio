<template>
  <div class="w-full flex justify-end space-x-4 md:order-2">
    <a v-for="item in filteredMenu" :key="item.name" :href="item.href" @click="useNuxtApp().$bus.$emit(item.click_event)" class="text-gray-400 hover:text-gray-500">
      <span class="sr-only">{{ item.name }}</span>
      <Icon :name="item.icon" class="h-6 w-6" :alt="item.name" />
    </a>
  </div>
</template>
<script setup>
import navigation from '@/public/footer-mini-nav.json';
const filteredMenu = computed(() => navigation.menu.filter((item) => item.enabled));
const props = defineProps({
  api_end_point: {
    type: String,
    required: true,
  },
});
</script>
