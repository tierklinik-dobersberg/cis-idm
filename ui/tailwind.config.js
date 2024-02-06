/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,ts}",
    "./node_modules/tw-elements/dist/js/**/*.js",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Source Code Pro", "sans-serif"],
      },
    },
  },
  plugins: [require("tw-elements/dist/plugin.cjs")],
};
