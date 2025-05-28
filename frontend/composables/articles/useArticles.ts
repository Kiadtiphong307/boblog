import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { Article, Category, Comment } from '~/types/article'

// 📦 Read: ดึงรายการบทความและหมวดหมู่
export function useProductList() {
  const articles = ref<Article[]>([])
  const categories = ref<Category[]>([])
  const selectedCategory = ref<string | number>('')
  const searchTerm = ref('')
  const localSearchTerm = ref('') // For immediate UI updates
  const loading = ref(true)
  
  let debounceTimer: NodeJS.Timeout | null = null

  // ดึงหมวดหมู่ทั้งหมด
  const fetchCategories = async () => {
    try {
      const res = await $fetch<{data: Category[]}>('/api/categories')
      categories.value = res.data || []
    } catch (err) {
      console.error('❌ Failed to load categories:', err)
    }
  }

  // ดึงบทความทั้งหมดโดยใช้ filter จาก searchTerm และ selectedCategory
  const fetchArticles = async () => {
    loading.value = true
    const query = new URLSearchParams()
    if (searchTerm.value.trim()) query.append('search', searchTerm.value.trim())
    if (selectedCategory.value) query.append('category_id', selectedCategory.value.toString())

    try {
      const res = await $fetch<{data: Article[]}>(`/api/articles?${query.toString()}`)
      articles.value = res.data || []
    } catch (err) {
      console.error('❌ Error loading articles:', err)
      articles.value = []
    } finally {
      loading.value = false
    }
  }

  // Debounced search function
  const debouncedSearch = (value: string) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    
    debounceTimer = setTimeout(() => {
      searchTerm.value = value
    }, 500) // Wait 500ms after user stops typing
  }

  // Update search functions
  const updateSearchTerm = (value: string) => {
    localSearchTerm.value = value
    debouncedSearch(value)
  }

  const updateSelectedCategory = (value: string | number) => {
    selectedCategory.value = value
  }

  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  // Cleanup on unmount
  onUnmounted(() => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
  })

  onMounted(() => {
    fetchCategories()
    fetchArticles()
  })

  // Watch for actual searchTerm changes (after debounce) and selectedCategory changes
  watch([searchTerm, selectedCategory], fetchArticles)

  return {
    articles,
    categories,
    selectedCategory,
    searchTerm: localSearchTerm, // Return localSearchTerm for immediate UI updates
    loading,
    fetchArticles,
    formatDate,
    updateSearchTerm,
    updateSelectedCategory,
  }
}

// 📝 Create: สำหรับสร้างบทความใหม่
export function useProductForm() {
  const router = useRouter()
  const title = ref('')
  const slug = ref('')
  const content = ref('')
  const categoryName = ref('')
  const tags = ref('')
  const error = ref<Record<string, string>>({})
  const success = ref(false)

  // สร้าง slug อัตโนมัติจาก title
  const slugify = (text: string): string =>
    text
      .toLowerCase()
      .trim()
      .normalize('NFD')
      .replace(/[^\p{L}\p{N}\s-]/gu, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '')

  watch(title, (newTitle) => {
    slug.value = slugify(newTitle)
  })

  // ส่งคำขอเพื่อสร้างบทความใหม่
  const handleSubmit = async () => {
    error.value = {}
    success.value = false

    const token = localStorage.getItem('token')
    if (!token) {
      error.value.general = 'คุณต้องเข้าสู่ระบบก่อนสร้างบทความ'
      return
    }

    try {
      await $fetch('/api/articles', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: {
          title: title.value,
          slug: slug.value,
          content: content.value,
          category_name: categoryName.value,
          tag_names: tags.value ? tags.value.split(',').map((t) => t.trim()).filter(Boolean) : [],
        },
      })

      success.value = true
      title.value = ''
      slug.value = ''
      content.value = ''
      categoryName.value = ''
      tags.value = ''

      setTimeout(() => router.push('/articles'), 1500)
    } catch (e: any) {
      if (e?.data?.errors) {
        error.value = e.data.errors
      } else if (e?.data?.error) {
        error.value.general = e.data.error
      } else {
        error.value.general = '❌ เกิดข้อผิดพลาดในการสร้างบทความ'
      }
    }
  }

  return {
    title,
    slug,
    content,
    categoryName,
    tags,
    error,
    success,
    handleSubmit,
  }
}

// 🔍 Read + 💬 Comment: แสดงรายละเอียดบทความและจัดการความคิดเห็น
export function useProductDetail() {
  const route = useRoute()
  const article = ref<Article | null>(null)
  const comments = ref<Comment[]>([])
  const newComment = ref('')

  // ดึงรายละเอียดบทความตาม slug
  const fetchArticle = async () => {
    const token = localStorage.getItem('token')
    try {
      const res = await $fetch<{data: Article}>(`/api/articles/${route.params.slug}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      article.value = (res as any).data
    } catch (err) {
      console.error('❌ Failed to load article:', err)
    }
  }

  // ดึงความคิดเห็นของบทความนั้น ๆ
  const fetchComments = async () => {
    const token = localStorage.getItem('token')
    try {
      const res = await $fetch(`/api/articles/${route.params.slug}/comments`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      comments.value = (res as any).data || []
    } catch (err) {
      console.error('❌ Failed to load comments:', err)
    }
  }

  // ส่งความคิดเห็นใหม่
  const submitComment = async () => {
    const token = localStorage.getItem('token')
    if (!token) return alert('กรุณาเข้าสู่ระบบก่อนแสดงความคิดเห็น')
    if (!newComment.value.trim()) return

    try {
      await $fetch(`/api/articles/${route.params.slug}/comments`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: {
          content: newComment.value,
        },
      })
      newComment.value = ''
      fetchComments()
    } catch (err) {
      console.error('❌ Failed to post comment:', err)
    }
  }

  const formatDate = (iso: string) =>
    new Date(iso).toLocaleString('th-TH', {
      dateStyle: 'medium',
      timeStyle: 'short',
    })

  onMounted(() => {
    fetchArticle()
    fetchComments()
  })

  return {
    article,
    comments,
    newComment,
    fetchArticle,
    fetchComments,
    submitComment,
    formatDate,
  }
}