import { defineNuxtPlugin } from '#app'
import underscorestring from 'underscore.string';

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.provide('s', underscorestring)
})
