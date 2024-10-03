import colors from "tailwindcss/colors"

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extends: {
      colors,
    },
  },
};
