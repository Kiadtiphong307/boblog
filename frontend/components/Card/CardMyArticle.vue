<template>
  <div class="bg-white rounded-xl shadow p-6 flex flex-col justify-between hover:shadow-md transition duration-200">
    <!-- Title -->
    <h2 class="text-lg font-bold mb-1">
      {{ article.title }}
    </h2>

    <!-- Date -->
    <p class="text-xs text-gray-400 mb-2">
      {{ formatDate(article.created_at || article.createdAt || '') }}
    </p>

    <!-- Summary -->
    <p class="text-gray-700 text-sm mb-3 line-clamp-3">
      {{ article.content || 'ไม่มีคำอธิบาย' }}
    </p>

    <!-- Action Buttons -->
    <div class="mt-4 flex flex-wrap gap-2">
      <NuxtLink
        :to="`/articles/${article.slug}`"
        class="bg-green-600 hover:bg-green-700 text-white text-sm px-4 py-1.5 rounded transition"
      >
        👁 ดูบทความ
      </NuxtLink>

      <NuxtLink
        :to="`/articles/my-articles/${article.slug}`"
        class="bg-blue-600 hover:bg-blue-700 text-white text-sm px-4 py-1.5 rounded transition"
      >
        ✏️ แก้ไข
      </NuxtLink>

      <button
        @click="onDelete && onDelete(article.slug)"
        class="bg-red-600 hover:bg-red-700 text-white text-sm px-4 py-1.5 rounded transition"
      >
        🗑 ลบ
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Article } from '~/types/article'

interface Props {
  article: Article
  formatDate: (date: string) => string
  onDelete?: (slug: string) => void
}

const props = defineProps<Props>()
</script>
