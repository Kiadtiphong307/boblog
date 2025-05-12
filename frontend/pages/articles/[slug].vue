<script setup>
import { useRoute } from "vue-router";
import { ref, onMounted } from "vue";

const route = useRoute();
const article = ref(null);

const fetchArticle = async () => {
  const token = localStorage.getItem('token') // ✅ ดึง token ที่ได้จาก login

  try {
    const res = await $fetch(`/api/articles/${route.params.slug}`, {
      headers: {
        Authorization: `Bearer ${token}` // ✅ แนบ token
      }
    })
    article.value = res.data
  } catch (err) {
    console.error('❌ Failed to load article:', err)
  }
}


const formatDate = (isoString) => {
  return new Date(isoString).toLocaleString();
};

onMounted(fetchArticle);
</script>

<template>
  <div class="max-w-3xl mx-auto py-8">
    <h1 class="text-3xl font-bold mb-4">{{ article?.title }}</h1>

    <div class="text-sm text-gray-500 mb-2">
      👤 {{ article?.author?.username }} | 📂 {{ article?.category?.name }} | 🕒
      {{ formatDate(article?.created_at) }}
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

    <hr class="my-6" />
    <h2 class="text-lg font-semibold">
      💬 Comments ({{ article?.comments?.length || 0 }})
    </h2>
    <ul class="mt-2 space-y-3">
      <li
        v-for="c in article?.comments"
        :key="c.id"
        class="text-sm text-gray-700"
      >
        - {{ c.content }}
      </li>
    </ul>
  </div>
</template>
