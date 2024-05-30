module.exports = {
  content: ["../httpserver/templates/*.gotmpl", "./src/**/*.tsx"],
  darkMode: ["class", '[data-theme="dark"]'],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
      fontFamily: {
        poppins: [
          "Poppins",
          "Inter",
          "-apple-system",
          "Lucida Grande",
          "Tahoma",
          "Sans-Serif",
        ],
        inter: [
          "Inter",
          "-apple-system",
          "Lucida Grande",
          "Tahoma",
          "Sans-Serif",
        ],
        body: [
          "Merriweather",
          "-apple-system",
          "Lucida Grande",
          "Tahoma",
          "Sans-Serif",
        ],
        ui: [
          "Roboto",
          "-apple-system",
          "Lucida Grande",
          "Tahoma",
          "Sans-Serif",
        ],
      },
      colors: {},
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/aspect-ratio"),
  ],
};
