<template>
  <div>
    <h5 class="text-gray-900 uppercase tracking-wide font-semibold mb-3 text-sm lg:text-xs">On this page</h5>
    <ol role="list" class="mt-4 space-y-3 text-sm">
      <li v-for="link of data" :key="link.id">
        <h3>
          <a :href="`#${link.id}`" class="block text-sm text-zinc-600 hover:text-primary-600 font-bold mb-3">{{ link.text }}</a>
        </h3>

        <ol v-if="link.children" role="list" class="mt-2 space-y-3 pl-5 text-slate-500 dark:text-slate-400">
          <li v-for="link of link.children" :key="link.id">
            <a :href="`#${link.id}`" class="text-slate-600 hover:text-primary-600">{{ link.text }}</a>
          </li>
        </ol>
      </li>
    </ol>
  </div>
</template>

<script>
export default {
  name: 'toc',
  data() {
    return {};
  },
  props: {
    data: Array,
  },
  methods: {},
  created() {},
};
</script>
