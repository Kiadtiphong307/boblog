<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-xl p-8 space-y-6">
      <h1 class="text-2xl font-bold text-center text-gray-800">
        📝 สมัครสมาชิก
      </h1>

      <form @submit.prevent="register" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.username" type="text" :placeholder="formPlaceholders.username"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.username" class="text-sm text-red-500 mt-1">{{ error.username }}</p>
          </div>
          <div>
            <input v-model="form.nickname" type="text" :placeholder="formPlaceholders.nickname"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.nickname" class="text-sm text-red-500 mt-1">{{ error.nickname }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.first_name" type="text" :placeholder="formPlaceholders.firstName"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.first_name" class="text-sm text-red-500 mt-1">{{ error.first_name }}</p>
          </div>
          <div>
            <input v-model="form.last_name" type="text" :placeholder="formPlaceholders.lastName"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.last_name" class="text-sm text-red-500 mt-1">{{ error.last_name }}</p>
          </div>
        </div>

        <div>
          <input v-model="form.email" type="email" :placeholder="formPlaceholders.email"
            class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
          <p v-if="error.email" class="text-sm text-red-500 mt-1">{{ error.email }}</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.password" type="password" :placeholder="formPlaceholders.password"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.password" class="text-sm text-red-500 mt-1">{{ error.password }}</p>
          </div>
          <div>
            <input v-model="form.confirm_password" type="password" :placeholder="formPlaceholders.confirmPassword"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.confirm_password" class="text-sm text-red-500 mt-1">{{ error.confirm_password }}</p>
          </div>
        </div>

        <button type="submit"
          class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition">
          ➕ สมัครสมาชิก
        </button>
      </form>

      <p v-if="success" class="text-green-600 text-center font-medium">{{ success }}</p>
      <p v-if="error.general" class="text-red-600 text-center font-medium">{{ error.general }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

import { useRegister } from '@/composables/useRegister'
import { formPlaceholders } from '~/constants/Placeholders'

const { form, error, success, register } = useRegister()
</script>