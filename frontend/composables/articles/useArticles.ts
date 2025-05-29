import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { Article, Category } from '~/types/article'

// 📦 Read: ดึงรายการบทความและหมวดหมู่
export function useArticles() {
  // raw data
  const articles = ref<Article[]>([])
  const categories = ref<Category[]>([])
  // filters
  const selectedCategory = ref<string | number>('')
  const searchTerm = ref('')
  const localSearchTerm = ref('') // for immediate v-model
  const loading = ref(true)
  let debounceTimer: NodeJS.Timeout | null = null

  // pagination state
  const currentPage = ref(1)
  const perPage = ref(4) // จำนวนบทความต่อหน้า
  const totalPages = computed(() =>
    Math.ceil(articles.value.length / perPage.value) || 1
  )
  const paginatedArticles = computed(() => {
    const start = (currentPage.value - 1) * perPage.value
    return articles.value.slice(start, start + perPage.value)
  })

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page
    }
  }
  const nextPage = () => {
    if (currentPage.value < totalPages.value) currentPage.value++
  }
  const prevPage = () => {
    if (currentPage.value > 1) currentPage.value--
  }

  // fetch categories
  const fetchCategories = async () => {
    try {
      const res = await $fetch<{ data: Category[] }>('/api/categories')
      categories.value = res.data || []
    } catch (err) {
      console.error('❌ Failed to load categories:', err)
    }
  }

  // fetch articles (reset to page 1 เมื่อ filter/search เปลี่ยน)
  const fetchArticles = async () => {
    loading.value = true
    currentPage.value = 1

    const query = new URLSearchParams()
    if (searchTerm.value.trim()) query.append('search', searchTerm.value.trim())
    if (selectedCategory.value) query.append('category_id', selectedCategory.value.toString())

    try {
      const res = await $fetch<{ data: Article[] }>(`/api/articles?${query.toString()}`)
      articles.value = res.data || []
    } catch (err) {
      console.error('❌ Error loading articles:', err)
      articles.value = []
    } finally {
      loading.value = false
    }
  }

  // debounced search term setter
  const debouncedSearch = (value: string) => {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
      searchTerm.value = value
    }, 500)
  }
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

  onUnmounted(() => {
    if (debounceTimer) clearTimeout(debounceTimer)
  })
  onMounted(() => {
    fetchCategories()
    fetchArticles()
  })
  watch([searchTerm, selectedCategory], fetchArticles)

  return {
    // เปลี่ยนให้ articles ชี้ไปที่ paginatedArticles
    articles: paginatedArticles,
    categories,
    selectedCategory,
    searchTerm: localSearchTerm,
    loading,
    fetchArticles,
    formatDate,
    updateSearchTerm,
    updateSelectedCategory,
    // exports สำหรับ pagination
    currentPage,
    totalPages,
    goToPage,
    nextPage,
    prevPage,
  }
}

// 📝 Create: สำหรับสร้างบทความใหม่
export function useCreateArticle() {
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

