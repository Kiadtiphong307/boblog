<script setup>
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

onMounted(() => {
  fetchMyArticles()
})
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">📝 บทความของฉัน</h1>

    <div v-if="loading">⏳ กำลังโหลดบทความ...</div>
    <div v-if="error" class="text-red-500">{{ error }}</div>
    <div v-if="articles.length === 0 && !loading">ยังไม่มีบทความ</div>

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="article in articles" :key="article.id" class="p-4 border rounded shadow hover:shadow-md transition">
        <h2 class="text-xl font-semibold">{{ article.title }}</h2>
        <p class="text-gray-600 mt-2 line-clamp-3">{{ article.description }}</p>
        <p class="text-sm text-gray-400 mt-4">
          🗓 {{ new Date(article.created_at).toLocaleDateString() }}
        </p>

        <!-- ปุ่มต่าง ๆ -->
        <div class="mt-4 flex gap-2 flex-wrap">
          <!-- 👁 ดูบทความ -->
          <NuxtLink :to="`/articles/${article.slug}`"
            class="px-3 py-1 bg-green-600 text-white rounded hover:bg-green-700">
            👁 ดูบทความ
          </NuxtLink>

          <!-- ✏️ แก้ไข -->
          <NuxtLink :to="`/articles/edit/${article.slug}`"
            class="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700">
            ✏️ แก้ไข
          </NuxtLink>

          <!-- 🗑 ลบ -->
          <button @click="deleteArticle(article.slug)" class="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700">
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
