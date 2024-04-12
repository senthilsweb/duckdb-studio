<template>
  <NuxtLayout name="landing">
    <div class="container p-4 lg:p-0 max-w-5xl mx-auto sm:grid grid-cols-12 gap-x-12">
      <div class="col-span-12 lg:col-span-8">
        <div class="mx-auto text-base max-w-prose lg:max-w-none">
          <TemplrJSTextOverlay :config="config" />
          <div class="mt-8 text-sm text-zinc-600 prose prose-zinc max-w-none dark:prose-invert dark:text-zinc-400 prose-headings:scroll-mt-28 prose-headings:font-display prose-headings:font-normal lg:prose-headings:scroll-mt-[8.5rem] prose-lead:text-zinc-500 dark:prose-lead:text-zinc-400 prose-a:font-semibold dark:prose-a:text-sky-400 prose-a:no-underline dark:prose-pre:ring-1 dark:prose-pre:ring-zinc-300/10 dark:prose-hr:border-zinc-800">
            <ContentRendererMarkdown :value="page.body" class="p-2" />
          </div>
        </div>
      </div>
      <div class="lg:col-span-3 sticky top-28 h-96 p-2 hidden lg:block justify-self-end">

        <Toc :data="page.toc.links" />
      </div>
    </div>
  </NuxtLayout>
</template>
<script setup>
import { ref, computed } from 'vue';
import { useHead, useState, useRoute, definePageMeta, parseMarkdown, useRuntimeConfig } from '#imports';

definePageMeta({});

const page = useState('page');
const data = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/blogs?title.eq=${useRoute().params.title}`);

if (data && Array.isArray(data.data)) {
  page.value = await parseMarkdown(data.data[0].content);
} else {
  page.value = { data: { title: 'Sorry, No page exists', body: '' }, toc: { links: [] } };
}

const config = ref({
  baseUrl: useRuntimeConfig().public.CLOUDINARY_BASE_URL,
  transformations: [
    {
      width: 800, // Default width
      crop: 'fit', // Default crop
      textOverlay: {
        font: 'Arial',
        size: 60,
        weight: 'bold',
        content: page.value.data.title.replace(/,/g, ' '), // Dynamically set content
        gravity: 'north_east',
        xOffset: 30,
        yOffset: 50,
      },
    },
  ],
  publicId: useRuntimeConfig().public.CLOUDINARY_PUBLIC_ID,
});

// Simplified transformation string generation, adhering to the model
function generateTransformationString(transformation) {
  const parts = [`w_${transformation.width || 800}`, `c_${transformation.crop || 'fit'}`];

  if (transformation.textOverlay) {
    const { textOverlay } = transformation;
    parts.push(`l_text:${textOverlay.font || 'Arial'}_${textOverlay.size || 60}_${textOverlay.weight || 'bold'}:${encodeURIComponent(textOverlay.content || 'Default Text')}`, `g_${textOverlay.gravity || 'south'}`, `x_${textOverlay.xOffset || 10}`, `y_${textOverlay.yOffset || 10}`);
  }

  return parts.join(',');
}

const cloudinaryURL = computed(() => {
  const transformationStrings = config.value.transformations.map(generateTransformationString).join('/');
  return `${config.value.baseUrl}/image/upload/${transformationStrings}/${config.value.publicId}`;
});

const siteDescription = 'SenthilsWeb - Data Engineering Hub: Simplifying Complex Technologies for Everyone';

useHead({
  title: page.value.data.title,
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1' },
    { name: 'description', content: siteDescription },
    { property: 'og:description', content: siteDescription },
    { property: 'og:title', content: page.value.data.title },
    { property: 'og:image', content: cloudinaryURL.value },
    { name: 'twitter:title', content: page.value.data.title },
    { name: 'twitter:image', content: cloudinaryURL.value },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:description', content: siteDescription },
  ],
});
</script>
