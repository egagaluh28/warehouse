// tailwind.config.js
const { nextui } = require("@nextui-org/react");

module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        "primary-green": "#005F02",
        "secondary-green": "#427A43",
      },
    },
  },
  darkMode: "class",
  plugins: [
    nextui({
      layout: {
        radius: {
          small: "4px",
          medium: "8px",
          large: "12px",
        },
      },
      themes: {
        light: {
          colors: {
            // Calm Nature Green Theme
            primary: {
              DEFAULT: "#15803d", // Green 700 (Deep, professional green)
              foreground: "#ffffff",
              50: "#f0fdf4",
              100: "#dcfce7",
              200: "#bbf7d0",
              300: "#86efac",
              400: "#4ade80",
              500: "#22c55e",
              600: "#16a34a",
              700: "#15803d",
              800: "#166534",
              900: "#14532d",
            },
            secondary: {
              DEFAULT: "#0f766e", // Teal 700
              foreground: "#ffffff",
            },
            background: "#FFFFFF",
          },
        },
        dark: {
          colors: {
            primary: {
              DEFAULT: "#22c55e", // Green 500 (Lighter for dark mode)
              foreground: "#000000",
            },
          },
        },
      },
    }),
  ],
};
