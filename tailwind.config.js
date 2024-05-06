/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["**/*.html", "**/*.go"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

