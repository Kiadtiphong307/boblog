<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded shadow-md w-full max-w-md">
      <h1 class="text-2xl font-bold mb-6 text-center">🔐 เข้าสู่ระบบ</h1>

      <!-- Form -->
      <form @submit.prevent="login" class="space-y-4">
        <input v-model="emailOrUsername" type="text" :placeholder="PlaceholdersLogin.emailOrUsername" required
          class="input input-bordered w-full" />
        <input v-model="password" type="password" :placeholder="PlaceholdersLogin.password" required
          class="input input-bordered w-full" />
        <button type="submit" class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition">
          เข้าสู่ระบบ
        </button>
      </form>

      <!-- Error -->
      <p v-if="error" class="text-red-600 mt-3 text-center">{{ error }}</p>

      <!-- Register -->
      <p class="mt-6 text-center text-sm text-gray-600">
        ยังไม่มีบัญชี?
        <NuxtLink to="/register" class="text-blue-600 hover:underline">
          สมัครสมาชิก
        </NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

import { useLogin } from '@/composables/useLogin'
import { PlaceholdersLogin } from '~/constants/Placeholders'

const { emailOrUsername, password, error, login } = useLogin()
</script>