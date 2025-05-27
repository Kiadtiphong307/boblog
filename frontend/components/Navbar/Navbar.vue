<template>
  <nav class="flex items-center justify-between px-6 md:px-24 py-6 bg-white shadow-md border-b border-blue-200">
    <div class="flex items-center space-x-6">
      <NuxtLink to="/" class="text-blue-500 font-bold text-3xl">BOBLOG</NuxtLink>
      <NuxtLink to="/articles/create" class="flex items-center text-lg text-black">✏️ เขียนบทความ</NuxtLink>
      <NuxtLink to="/articles" class="flex items-center text-lg text-black">🔍 ค้นหา</NuxtLink>
    </div>

    <div class="relative">
      <NuxtLink v-if="!nickname" to="/login">
        <button class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-6 rounded transition">
          เข้าสู่ระบบ
        </button>
      </NuxtLink>

      <template v-else>
        <div @click="showDropdown = !showDropdown"
          class="flex items-center gap-2 cursor-pointer px-3 py-2 rounded hover:bg-gray-100 transition">
          <div v-if="imageUrl" class="w-8 h-8 rounded-full overflow-hidden border">
            <img :src="imageUrl" alt="avatar" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-8 h-8 bg-blue-200 rounded-full flex items-center justify-center text-white font-bold">
            {{ nickname.charAt(0).toUpperCase() }}
          </div>
          <span class="text-gray-800 font-medium hidden md:inline">
            สวัสดี, {{ nickname }}
          </span>
          <svg class="h-4 w-4 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
          </svg>
        </div>

        <transition name="fade">
          <div v-show="showDropdown"
            class="absolute right-0 mt-2 w-48 bg-white border rounded shadow-lg z-50 overflow-hidden text-sm">
            <NuxtLink to="/profile" class="block px-4 py-3 text-gray-700 hover:bg-gray-100 transition">👤 โปรไฟล์</NuxtLink>
            <NuxtLink to="/articles/my-articles" class="block px-4 py-3 text-gray-700 hover:bg-gray-100 transition">📚 จัดการบทความ</NuxtLink>
            <button @click="logout" class="block w-full text-left px-4 py-3 text-red-600 hover:bg-red-100 transition">
              🚪 ออกจากระบบ
            </button>
          </div>
        </transition>
      </template>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useAuthInfo } from '@/composables/useAuthInfo'
const { nickname, imageUrl, showDropdown, logout } = useAuthInfo()
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>