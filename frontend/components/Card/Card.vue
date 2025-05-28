<template>
  <div class="bg-white rounded-xl shadow p-6 flex flex-col justify-between">
    <!-- Category -->
    <div class="mb-2">
      <span class="text-sm bg-blue-100 text-blue-700 px-2 py-1 rounded font-medium">
        {{ article.category?.name || 'ไม่มีหมวดหมู่' }}
      </span>
    </div>
    <!-- Article Title -->
    <h2 class="text-lg font-bold mb-1">
      {{ article.title }}
    </h2>
    <!-- Date -->
    <p class="text-xs text-gray-400 mb-2">
      {{ formatDate(article.created_at || article.createdAt || '') }}
    </p>
    <!-- Summary -->
    <p class="text-gray-700 text-sm mb-3">
      {{ (article.content || '').slice(0, 120) }}...
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
    <!-- Read More Button -->
    <NuxtLink
      :to="`/articles/${article.slug}`"
      class="self-start inline-block bg-blue-500 text-white text-sm px-4 py-2 rounded hover:bg-blue-600 transition"
    >
      อ่านเพิ่มเติม
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import type { Article } from '~/types/article'

interface Props {
  article: Article
  formatDate: (date: string) => string
}

defineProps<Props>()
</script>