<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">📝 บทความของฉัน</h1>

    <div v-if="loading" class="text-gray-600">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500 font-semibold">{{ error }}</div>
    <div v-if="articles.length === 0 && !loading" class="text-gray-500">ยังไม่มีบทความ</div>

    <!-- แสดงจำนวนบทความ
    <div v-if="articles.length > 0 && !loading" class="mb-6 text-gray-600">
      แสดง {{ pageInfo.showing.from }}-{{ pageInfo.showing.to }} จาก {{ pageInfo.showing.total }} บทความ
    </div> -->

    <!-- Articles Grid -->
    <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
      <Card
        v-for="article in paginatedItems"
        :key="article.id"
        :article="article"
        :formatDate="formatDateTime"
        :onDelete="confirmDelete"
      />
    </div>

    <!-- Pagination Controls -->
    <div v-if="totalPages > 1" class="flex flex-col items-center space-y-4">
      <!-- Navigation Buttons -->
      <div class="flex items-center space-x-2">
        <button
          @click="prevPage"
          :disabled="!hasPrevPage"
          class="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300 disabled:cursor-not-allowed hover:bg-blue-600"
        >
          ← ก่อนหน้า
        </button>
        
        <!-- Page Numbers -->
        <div class="flex space-x-1">
          <button
            v-for="page in getPageNumbers()"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-3 py-2 rounded',
              page === currentPage 
                ? 'bg-blue-500 text-white font-bold' 
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            {{ page }}
          </button>
        </div>
        
        <button
          @click="nextPage"
          :disabled="!hasNextPage"
          class="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300 disabled:cursor-not-allowed hover:bg-blue-600"
        >
          ถัดไป →
        </button>
      </div>

      <!-- Items per page selector -->
      <div class="flex items-center space-x-2 text-sm text-gray-600">
        <label for="itemsPerPage">แสดงต่อหน้า:</label>
        <select
          id="itemsPerPage"
          v-model="itemsPerPage"
          @change="setItemsPerPage(itemsPerPage)"
          class="border border-gray-300 rounded px-2 py-1"
        >
          <option :value="3">3</option>
          <option :value="6">6</option>
          <option :value="9">9</option>
          <option :value="12">12</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from '~/components/Card/CardMyArticle.vue'
import { useMyArticles } from '~/composables/articles/useMyArticles'
import { usePagination } from '~/composables/usePagination'

const {
  articles,
  loading,
  error,
  deleteArticle,
  formatDateTime,
} = useMyArticles()

// Setup pagination with 6 items per page by default
const {
  currentPage,
  itemsPerPage,
  totalPages,
  paginatedItems,
  hasNextPage,
  hasPrevPage,
  pageInfo,
  goToPage,
  nextPage,
  prevPage,
  setItemsPerPage,
  getPageNumbers
} = usePagination(articles, { perPage: 6 })

const confirmDelete = async (slug: string) => {
  if (confirm('คุณแน่ใจหรือไม่ว่าต้องการลบบทความนี้?')) {
    await deleteArticle(slug)
    // หากลบแล้วหน้าปัจจุบันไม่มี items ให้กลับไปหน้าก่อนหน้า
    if (paginatedItems.value.length === 0 && currentPage.value > 1) {
      prevPage()
    }
  }
}
</script>