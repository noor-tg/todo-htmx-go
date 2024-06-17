const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./**/*.go"],
  theme: {
    extends: {
      colors,
    },
  },
};
