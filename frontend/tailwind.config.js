/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        BrandBlue: '#189aca',
        BrandOrange: '#ffa200',
      }
    },
  },
  plugins: [],
}

