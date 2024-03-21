import { defineNuxtPlugin } from '#app'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime.js'
import localeData from 'dayjs/plugin/localeData.js'
import updateLocale from 'dayjs/plugin/updateLocale.js'
import utc from 'dayjs/plugin/utc'
export default defineNuxtPlugin((nuxtApp) => {
    dayjs.extend(relativeTime)
    dayjs.extend(localeData)
    dayjs.extend(updateLocale)
    dayjs.extend(utc)
    nuxtApp.provide('dayjs', dayjs)
})

