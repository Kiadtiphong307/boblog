// ✅ FRONTEND: Nuxt 3 (Vue) ค้นหา articles จากชื่อหรือ tags หลายคำ
<script setup lang="ts">
import { ref, onMounted, watch } from "vue";

const articles = ref([]);
const categories = ref([]);
const selectedCategory = ref("");
const searchTerm = ref("");
const loading = ref(true);

const fetchCategories = async () => {
  try {
    const res = await $fetch("/api/categories");
    categories.value = res.data || [];
  } catch (err) {
    console.error("❌ Failed to load categories:", err);
  }
};

const fetchArticles = async () => {
  loading.value = true;
  const query = new URLSearchParams();

  if (searchTerm.value.trim()) {
    query.append("search", searchTerm.value.trim()); // ✅ ส่งไปเป็น search keyword เดียว
  }

  if (selectedCategory.value) {
    query.append("category_id", selectedCategory.value);
  }

  try {
    const res = await $fetch(`/api/articles?${query.toString()}`);
    articles.value = res.data || [];
  } catch (err) {
    console.error("❌ Error loading articles:", err);
    articles.value = [];
  } finally {
    loading.value = false;
  }
};


onMounted(() => {
  fetchCategories();
  fetchArticles();
});

watch([searchTerm, selectedCategory], fetchArticles);
</script>

<template>
  <div class="max-w-4xl mx-auto py-8">
    <h1 class="text-2xl font-bold mb-6">📚 บทความทั้งหมด</h1>

    <div class="flex flex-col md:flex-row md:items-center md:space-x-4 space-y-4 md:space-y-0 mb-6">
      <input
        v-model="searchTerm"
        type="text"
        placeholder="🔍 ค้นหาชื่อบทความหรือแท็ก (เว้นวรรคคั่น)"
        class="border px-4 py-2 rounded w-full md:w-2/3"
      />

      <select v-model="selectedCategory" class="border px-4 py-2 rounded w-full md:w-1/3">
        <option value="">📂 ทุกหมวดหมู่</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <div v-if="loading" class="text-gray-500">⏳ กำลังโหลดบทความ...</div>
    <div v-else-if="articles.length === 0" class="text-red-500">❌ ไม่พบบทความ</div>
    <div v-else>
      <div v-for="article in articles" :key="article.id" class="border-b py-6">
        <NuxtLink :to="`/articles/${article.slug}`" class="text-xl font-semibold text-blue-600 hover:underline">
          {{ article.title }}
        </NuxtLink>
        <div class="text-sm text-gray-500 mt-1">
          👤 {{ article.author?.username || "Unknown" }} |
          📂 {{ article.category?.name || "ไม่มีหมวดหมู่" }} |
          🏷️ <span v-for="tag in article.tags" :key="tag.id" class="bg-gray-200 rounded px-2 py-0.5 text-xs mr-1">
            #{{ tag.name }}
          </span>
        </div>
        <div class="text-gray-700 mt-3">
          {{ (article.content || '').slice(0, 120) }}...
        </div>
      </div>
    </div>
  </div>
</template>
