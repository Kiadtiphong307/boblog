<script setup>
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

  if (searchTerm.value) query.append("search", searchTerm.value);
  if (selectedCategory.value)
    query.append("category_id", selectedCategory.value);

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
  fetchArticles();
});

// fetch ใหม่ทุกครั้งที่ filter เปลี่ยน
watch([searchTerm, selectedCategory], fetchArticles);
</script>

<template>
  <div class="max-w-4xl mx-auto py-8">
    <h1 class="text-2xl font-bold mb-6">📚 Articles</h1>

    <!-- Filters -->
    <div
      class="flex flex-col md:flex-row md:items-center md:space-x-4 space-y-4 md:space-y-0 mb-6"
    >
      <input
        v-model="searchTerm"
        type="text"
        placeholder="🔍 Search articles..."
        class="border px-4 py-2 rounded w-full md:w-1/2"
      />
      <select
        v-model="selectedCategory"
        class="border px-4 py-2 rounded w-full md:w-1/3"
      >
        <option value="">All Categories</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- Articles List -->
    <div v-if="loading" class="text-gray-500">Loading articles...</div>
    <div v-else-if="articles.length === 0" class="text-gray-500">
      No articles found.
    </div>
    <div v-else>
      <div v-for="article in articles" :key="article.id" class="border-b py-6">
        <NuxtLink
          :to="`/articles/${article.slug}`"
          class="text-xl font-semibold text-blue-600 hover:underline"
        >
          {{ article.title }}
        </NuxtLink>

        <div class="text-sm text-gray-500 mt-1">
          👤 {{ article.author?.username || "Unknown" }} | 📂
          {{ article.category?.name || "Uncategorized" }} | 🏷️
          <span
            v-for="tags in article.tags"
            :key="tags.id"
            class="bg-gray-200 rounded px-2 py-0.5 text-xs mr-1"
          >
            #{{ tags.name }}
          </span>
        </div>

        <div class="text-gray-700 mt-3">
          {{ (article.content || "").slice(0, 120) }}...
        </div>

        <div class="text-xs text-gray-400 mt-2">
          🕒 {{ new Date(article.created_at).toLocaleString() }}
        </div>
      </div>
    </div>
  </div>
</template>
