<template>
    <div>
      <div class="bg-blue-500 min-h-[400px] flex flex-col items-center justify-center w-full py-16">
        <div class="text-white text-4xl md:text-5xl font-bold mb-4 pb-2">
          {{ HeroSection.Section }}
        </div>
        <NuxtLink to="/articles/create" class="bg-white text-blue-500 font-semibold px-8 py-3 rounded hover:bg-blue-100 shadow-lg 
          transition hover:scale-110 hover:ease-in-out hover:shadow-xl hover:duration-700">
          {{ Menu.title_create_new_articel }}
        </NuxtLink>
      </div>
    </div>
  </template>
  
<script setup>
import { HeroSection } from '~/constants/HeroSection';
import { Menu } from '~/constants/Menu';
</script>

