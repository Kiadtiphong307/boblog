<script setup lang="ts">
import { useMyArticles } from '~/composables/articles/useMyArticles'

const {
  articles,
  loading,
  error,
  fetchMyArticles,
  deleteArticle,
  formatDateTime,
} = useMyArticles()
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">📝 บทความของฉัน</h1>

    <!-- Loading/Error -->
    <div v-if="loading" class="text-gray-600">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500 font-semibold">{{ error }}</div>
    <div v-if="articles.length === 0 && !loading" class="text-gray-500">ยังไม่มีบทความ</div>

    <!-- Article Cards -->
    <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="article in articles"
        :key="article.id"
        class="bg-white border border-gray-200 rounded-2xl shadow-sm p-6 hover:shadow-md transition duration-200"
      >
        <h2 class="text-xl font-semibold text-gray-800 mb-2">{{ article.title }}</h2>
        <p class="text-gray-600 text-sm line-clamp-3 mb-4">
          {{ article.content }}
        </p>
        <p class="text-xs text-gray-400">
          🗓 {{ formatDateTime(article.created_at) }}
        </p>

        <!-- Actions -->
        <div class="mt-6 flex flex-wrap gap-2">
          <NuxtLink
            :to="`/articles/${article.slug}`"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-green-600 text-white hover:bg-green-700 transition"
          >
            👁 ดูบทความ
          </NuxtLink>

          <NuxtLink
            :to="`/articles/my-articles/${article.slug}`"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition"
          >
            ✏️ แก้ไข
          </NuxtLink>

          <button
            @click="deleteArticle(article.slug)"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-red-600 text-white hover:bg-red-700 transition"
          >
            🗑 ลบ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
