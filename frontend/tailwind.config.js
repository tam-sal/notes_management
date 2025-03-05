/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx}',
  ],
  theme: {
    extend: {
    },
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      {
        nord: {
          primary: "#f6d860",
          "base-100": "#ffffff",
          "base-content": "#1f2937",
          info: "#F1F5F9",
          neutral: "#2a2a2a",
        },
      },
      {
        luxury: {
          primary: "#ffffff",
          "base-100": "#1e1e28",
          "base-content": "#f3f4f6",
          info: "#f6d860",
          neutral: "#2a2a2a",
        },
      },
    ],
    darkTheme: "luxury",
  },
  styled: true,
  themeRoot: "*",
}


