
const plugin = require('tailwindcss/plugin')
const colors = require('tailwindcss/colors')

module.exports = {
  mode: 'jit',
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    './src/**/*.{js,jsx,ts,tsx,vue}',
    "node_modules/tailvue/dist/tailvue.es.js"
  ],
  safelist: [
    //orange
    'bg-[#fefce8]',
    'bg-primary-100',
    'bg-primary-200',
    'bg-primary-200',
    'bg-[#f97316]',
    'bg-[#ca8a04]',
    //teal
    'bg-[#f0fdfa]', 'bg-[#ccfbf1]', 'bg-[#99f6e4]', 'bg-[#14b8a6]', 'bg-[#0d9488]',
    'bg-[#eff6ff]', 'bg-[#ebf8ff]', 'bg-[#bee3f8]', 'bg-[#90cdf4]', 'bg-[#63b3ed]', 'bg-[#4299e1]', 'bg-[#3182ce]', 'bg-[#2b6cb0]', 'bg-[#2c5282]', 'bg-[#2a4365]',
    'bg-[#fff7ed]', 'bg-[#fffaf0]', 'bg-[#feebc8]', 'bg-[#fbd38d]', 'bg-[#f6ad55]', 'bg-[#ed8936]', 'bg-[#dd6b20]', 'bg-[#c05621]', 'bg-[#9c4221]', 'bg-[#7b341e]',
    'bg-[#f0fdfa]', 'bg-[#e6fffa]', 'bg-[#b2f5ea]', 'bg-[#81e6d9]', 'bg-[#4fd1c5]', 'bg-[#38b2ac]', 'bg-[#319795]', 'bg-[#2c7a7b]', 'bg-[#285e61]', 'bg-[#234e52]',
    'bg-[#f0fdf4]', 'bg-[#c6f6d5]', 'bg-[#9ae6b4]', 'bg-[#68d391]', 'bg-[#48bb78]', 'bg-[#38a169]', 'bg-[#2f855a]', 'bg-[#276749]', 'bg-[#22543d]',
    'bg-[#eef2ff]', 'bg-[#ebf4ff]', 'bg-[#c3dafe]', 'bg-[#a3bffa]', 'bg-[#7f9cf5]', 'bg-[#667eea]', 'bg-[#5a67d8]', 'bg-[#4c51bf]', 'bg-[#434190]', 'bg-[#3c366b]',
    'bg-[#fefce8]', 'bg-[#fffff0]', 'bg-[#fefcbf]', 'bg-[#faf089]', 'bg-[#f6e05e]', 'bg-[#ecc94b]', 'bg-[#d69e2e]', 'bg-[#b7791f]', 'bg-[#975a16]', 'bg-[#744210]',
    'bg-[#fef2f2]', 'bg-[#fff5f5]', 'bg-[#fed7d7]', 'bg-[#feb2b2]', 'bg-[#fc8181]', 'bg-[#f56565]', 'bg-[#e53e3e]', 'bg-[#c53030]', 'bg-[#9b2c2c]', 'bg-[#742a2a]',
    'bg-[#fdf2f8]', 'bg-[#fff5f7]', 'bg-[#fed7e2]', 'bg-[#fbb6ce]', 'bg-[#f687b3]', 'bg-[#ed64a6]', 'bg-[#d53f8c]', 'bg-[#b83280]', 'bg-[#97266d]', 'bg-[#702459]',
    'bg-[#faf5ff]', 'bg-[#faf5ff]', 'bg-[#e9d8fd]', 'bg-[#d6bcfa]', 'bg-[#b794f4]', 'bg-[#9f7aea]', 'bg-[#805ad5]', 'bg-[#6b46c1]', 'bg-[#553c9a]', 'bg-[#44337a]',
    'bg-[#f9fafb]', 'bg-[#f7fafc]', 'bg-[#edf2f7]', 'bg-[#e2e8f0]', 'bg-[#cbd5e0]', 'bg-[#a0aec0]', 'bg-[#718096]', 'bg-[#4a5568]', 'bg-[#2d3748]', 'bg-[#1a202c]',
    'hover:bg-[#5a67d8]', 'hover:bg-[#f9fafb]', 'hover:bg-[#fef2f2]', 'hover:bg-[#fff7ed]', 'hover:bg-[#fefce8]', 'hover:bg-[#f0fdf4]', 'hover:bg-[#f0fdfa]', 'hover:bg-[#eff6ff]', 'hover:bg-[#eef2ff]', 'hover:bg-[#faf5ff]', 'hover:bg-[#fdf2f8]',
    'text-[#4a5568]', 'text-[#c53030]', 'text-[#c05621]', 'text-[#b7791f]', 'text-[#2f855a]', 'text-[#2c7a7b]', 'text-[#2b6cb0]',
    'text-[#4c51bf]', 'text-[#6b46c1]', 'text-[#b83280]',
    'hover:text-[#4a5568]', 'hover:text-[#c53030]', 'hover:text-[#c05621]', 'hover:text-[#b7791f]', 'hover:text-[#2f855a]', 'hover:text-[#2c7a7b]', 'hover:text-[#2b6cb0]', 'hover:text-[#4c51bf]', 'hover:text-[#6b46c1]', 'hover:text-[#b83280]',
    'hover:bg-[#4a5568]', 'hover:bg-[#c53030]', 'hover:bg-[#c05621]', 'hover:bg-[#b7791f]', 'hover:bg-[#2f855a]', 'hover:bg-[#2c7a7b]', 'hover:bg-[#2b6cb0]', 'hover:bg-[#4c51bf]', 'hover:bg-[#6b46c1]', 'hover:bg-[#b83280]',
    'via-[#5a67d8]', 'via-[#f9fafb]', 'via-[#fef2f2]', 'via-[#fff7ed]', 'via-[#fefce8]', 'via-[#f0fdf4]', 'via-[#f0fdfa]', 'via-[#eff6ff]', 'via-[#eef2ff]', 'via-[#faf5ff]', 'via-[#fdf2f8]',
    'hover:text-[#5a67d8]', 'hover:text-[#f9fafb]', 'hover:text-[#fef2f2]', 'hover:text-[#fff7ed]', 'hover:text-[#fefce8]', 'hover:text-[#f0fdf4]', 'hover:text-[#f0fdfa]', 'hover:text-[#eff6ff]', 'hover:text-[#eef2ff]', 'hover:text-[#faf5ff]', 'hover:text-[#fdf2f8]',

  ],
  corePlugins: {
    // due to https://github.com/tailwindlabs/tailwindcss/issues/6602 - buttons disappear
    preflight: false,
  },
  darkMode: "class",
  theme: {
    extend: {
      aspectRatio: {
        '4/3': '4 / 3',
      },
      colors: {
        transparent: 'transparent',
        current: 'currentColor',
        'white': '#ffffff',
        //'purple': '#3f3cbb',
        'midnight': '#121063',
        'metal': '#565584',
        'tahiti': '#3ab7bf',
        'silver': '#ecebff',
        'bubble-gum': '#ff77e9',
        'bermuda': '#78dcca',
        fuchsia: colors.fuchsia,
        pink: colors.pink,
        purple: colors.purple,
        indigo: colors.indigo,
        emerald: colors.emerald,
        gray: colors.gray,
        red: colors.red,
        yellow: colors.yellow,
        green: colors.green,
        lime: colors.lime,
        cyan: colors.cyan,
        blue: colors.blue,
        purple: colors.purple,
        pink: colors.pink,
        rose: colors.rose,
        teal: colors.teal,
        orange: colors.orange,
        magenta: colors.magenta,
        primary: colors.teal,
        secondary: colors.gray,
        footer: colors.slate,
        sy: colors.yellow,
      },
      animation: {
        marquee: 'marquee 25s linear infinite',
        marquee2: 'marquee2 25s linear infinite',
      },
      keyframes: {
        marquee: {
          '0%': { transform: 'translateX(0%)' },
          '100%': { transform: 'translateX(-100%)' },
        },
        marquee2: {
          '0%': { transform: 'translateX(100%)' },
          '100%': { transform: 'translateX(0%)' },
        },
      },
      fontFamily: {
        sans: ['Poppins', 'sans-serif'],
        spotlight: ['Hey August'],
      },
    },
  },
  variants: {
    textColor: ['responsive', 'hover', 'focus', 'group-hover'],
    backgroundColor: ['responsive', 'hover', 'focus', 'even', 'odd'],
  },
  corePlugins: {
    aspectRatio: false,
  },
  plugins: [plugin, require('@tailwindcss/forms'),require('@tailwindcss/typography')],
};