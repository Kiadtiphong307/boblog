<script setup lang="ts">
import { ref, onMounted, watch } from "vue";

const articles = ref([]);
const categories = ref([]);
const tags = ref([]);

const selectedCategory = ref("");
const selectedTag = ref("");
const searchTerm = ref("");
const loading = ref(true);

// ✅ โหลดหมวดหมู่
const fetchCategories = async () => {
  try {
    const res = await $fetch("/api/categories");
    categories.value = res.data || [];
  } catch (err) {
    console.error("❌ Failed to load categories:", err);
  }
};

// ✅ โหลดแท็ก
const fetchTags = async () => {
  try {
    const res = await $fetch("/api/tags");
    tags.value = res.data || [];
  } catch (err) {
    console.error("❌ Failed to load tags:", err);
  }
};

// ✅ โหลดบทความ
const fetchArticles = async () => {
  loading.value = true;
  const query = new URLSearchParams();

  if (searchTerm.value) query.append("search", searchTerm.value);
  if (selectedCategory.value) query.append("category_id", selectedCategory.value);
  if (selectedTag.value) query.append("tag", selectedTag.value);

  try {
    const res = await $fetch(`/api/articles?${query.toString()}`);
    articles.value = res.data || [];
  } catch (err) {
    console.error("❌ Error loading articles:", err);
    articles.value = [];
  }

  loading.value = false;
};

onMounted(() => {
  fetchCategories();
  fetchTags();
  fetchArticles();
});

watch([searchTerm, selectedCategory, selectedTag], fetchArticles);
</script>

<template>
  <div class="max-w-4xl mx-auto py-8">
    <h1 class="text-2xl font-bold mb-6">📚 บทความทั้งหมด</h1>

    <!-- Filters -->
    <div class="flex flex-col md:flex-row md:items-center md:space-x-4 space-y-4 md:space-y-0 mb-6">
      <input
        v-model="searchTerm"
        type="text"
        placeholder="🔍 ค้นหาบทความ..."
        class="border px-4 py-2 rounded w-full md:w-1/2"
      />

      <select v-model="selectedCategory" class="border px-4 py-2 rounded w-full md:w-1/4">
        <option value="">📂 ทุกหมวดหมู่</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>

      <select v-model="selectedTag" class="border px-4 py-2 rounded w-full md:w-1/4">
        <option value="">🏷️ ทุกแท็ก</option>
        <option v-for="tag in tags" :key="tag.id" :value="tag.name">
          #{{ tag.name }}
        </option>
      </select>
    </div>

    <!-- Articles List -->
    <div v-if="loading" class="text-gray-500">⏳ กำลังโหลดบทความ...</div>
    <div v-else-if="articles.length === 0" class="text-gray-500">❌ ไม่พบบทความที่ตรงกัน</div>
    <div v-else>
      <div
        v-for="article in articles"
        :key="article.id"
        class="border-b py-6"
      >
        <NuxtLink :to="`/articles/${article.slug}`" class="text-xl font-semibold text-blue-600 hover:underline">
          {{ article.title }}
        </NuxtLink>

        <div class="text-sm text-gray-500 mt-1">
          👤 {{ article.author?.username || "Unknown" }} |
          📂 {{ article.category?.name || "ไม่มีหมวดหมู่" }} |
          🏷️
          <span
            v-for="tag in article.tags"
            :key="tag.id"
            class="bg-gray-200 rounded px-2 py-0.5 text-xs mr-1"
          >
            #{{ tag.name }}
          </span>
        </div>

        <div class="text-gray-700 mt-3">
          {{ (article.content || "").slice(0, 120) }}...
        </div>

        <div class="text-xs text-gray-400 mt-2">
          🕒 {{ new Date(article.created_at).toLocaleString() }}
        </div>

        <div class="mt-2">
          <NuxtLink
            :to="`/articles/${article.slug}`"
            class="text-sm text-white bg-blue-600 hover:bg-blue-700 px-4 py-1 rounded inline-block"
          >
            อ่านต่อ
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
