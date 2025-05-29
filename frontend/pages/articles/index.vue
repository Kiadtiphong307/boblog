<template>
  <div class="max-w-6xl mx-auto py-10 px-4">
    <h1 class="text-2xl font-bold mb-6">📚 บทความทั้งหมด</h1>
    
    <Filter
      :categories="categories"
      :search-term="searchTerm"
      :selected-category="selectedCategory"
      :update-search-term="updateSearchTerm"
      :update-selected-category="updateSelectedCategory"
    />
    
    <div v-if="loading" class="text-gray-500">⏳ กำลังโหลดบทความ...</div>
    <div v-else-if="articles.length === 0" class="text-red-500">❌ ไม่พบบทความ</div>
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card
          v-for="article in articles"
          :key="article.id"
          :article="article"
          :format-date="formatDate"
        />
      </div>

      <!-- Pagination Controls -->
      <div class="mt-8 flex justify-center items-center space-x-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-3 py-1 rounded border"
        >
          ก่อนหน้า
        </button>

        <button
          v-for="page in totalPages"
          :key="page"
          @click="goToPage(page)"
          :class="[
            'px-3 py-1 rounded border',
            page === currentPage ? 'bg-blue-600 text-white' : 'bg-white'
          ]"
        >
          {{ page }}
        </button>

        <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="px-3 py-1 rounded border"
        >
          ถัดไป
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useArticles } from '~/composables/articles/useArticles'
import Card from '~/components/Card/Card.vue'
import Filter from '~/components/Searching/Filter.vue'

const {
  articles,
  categories,
  selectedCategory,
  searchTerm,
  loading,
  formatDate,
  updateSearchTerm,
  updateSelectedCategory,
  // pagination
  currentPage,
  totalPages,
  goToPage,
  nextPage,
  prevPage,
} = useArticles()
</script>
