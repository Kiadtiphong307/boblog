<script setup>
import { ref, onMounted } from 'vue'

const articles = ref([])
const loading = ref(false)
const error = ref(null)

const fetchMyArticles = async () => {
  loading.value = true
  error.value = null

  try {
    const token = localStorage.getItem('token')
    const res = await $fetch('/api/articles/my-articles', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    articles.value = res?.data || []
  } catch (err) {
    error.value = '❌ ไม่สามารถโหลดบทความของคุณได้'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchMyArticles()
})
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">📝 บทความของฉัน</h1>

    <div v-if="loading">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500">{{ error }}</div>

    <div v-if="articles.length === 0 && !loading">ยังไม่มีบทความ</div>

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="article in articles"
        :key="article.id"
        class="p-4 border rounded shadow-sm hover:shadow-md transition"
      >
        <h2 class="text-xl font-semibold">{{ article.title }}</h2>
        <p class="text-gray-600 mt-2 line-clamp-3">{{ article.description }}</p>
        <p class="text-sm text-gray-400 mt-4">🗓 {{ new Date(article.created_at).toLocaleDateString() }}</p>
      </div>
    </div>
  </div>
</template>
