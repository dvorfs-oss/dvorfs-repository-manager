import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{js,ts,jsx,tsx,mdx}"],
  theme: {
    extend: {
      colors: {
        ink: {
          50: "#f4f7fb",
          100: "#e6edf7",
          200: "#c8d4ea",
          300: "#9eb0d7",
          400: "#748dc0",
          500: "#526ca3",
          600: "#425787",
          700: "#344465",
          800: "#24314a",
          900: "#101a2b",
          950: "#070b12",
        },
        accent: {
          50: "#eefbf6",
          100: "#d7f6ea",
          200: "#afe9d4",
          300: "#76d5b3",
          400: "#3fbf93",
          500: "#159e74",
          600: "#0d7c5b",
          700: "#0a6249",
          800: "#084c39",
          900: "#053226",
        },
      },
      boxShadow: {
        glow: "0 24px 80px rgba(6, 12, 22, 0.45)",
      },
      backgroundImage: {
        "mesh-radial":
          "radial-gradient(circle at top left, rgba(63, 191, 147, 0.25), transparent 34%), radial-gradient(circle at top right, rgba(82, 108, 163, 0.18), transparent 26%), linear-gradient(180deg, rgba(7, 11, 18, 1), rgba(10, 15, 25, 1))",
      },
      fontFamily: {
        sans: [
          '"Avenir Next"',
          '"Segoe UI"',
          '"Helvetica Neue"',
          "Arial",
          "sans-serif",
        ],
      },
    },
  },
  plugins: [],
};

export default config;
