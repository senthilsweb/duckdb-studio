<template>
  <div>
    <!--<label v-if="show_label" :for="name" class="block text-sm font-medium text-gray-700">
      {{ label }}
    </label>-->
    <div class="mt-1">
      <select @change="onChangeDropdown" :id="name" :autocomplete="name" :name="name" :readonly="readonly" :disabled="readonly" class="h-8 mt-2 block w-full border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6">
        <option value="">Select Schema</option>
        <option v-for="(item, index) in items" :key="index" :value="item.code" :selected="item.code === selectedValue">{{ label }} : {{ item.name }}</option>
      </select>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  type: String,
  label: String,
  show_label: Boolean,
  readonly: Boolean,
  name: String,
  selectedValue: String,
  selectedText: String,
  api: String,
  data: Array,
});

const emit = defineEmits(['selectedItem']);
let items = ref([]);

onMounted(async () => {
  if (props.api) {
    items.value = await useNuxtApp().$fetch(props.api);
  } else if (props.data) {
    items.value = props.data;
  }
});

function onChangeDropdown(e) {
  const selectedCode = e.target.value;
  const option = items.value.find((option) => selectedCode === option.code);
  if (option) {
    emit('selectedItem', option.code);
  } else {
    emit('selectedItem', '');
  }
}
</script>
