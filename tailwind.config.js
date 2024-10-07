/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.templ"],
  theme: {
    extend: {},
  },
  // plugins: [require('daisyui')],
  plugins: [],
  // https://daisyui.com/docs/themes/
  daisyui: {
    themes: ["retro", "dim"],
    darkTheme: "dim"
  },
}

