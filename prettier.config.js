module.exports = {
  plugins: ["prettier-plugin-tailwindcss-extra", "prettier-plugin-tailwindcss"],
  overrides: [
    {
      files: "*.templ",
      options: { parser: "tailwindcss-extra" },
    },
  ],
};
