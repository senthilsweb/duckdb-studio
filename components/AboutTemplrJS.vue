<template>
	<TransitionRoot appear :show="isOpen" as="template">
		<Dialog as="div" @close="closeModal" class="relative z-50">
			<TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0" enter-to="opacity-100"
				leave="duration-200 ease-in" leave-from="opacity-100" leave-to="opacity-0">
				<div class="fixed inset-0 bg-black bg-opacity-25" />
			</TransitionChild>

			<div class="fixed inset-0 overflow-y-auto">
				<div class="flex min-h-full items-center justify-center p-4 text-center">
					<TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0 scale-95"
						enter-to="opacity-100 scale-100" leave="duration-200 ease-in" leave-from="opacity-100 scale-100"
						leave-to="opacity-0 scale-95">
						<DialogPanel
							class="w-full max-w-md transform overflow-hidden rounded-2xl bg-primary-50 p-6 text-left align-middle shadow-xl transition-all">
							<DialogTitle as="h3" class="text-md font-medium leading-6 text-gray-900">DuckDB Studio! A simple, yet powerful web utility designed to explore and interact with DuckDB databases efficiently. </DialogTitle>
							<div class="mt-2">
								<p class="text-sm text-gray-500">version 1.1.20240411.1</p>
							</div>
						</DialogPanel>
					</TransitionChild>
				</div>
			</div>
		</Dialog>
	</TransitionRoot>
</template>

<script setup>
import { ref } from "vue";
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from "@headlessui/vue";

const isOpen = ref(false);

useNuxtApp().$bus.$on("evtAboutTemplrJS", (data) => {
	//alert("Hello");
	isOpen.value = true;
});

function closeModal() {
	isOpen.value = false;
}
function openModal() {
	isOpen.value = true;
}
</script>
