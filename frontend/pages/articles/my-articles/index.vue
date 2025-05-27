เิอออออออออออออออ <script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const articles = ref([])
const loading = ref(false)
const error = ref(null)
const router = useRouter()

const fetchMyArticles = async () => {
  loading.value = true
  error.value = null

  try {
    const token = localStorage.getItem('token')
    const res = await $fetch('/api/articles/my-articles', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    articles.value = res?.data || []
  } catch (err) {
    error.value = '❌ ไม่สามารถโหลดบทความของคุณได้'
  } finally {
    loading.value = false
  }
}

const deleteArticle = async (slug) => {
  const confirmed = confirm('คุณแน่ใจหรือไม่ว่าต้องการลบบทความนี้?')
  if (!confirmed) return

  try {
    const token = localStorage.getItem('token')
    if (!token) {
      alert('❌ ไม่พบ token โปรดเข้าสู่ระบบอีกครั้ง')
      return
    }

    await $fetch(`/api/articles/${slug}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    articles.value = articles.value.filter((a) => a.slug !== slug)
    alert('✅ ลบบทความเรียบร้อยแล้ว')
  } catch (err) {
    console.error('❌ ลบบทความไม่สำเร็จ:', err)
    const message = err?.data?.message || 'ลบบทความไม่สำเร็จ กรุณาลองใหม่'
    alert(`❌ ${message}`)
  }
}

const formatDateTime = (input) => {
  const date = new Date(input)
  return date.toLocaleString('th-TH', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })

  .465278

  ล

















  ฃฃ
}

onMounted(() => {
  fetchMyArticles()
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">📝 บทความของฉัน</h1>

    <div v-if="loading" class="text-gray-600">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500 font-semibold">{{ error }}</div>
    <div v-if="articles.length === 0 && !loading" class="text-gray-500">ยังไม่มีบทความ</div>

    <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="article in articles"
        :key="article.id"
        class="bg-white border border-gray-200 rounded-2xl shadow-sm p-6 hover:shadow-md transition duration-200"
      >
        <h2 class="text-xl font-semibold text-gray-800 mb-2">{{ article.title }}</h2>
        <p class="text-gray-600 text-sm line-clamp-3 mb-4">
          {{ article.description }}
        </p>
        <p class="text-xs text-gray-400">
          🗓 {{ formatDateTime(article.created_at) }}
        </p>
        

        <div class="mt-6 flex flex-wrap gap-2">
          <NuxtLink
            :to="`/articles/${article.slug}`"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-green-600 text-white hover:bg-green-700 transition"
          >
            👁 ดูบทความ
          </NuxtLink>

          <NuxtLink
            :to="`/articles/edit/${article.slug}`"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition"
          >
            ✏️ แก้ไข
          </NuxtLink>

          <button
            @click="deleteArticle(article.slug)"
            class="inline-block px-4 py-1.5 text-sm rounded-xl bg-red-600 text-white hover:bg-red-700 transition"
          >
            🗑 ลบ
          </button>
        </div>
      </div>
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
