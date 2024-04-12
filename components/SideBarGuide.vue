<!--
  This example requires some changes to your config:
  
  ```
  // tailwind.config.js
  module.exports = {
    // ...
    plugins: [
      // ...
      require('@tailwindcss/forms'),
    ],
  }
  ```
-->
<template>
  <TransitionRoot as="template" :show="isSideBarGuideOpen">
    <Dialog as="div" class="relative">
      <div class="fixed inset-0" />
      <div class="fixed inset-0 overflow-hidden">
        <div class="absolute inset-0 overflow-hidden">
          <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
            <TransitionChild as="template" enter="transform transition ease-in-out duration-500 sm:duration-700" enter-from="translate-x-full" enter-to="translate-x-0" leave="transform transition ease-in-out duration-500 sm:duration-700" leave-from="translate-x-0" leave-to="translate-x-full">
              <DialogPanel class="pointer-events-auto w-screen max-w-xl">
                <form class="flex h-full flex-col divide-y divide-gray-200 bg-white shadow-xl" @submit.prevent="saveSideBarGuide">
                  <div class="h-0 flex-1 overflow-y-auto">
                    <div class="bg-gradient-to-r from-primary-900 via-secondary-700 to-primary-900 py-6 px-4 sm:px-6">
                      <div class="flex items-center justify-between">
                        <DialogTitle class="text-lg font-medium text-white">{{ form_title ? form_title : 'Contact us' }} </DialogTitle>
                        <div class="ml-3 flex h-7 items-center">
                          <button type="button" class="rounded-md bg-primary-700 text-primary-200 hover:text-white focus:outline-none focus:ring-2 focus:ring-white" @click="isSideBarGuideOpen = !isSideBarGuideOpen">
                            <span class="sr-only">Close panel</span>
                            <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                          </button>
                        </div>
                      </div>
                      <div class="mt-1">
                        <p class="text-sm text-primary-300">{{ form_description ? form_description : 'Find out how we can help you?' }}</p>
                      </div>
                    </div>
                    <!--Progress bar (start)-->
                    <div v-if="progress" class="z-50 h-full overflow-hidden flex flex-col items-center justify-center">
                      <div class="loader animate-ping text-indigo-900 ease-linear rounded-full border-4 border-t-4 border-primary-900 h-12 w-12 mb-4"></div>
                    </div>
                    <!--Progress bar (end)-->
                    <!--Form Body (start)-->
                    <div class="flex-1 flex-col justify-between">
                      <div class="min-w-0 max-w-2xl flex-auto px-4 py-4 lg:max-w-none lg:pr-0 lg:pl-8 xl:px-8">
                        <ClientOnly>
                          <article>
                            <header class="mb-9 space-y-1">
                              <p class="font-display text-sm font-medium text-sky-500">{{ page.title }}</p>
                              <h1 class="font-display text-3xl tracking-tight text-gray-900 dark:text-white">{{ page.pageTitle || page.title }}</h1>
                            </header>
                            <div class="mt-8 text-sm text-zinc-600 prose prose-zinc max-w-none dark:prose-invert dark:text-zinc-400 prose-headings:scroll-mt-28 prose-headings:font-display prose-headings:font-normal lg:prose-headings:scroll-mt-[8.5rem] prose-lead:text-zinc-500 dark:prose-lead:text-zinc-400 prose-a:font-semibold dark:prose-a:text-sky-400 prose-a:no-underline dark:prose-pre:ring-1 dark:prose-pre:ring-zinc-300/10 dark:prose-hr:border-zinc-800">
                              <ContentRendererMarkdown :value="page" ref="doc_content"/>
                            </div>
                          </article>
                        </ClientOnly>
                      </div>
                      <!--<div class="hidden xl:sticky xl:top-[4.5rem] xl:-mr-6 xl:block xl:h-[calc(100vh-4.5rem)] xl:flex-none xl:overflow-y-auto xl:py-16">
                        <nav aria-labelledby="on-this-page-title" class="w-56">
                          <Toc :data="data.body.toc.links" />
                        </nav>
                      </div>-->
                    </div>
                    <!--Form Body (end)-->
                  </div>
                  <div class="flex flex-shrink-0 justify-end px-4 py-4">
                    <button type="button" class="border border-gray-300 bg-white py-2 px-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2" @click="cancel">Cancel</button>
                    <!--<button type="submit" class="ml-4 zyn-button">Submit</button>-->
                  </div>
                </form>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
<!--Server and client side javascript-->
<script setup>
import { parseMarkdown } from '@nuxtjs/mdc/runtime';
//-----------------------------------------------------------------------------------------------------
//1. Page level imports
//-----------------------------------------------------------------------------------------------------
import { ref } from 'vue';
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { XMarkIcon } from '@heroicons/vue/24/outline';

//-----------------------------------------------------------------------------------------------------
//2. Page or component properties
//-----------------------------------------------------------------------------------------------------
const props = defineProps({
  form_title: {
    type: String,
    required: true,
  },
  form_description: {
    type: String,
    required: true,
  },
  content_identifier: {
    type: String,
    required: true,
  },
});

//-----------------------------------------------------------------------------------------------------
//2. Variable declaration and initialization
//-----------------------------------------------------------------------------------------------------
const isSideBarGuideOpen = ref(false);
const progress = ref(false);

const page = ref({});

//-----------------------------------------------------------------------------------------------------
//3. Data initialization and server side rendering
//-----------------------------------------------------------------------------------------------------
/*const { data: inquiry_form_reasons } = await useAsyncData('inquiry_form_reasons-' + Math.random, () =>
  $fetch('/api/generic', {
    initialCache: false,
    method: 'post',
    body: { collection: 'configurations', projection: {}, filter: { property_type: { $eq: 'inquiry-form-reasons' } }, limit: 10000 },
    onResponse({ request, response, options }) {
      //console.log(response._data.documents);
    },
  })
);*/

//-----------------------------------------------------------------------------------------------------
//4. Event Subscriptions
//-----------------------------------------------------------------------------------------------------
useNuxtApp().$bus.$on('evtSideBarGuide', async (args) => {
  isSideBarGuideOpen.value = !isSideBarGuideOpen.value;
  const data = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/blogs?title.eq=${props.content_identifier}`);
  console.log(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/blogs?title.eq=${props.content_identifier}`)
  if (data && Array.isArray(data.data)) {
    page.value = await parseMarkdown(data.data[0].content);
  } else {
    page.value = { data: { title: 'Sorry, No page exists', body: '' }, toc: { links: [] } };
  }

  //data.value = {};
  progress.value = false;
});

function cancel() {
  try {
    isSideBarGuideOpen.value = false;
  } catch (e) {
    console.log(e);
  }

  return false;
}

//-----------------------------------------------------------------------------------------------------
//5. Page level events handlers
//-----------------------------------------------------------------------------------------------------

function handleSelectedInInquiry_Reason(e) {
  data.value.inquiry_reason = e.target.value;
}

//-----------------------------------------------------------------------------------------------------
//6. Page level action methods
//-----------------------------------------------------------------------------------------------------
async function saveSideBarGuide(args) {
  try {
    // Validate form data
    if (!modelValidate(data.value)) {
      //Show user friendly error message
      useNuxtApp().$toast.show({
        type: 'error',
        message: 'Please fill all required inputs',
        timeout: 10,
      });
      return;
    }
    progress.value = true;

    //attach tenant into the payload (Required for SaaS)
    data.value.tenant = 'www.xyz-company.com'; //todo: dynamic get from settings

    //Send to server
    await useFetch('/api/inquiries/create', {
      method: 'post',
      body: data.value,
      initialCache: false,
      onResponse({ request, response, options }) {
        console.log(JSON.stringify(response));
        if (response._data) {
          if (response._data.code == 200) {
            useNuxtApp().$toast.show({ type: 'success', message: `Thank you for your inquiry, one of our associate will contact you soon`, timeout: 6 });
            //data.value = {};
            isSideBarGuideOpen.value = !isSideBarGuideOpen.value;
          } else {
            useNuxtApp().$toast.show({ type: 'error', message: `Failure during inquiry save`, timeout: 6 });
          }
          progress.value = false;
        }
      },
    });
  } catch (error) {
    //Show user friendly error message
    useNuxtApp().$toast.show({
      type: 'error',
      message: `Oops!... Something went wrong....<br/>[${error.message}]`,
      timeout: 6,
    });
    //reset progress
    progress.value = false;
  } finally {
  }
}

function modelValidate(args) {
  let name = data.value.name != undefined && data.value.name != '';
  let phone_number = data.value.phone_number != undefined && data.value.phone_number != '';
  let email = data.value.email != undefined && data.value.email != '';
  let description = data.value.description != undefined && data.value.description != '';
  console.log('name && phone_number && email && description', name && phone_number && email && description);
  return name && phone_number && email && description;
}

async function getContent() {
  data.value = await useAsyncData(`content`, () => {
    return queryContent().where({ _path: '/docs/getting-started' }).findOne();
  });
}
</script>
