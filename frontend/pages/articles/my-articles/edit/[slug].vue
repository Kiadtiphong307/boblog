<script setup lang="ts">
import { useEditArticles } from '~/composables/articles/useEditArticles'

const {
  title,
  content,
  selectedCategory,
  selectedTags,
  tagInput,
  showSuggestions,
  categories,
  tags,
  loading,
  error,
  filteredTagSuggestions,
  handleBlur,
  handleTagInput,
  selectTag,
  removeTag,
  updateArticle,
} = useEditArticles()

const handleSubmit = async () => {
  if (!title.value || !content.value) {
    alert('กรุณากรอกชื่อและเนื้อหาบทความ')
    return
  }
  await updateArticle()
}

useHead({
  title: () => `แก้ไข: ${title.value || 'Loading...'}`,
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 py-10 px-6">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-10">
        <h1 class="text-4xl font-bold text-gray-900">📝 แก้ไขบทความ</h1>
        <p class="mt-2 text-gray-600">กรุณาปรับปรุงเนื้อหาและรายละเอียดของบทความของคุณ</p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="bg-white rounded-2xl shadow p-12 text-center">
        <div class="animate-spin h-10 w-10 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
        <p class="mt-4 text-gray-600">กำลังโหลดข้อมูลบทความ...</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="bg-white rounded-2xl shadow p-8 text-red-700">
        <p class="text-center text-lg font-semibold">❌ {{ error }}</p>
        <div class="mt-6 text-center">
          <NuxtLink to="/articles/my-articles" class="text-blue-600 hover:underline">กลับไปหน้าบทความของฉัน</NuxtLink>
        </div>
      </div>

      <!-- Form -->
      <div v-else class="bg-white rounded-2xl shadow-lg p-8">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Title -->
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">ชื่อบทความ <span class="text-red-500">*</span></label>
            <input
              v-model="title"
              type="text"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
              placeholder="ใส่ชื่อบทความที่น่าสนใจ"
              required
            />
          </div>

          <!-- Content -->
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">เนื้อหาบทความ <span class="text-red-500">*</span></label>
            <textarea
              v-model="content"
              rows="10"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition resize-none"
              placeholder="เขียนเนื้อหาบทความของคุณที่นี่..."
              required
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">รองรับ HTML tags พื้นฐาน</p>
          </div>

          <!-- Category -->
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">หมวดหมู่</label>
            <select
              v-model="selectedCategory"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
            >
              <option :value="null">-- เลือกหมวดหมู่ --</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <!-- Tags -->
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">แท็ก</label>

            <!-- Selected Tags -->
            <div v-if="selectedTags.length" class="flex flex-wrap gap-2 mb-2">
              <span
                v-for="tag in selectedTags"
                :key="tag.name"
                class="bg-blue-100 text-blue-800 px-3 py-1 text-sm rounded-full inline-flex items-center"
              >
                {{ tag.name }}
                <button @click="removeTag(tag)" class="ml-2 text-blue-500 hover:text-red-500 focus:outline-none">
                  &times;
                </button>
              </span>
            </div>

            <!-- Tag Input -->
            <div class="relative">
              <input
                v-model="tagInput"
                @input="showSuggestions = true"
                @blur="handleBlur"
                @keypress.enter.prevent="handleTagInput"
                class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
                placeholder="พิมพ์แท็กแล้วกด Enter"
              />
              <div
                v-if="showSuggestions && filteredTagSuggestions.length"
                class="absolute z-10 mt-1 w-full bg-white border border-gray-300 rounded-xl shadow max-h-60 overflow-auto"
              >
                <button
                  v-for="suggestion in filteredTagSuggestions"
                  :key="suggestion.id"
                  type="button"
                  @click="selectTag(suggestion)"
                  class="w-full text-left px-4 py-2 hover:bg-gray-100 transition"
                >
                  {{ suggestion.name }}
                </button>
              </div>
            </div>
          </div>

          <!-- Buttons -->
          <div class="flex flex-col sm:flex-row gap-4 pt-6 border-t border-gray-200">
            <button
              type="submit"
              :disabled="!title || !content"
              class="w-full sm:w-auto px-6 py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition disabled:opacity-50"
            >
              💾 บันทึกการแก้ไข
            </button>
            <NuxtLink
              to="/articles/my-articles"
              class="w-full sm:w-auto px-6 py-3 bg-gray-100 text-gray-700 rounded-xl font-medium text-center hover:bg-gray-200 transition"
            >
              ยกเลิก
            </NuxtLink>
          </div>
        </form>
      </div>

      <!-- Back -->
      <div class="mt-8 text-center">
        <NuxtLink to="/articles/my-articles" class="text-gray-600 hover:text-blue-600 transition">
          ← กลับไปหน้าบทความของฉัน
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
