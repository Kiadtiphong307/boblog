<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">📝 บทความของฉัน</h1>

    <div v-if="loading" class="text-gray-600">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500 font-semibold">{{ error }}</div>
    <div v-if="articles.length === 0 && !loading" class="text-gray-500">ยังไม่มีบทความ</div>

    <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card
        v-for="article in articles"
        :key="article.id"
        :article="article"
        :formatDate="formatDateTime"
        :onDelete="confirmDelete"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import Card from '~/components/Card/CardMyArticle.vue'
import { useMyArticles } from '~/composables/articles/useMyArticles'

const {
  articles,
  loading,
  error,
  deleteArticle,
  formatDateTime,
} = useMyArticles()

const confirmDelete = async (slug: string) => {
  if (confirm('คุณแน่ใจหรือไม่ว่าต้องการลบบทความนี้?')) {
    await deleteArticle(slug)
  }
}
</script>
