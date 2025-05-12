
<script setup>
import { ref, onMounted } from 'vue'

const articles = ref([])
const loading = ref(true)

const fetchArticles = async () => {
  loading.value = true
  const token = localStorage.getItem('token')

  try {
    const res = await $fetch('/api/articles', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    articles.value = res.data || []
  } catch (err) {
    console.error('❌ Error loading articles:', err)
    articles.value = []
  }

  loading.value = false
}

onMounted(fetchArticles)
</script>



<template>
  <div class="max-w-4xl mx-auto py-8">
    <h1 class="text-2xl font-bold mb-4">📚 Articles</h1>

    <NuxtLink to="/articles/create" class="bg-green-600 text-white px-4 py-2 rounded mb-6 inline-block">
      + Create New Article
    </NuxtLink>

    <div v-if="loading" class="text-gray-500">Loading articles...</div>
    <div v-else-if="articles.length === 0" class="text-gray-500">No articles found.</div>

    <div v-else>
      <div
        v-for="article in articles"
        :key="article.id"
        class="border-b py-6"
      >
        <NuxtLink
          :to="`/articles/${article.slug}`"
          class="text-xl font-semibold text-blue-600 hover:underline"
        >
          {{ article.title }}
        </NuxtLink>

        <div class="text-sm text-gray-500 mt-1">
          👤 {{ article.author?.username || 'Unknown Author' }} |
          📂 {{ article.category?.name || 'Uncategorized' }} |
          🏷️
          <span v-for="tag in article.tags" :key="tag.id" class="inline-block mr-1">
            <span class="bg-gray-200 rounded px-2 py-0.5 text-xs">#{{ tag.name }}</span>
          </span>
        </div>

        <div class="text-gray-700 mt-3">
          {{ (article.content || '').slice(0, 120) }}...
        </div>

        <div class="text-xs text-gray-400 mt-2">
          🕒 {{ new Date(article.created_at).toLocaleString() }}
        </div>
      </div>
    </div>
  </div>
</template>
