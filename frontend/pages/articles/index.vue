<script setup lang="ts">
import { ref, onMounted, watch } from "vue";

const articles = ref([]);
const categories = ref([]);
const selectedCategory = ref("");
const searchTerm = ref("");
const loading = ref(true);

// Fetch Categories
const fetchCategories = async () => {
  try {
    const res = await $fetch("/api/categories");
    categories.value = res.data || [];
  } catch (err) {
    console.error("❌ Failed to load categories:", err);
  }
};

// โหลดบทความทั้งหมด (ไม่ใช้ pagination)
const fetchArticles = async () => {
  loading.value = true;
  const query = new URLSearchParams();

  if (searchTerm.value.trim()) {
    query.append("search", searchTerm.value.trim());
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

// Format Date
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("th-TH", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

onMounted(() => {
  fetchCategories();
  fetchArticles();
});

watch([searchTerm, selectedCategory], fetchArticles);
</script>

<template>
  <div class="max-w-6xl mx-auto py-10 px-4">
    <h1 class="text-2xl font-bold mb-6">📚 บทความทั้งหมด</h1>

    <!-- Filter -->
    <div class="flex flex-col md:flex-row md:items-center md:space-x-4 space-y-4 md:space-y-0 mb-8">
      <input
        v-model="searchTerm"
        type="text"
        placeholder="🔍 ค้นหาชื่อบทความหรือแท็ก"
        class="border px-4 py-2 rounded w-full md:w-2/3"
      />

      <select v-model="selectedCategory" class="border px-4 py-2 rounded w-full md:w-1/3">
        <option value="">📂 ทุกหมวดหมู่</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- Content -->
    <div v-if="loading" class="text-gray-500">⏳ กำลังโหลดบทความ...</div>
    <div v-else-if="articles.length === 0" class="text-red-500">❌ ไม่พบบทความ</div>
    <div v-else>
      <!-- Cards Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div
          v-for="article in articles"
          :key="article.id"
          class="bg-white rounded-xl shadow p-6 flex flex-col justify-between"
        >
          <!-- Category -->
          <div class="mb-2">
            <span class="text-sm bg-blue-100 text-blue-700 px-2 py-1 rounded font-medium">
              {{ article.category?.name || "ไม่มีหมวดหมู่" }}
            </span>
          </div>

          <!-- Article Title -->
          <h2 class="text-lg font-bold mb-1">
            {{ article.title }}
          </h2>

          <!-- Date -->
          <p class="text-xs text-gray-400 mb-2">
            {{ formatDate(article.created_at || article.createdAt) }}
          </p>

          <!-- Summary -->
          <p class="text-gray-700 text-sm mb-3">
            {{ (article.content || '').slice(0, 120) }}...
          </p>

          <!-- Author -->
          <p class="text-sm text-gray-500 mb-2">
            ผู้เขียน: <span class="font-medium">{{ article.author?.username || "ไม่ทราบชื่อ" }}</span>
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
      </div>
    </div>
  </div>
</template>
