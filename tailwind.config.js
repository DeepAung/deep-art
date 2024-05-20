const defaultTheme = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js,templ,go}",
    "node_modules/preline/dist/*.js",
  ],
  theme: {
    fontFamily: {
      sans: ["Inter"],
      display: ["Inter"],
      body: ["Oswald"],
    },
  },
  plugins: [require("preline/plugin")],
};
