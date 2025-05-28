<template>
  <div class="max-w-6xl mx-auto py-10 px-4">
    <h1 class="text-2xl font-bold mb-6">📚 บทความทั้งหมด</h1>
    
    <!-- Filter Component -->
    <Filter
      :categories="categories"
      :search-term="searchTerm"
      :selected-category="selectedCategory"
      :update-search-term="updateSearchTerm"
      :update-selected-category="updateSelectedCategory"
    />
    
    <!-- Content -->
    <div v-if="loading" class="text-gray-500">⏳ กำลังโหลดบทความ...</div>
    <div v-else-if="articles.length === 0" class="text-red-500">❌ ไม่พบบทความ</div>
    <div v-else>
      <!-- Cards Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card
          v-for="article in articles"
          :key="article.id"
          :article="article"
          :format-date="formatDate"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProductList } from '~/composables/articles/useArticles'
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
} = useProductList()
</script>