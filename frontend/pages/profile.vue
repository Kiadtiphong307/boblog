<template>
  <div class="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-2xl mx-auto bg-white shadow-xl rounded-2xl p-8">
      <h2 class="text-3xl font-bold text-center text-blue-700 mb-8">แก้ไขโปรไฟล์</h2>

      <div v-if="success" class="bg-green-100 text-green-800 px-4 py-2 rounded mb-4 text-sm font-medium">
        {{ success }}
      </div>
      <div v-if="error" class="bg-red-100 text-red-800 px-4 py-2 rounded mb-4 text-sm font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="updateProfile" class="space-y-6">
        <!-- Avatar -->
        <div class="flex flex-col items-center">
          <div class="w-24 h-24 rounded-full border overflow-hidden shadow">
            <img :src="previewImage" v-if="previewImage" class="w-full h-full object-cover" />
            <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center text-gray-500 text-sm">
              {{ PlaceholdersImage.noImage }}
            </div>
          </div>
          <label class="mt-4 text-sm font-medium text-gray-700">{{ ProfileText.uploadLabel }}</label>
          <input type="file" accept="image/*" @change="handleFileChange" class="mt-1 text-sm text-gray-600" />
          <p class="text-xs text-gray-400 mt-1">{{ PlaceholdersImage.imageNote }}</p>
        </div>

        <!-- First Name -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">{{ ProfileText.firstName }}</label>
          <input v-model="form.first_name" type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>

        <!-- Last Name -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">{{ ProfileText.lastName }}</label>
          <input v-model="form.last_name" type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>

        <!-- Nickname -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">{{ ProfileText.nickname }}</label>
          <input v-model="form.nickname" type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>

        <!-- Bio -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">{{ ProfileText.bio }}</label>
          <textarea v-model="form.bio" rows="4"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"></textarea>
        </div>

        <!-- Submit -->
        <div class="pt-4 text-right">
          <button type="submit"
            class="inline-block bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-2 rounded-lg shadow-md transition">
            💾 บันทึกข้อมูล
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProfile } from '@/composables/useProfile'
import { ProfileText } from '@/constants/Profile/Profile'
import { PlaceholdersImage } from '@/constants/Placeholders'

const {
  user,
  form,
  selectedFile,
  previewImage,
  success,
  error,
  loading,
  handleFileChange,
  updateProfile,
} = useProfile()
</script>