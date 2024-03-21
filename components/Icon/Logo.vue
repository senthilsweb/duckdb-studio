<template>
  <img tag="img" to="#" alt="Logo" :src="logo_url ? logo_url : (usePropertiesStore().logo_url ? usePropertiesStore().logo_url : '/logo.svg')" :class="[usePropertiesStore().logo_css  ? usePropertiesStore().logo_css : 'h-10 mx-auto w-auto']"/>
</template>
<script setup>
import { usePropertiesStore } from '~/stores/properties';

//-----------------------------------------------------------------------------------------------------
//2. Page or component properties
//-----------------------------------------------------------------------------------------------------
const props = defineProps({
  logo_url: {
    type: String,
  },
});
</script>
