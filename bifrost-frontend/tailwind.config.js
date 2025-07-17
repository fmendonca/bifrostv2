/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        bifrostBlue: '#1E40AF',
        bifrostGreen: '#10B981',
        bifrostRed: '#EF4444',
        bifrostYellow: '#F59E0B'
      }
    },
  },
  plugins: [],
}
