<template>
  <div class="max-w-3xl mx-auto py-8 px-4">
    <!-- Title -->
    <h1 class="text-3xl font-bold mb-4">{{ article?.title }}</h1>
    <!-- Author, Category, Date -->
    <div class="text-sm text-gray-500 mb-2">
      👤 {{ article?.author?.username }} |
      📂 {{ article?.category?.name }} |
      🕒 {{ formatDate(article?.created_at || '') }}
    </div>
    <div class="mb-4">
      <span
        v-for="tag in article?.tags || []"
        :key="tag.id"
        class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded mr-2"
      >
        #{{ tag.name }}
      </span>
    </div>
    <div class="prose max-w-none" v-html="article?.content"></div>

    <!-- Comments -->
    <hr class="my-6" />
    <h2 class="text-lg font-semibold mb-2">💬 Comments ({{ comments.length }})</h2>
    <ul class="space-y-3 mb-6">
      <li
        v-for="c in comments"
        :key="c.id"
        class="bg-gray-100 p-3 rounded text-sm text-gray-800"
      >
        <div class="font-semibold text-blue-700 mb-1">
          👤 {{ c.user?.username || 'ไม่ทราบชื่อ' }}
          <span class="text-gray-500 text-xs ml-2">🕒 {{ formatDate(c.created_at) }}</span>
        </div>
        <div>{{ c.content }}</div>
      </li>
    </ul>

    <!-- Comment Form -->
    <form @submit.prevent="submitComment" class="space-y-2">
      <textarea
        v-model="newComment"
        rows="3"
        class="w-full border rounded p-2"
        placeholder="แสดงความคิดเห็น..."
      ></textarea>
      <button
        type="submit"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
       แสดงความคิดเห็น
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useProductDetail } from '~/composables/articles/useArticles'

const {
  article,
  comments,
  newComment,
  fetchArticle,
  fetchComments,
  submitComment,
  formatDate,
} = useProductDetail()
</script>
