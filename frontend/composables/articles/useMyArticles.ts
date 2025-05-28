import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Article } from '@/types/article'

export function useMyArticles() {
  const articles = ref<Article[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const router = useRouter()

  // ดึงบทความของฉัน
  const fetchMyArticles = async () => {
    loading.value = true
    error.value = null

    try {
      const token = localStorage.getItem('token')
      const res = await $fetch<{ data: Article[] }>('/api/articles/my-articles', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      articles.value = res.data || []
    } catch (err) {
      error.value = '❌ ไม่สามารถโหลดบทความของคุณได้'
    } finally {
      loading.value = false
    }
  }

  // ลบบทความ
  const deleteArticle = async (slug: string) => {
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
    } catch (err: any) {
      console.error('❌ ลบบทความไม่สำเร็จ:', err)
      const message = err?.data?.message || 'ลบบทความไม่สำเร็จ กรุณาลองใหม่'
      alert(`❌ ${message}`)
    }
  }

  const formatDateTime = (input: string): string => {
    const date = new Date(input)
    return date.toLocaleString('th-TH', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  onMounted(() => {
    fetchMyArticles()
  })

  return {
    articles,
    loading,
    error,
    fetchMyArticles,
    deleteArticle,
    formatDateTime,
  }
}
