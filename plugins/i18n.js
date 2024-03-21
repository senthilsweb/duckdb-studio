import { createI18n } from "vue-i18n";
import en from '../locales/en.json'
import ja from '../locales/ja.json'

export default defineNuxtPlugin(({ vueApp }) => {
	const i18n = createI18n({
		legacy: false,
		globalInjection: true,
		locale: "en",
		messages: {
			en,
			ja
		},
	});

	vueApp.use(i18n);
});
