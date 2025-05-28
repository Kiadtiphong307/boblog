<script setup lang="ts">
import type { Article } from '~/types/article'

interface Props {
  article: Article
  formatDate: (date: string) => string
  onDelete?: (slug: string) => void
}

const props = defineProps<Props>()
</script>

<template>
  <div class="bg-white rounded-xl shadow p-6 flex flex-col justify-between hover:shadow-md transition duration-200">
    <!-- Category -->
    <div class="mb-2">
      <span class="text-sm bg-blue-100 text-blue-700 px-2 py-1 rounded font-medium">
        {{ article.category?.name || 'ไม่มีหมวดหมู่' }}
      </span>
    </div>

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
      {{ article.content || '' }}
    </p>

    <!-- Author -->
    <p class="text-sm text-gray-500 mb-2">
      ผู้เขียน: <span class="font-medium">{{ article.author?.username || 'ไม่ทราบชื่อ' }}</span>
    </p>

    <!-- Tags -->
    <div class="mb-4">
      <span
        v-for="tag in article.tags"
        :key="tag.id"
        class="inline-block bg-gray-100 text-gray-700 text-xs px-2 py-1 rounded mr-2"
      >
        #{{ tag.name }}
      </span>
    </div>

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

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
