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
</script>

<template>
  <div class="max-w-4xl mx-auto py-10 px-6">
    <h1 class="text-3xl font-bold text-gray-800 mb-8 border-b pb-4">✏️ แก้ไขบทความ</h1>

    <div v-if="loading" class="text-gray-500 text-center py-6">⏳ กำลังโหลด...</div>
    <div v-if="error" class="text-red-600 bg-red-100 p-4 rounded mb-6">{{ error }}</div>

    <form @submit.prevent="updateArticle" v-if="!loading" class="space-y-6">
      <!-- Title -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">หัวข้อบทความ</label>
        <input v-model="title" type="text"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" required />
      </div>

      <!-- Content -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">เนื้อหา</label>
        <textarea v-model="content" rows="10"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" required></textarea>
      </div>

      <!-- Category -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">หมวดหมู่</label>
        <select v-model="selectedCategory"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" required>
          <option value="" disabled>-- เลือกหมวดหมู่ --</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>

      <!-- Tags -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Tags</label>

        <div class="flex flex-wrap gap-2 mb-2">
          <span v-for="tag in selectedTags" :key="tag.id || tag.name"
            class="bg-blue-100 text-blue-700 px-3 py-1 rounded-full flex items-center text-sm">
            {{ tag.name }}
            <button type="button" class="ml-2 text-blue-500 hover:text-red-500" @click="removeTag(tag)">
              ✕
            </button>
          </span>
        </div>

        <div class="relative">
          <input v-model="tagInput" @keydown.enter.prevent="handleTagInput" @keydown.tab.prevent="handleTagInput"
            @focus="showSuggestions = true" @blur="handleBlur" type="text" placeholder="พิมพ์แท็ก เช่น Go, Docker, Vue"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400" />

          <!-- Dropdown Suggestion -->
          <ul v-if="showSuggestions && filteredTagSuggestions.length > 0"
            class="absolute z-10 bg-white border border-gray-300 w-full mt-1 rounded-md shadow-lg max-h-48 overflow-auto">
            <li v-for="tag in filteredTagSuggestions" :key="tag.id" @mousedown.prevent="selectTag(tag)"
              class="px-4 py-2 cursor-pointer hover:bg-blue-100">
              {{ tag.name }}
            </li>
          </ul>
        </div>
      </div>

      <!-- Save Button -->
      <div class="text-right">
        <button type="submit"
          class="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-6 rounded-lg transition">
          💾 บันทึกการแก้ไข
        </button>
      </div>
    </form>
  </div>
</template>
