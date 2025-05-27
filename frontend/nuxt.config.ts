import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  app:{
    head:{
      title: 'BoBlog', // ชื่อแท็บเบราว์เซอร์
      meta: [
        { name: 'description', content: 'ระบบจัดการบทความออนไลน์' }
      ]
    }
  },
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  
});