<script setup>
import { ref, onMounted } from "vue";

const articles = ref([]);
const loading = ref(false);

const fetchMyArticles = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem("token");
    const res = await $fetch("/api/articles/my-articles", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    console.log("📦 Articles:", res.data);
    articles.value = res.data || [];
  } catch (error) {
    console.error("❌ ดึงบทความล้มเหลว:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchMyArticles();
});
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">บทความของฉัน</h1>

    <div v-if="loading">กำลังโหลด...</div>
    <div v-else-if="articles.length === 0">ยังไม่มีบทความ</div>
    <div v-else class="grid gap-4">
      <div
        v-for="article in articles"
        :key="article.id"
        class="border p-4 rounded-lg shadow hover:shadow-lg transition"
      >
        <h2 class="text-xl font-semibold text-blue-600">{{ article.title }}</h2>

        <p class="text-gray-700 mt-2">
          {{
            article.content.length > 150
              ? article.content.substring(0, 150) + "..."
              : article.content
          }}
        </p>

        <div class="text-sm text-gray-500 mt-2">
          ✍️ โดย
          {{
            article.author?.username || article.author?.nickname || "ไม่ระบุ"
          }}
          | 🏷️ หมวด {{ article.category?.name || "ไม่ระบุ" }}
        </div>

        <div
          v-if="article.tags?.length"
          class="mt-2 flex flex-wrap gap-2 text-sm"
        >
          <span
            v-for="tag in article.tags"
            :key="tag.id"
            class="bg-blue-100 text-blue-800 px-2 py-1 rounded"
          >
            #{{ tag.name }}
          </span>
        </div>

        <NuxtLink
          :to="`/articles/${article.slug}`"
          class="mt-4 inline-block text-blue-500 hover:underline"
        >
          อ่านเพิ่มเติม →
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
