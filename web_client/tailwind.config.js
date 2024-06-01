module.exports = {
  content: ["../httpserver/templates/*.gotmpl", "./src/**/*.tsx"],
  darkMode: ["class", '[data-theme="dark"]'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: "#eef2ff",
          100: "#e0e7ff",
          200: "#c7d2fe",
          300: "#a5b4fc",
          400: "#818cf8",
          500: "#6366f1",
          600: "#4f46e5",
          700: "#4338ca",
          800: "#3730a3",
          900: "#312e81",
          950: "#1e1b4b",
        },
      },
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
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/aspect-ratio"),
  ],
};
