import { ref, watch, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Article, Category } from '~/types/article'
import { usePagination } from '~/composables/usePagination'
import { useFilter } from '~/composables/useFilter'

// 📦 Read: ดึงรายการบทความและหมวดหมู่
export function useArticles() {
  // Raw data
  const allArticles = ref<Article[]>([])
  const categories = ref<Category[]>([])
  const loading = ref(true)
  const error = ref<string>('')
  
  // Setup filter functionality
  const filter = useFilter({ debounceDelay: 500 })
  
  // Setup pagination functionality
  const pagination = usePagination(allArticles, { perPage: 4 })
  
  // API methods
  const fetchCategories = async () => {
    try {
      const res = await $fetch<{ data: Category[] }>('/api/categories')
      categories.value = res.data || []
    } catch (err) {
      console.error('❌ Failed to load categories:', err)
      error.value = 'ไม่สามารถโหลดหมวดหมู่ได้'
    }
  }
  
  const fetchArticles = async () => {
    loading.value = true
    error.value = ''
    
    // Reset to first page when filters change
    pagination.resetToFirstPage()
    
    try {
      const queryParams = filter.buildQueryParams()
      const url = `/api/articles${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      
      const res = await $fetch<{ data: Article[] }>(url)
      allArticles.value = res.data || []
    } catch (err) {
      console.error('❌ Error loading articles:', err)
      allArticles.value = []
      error.value = 'ไม่สามารถโหลดบทความได้'
    } finally {
      loading.value = false
    }
  }
  
  // Utility functions
  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }
  
  const refreshArticles = async () => {
    await fetchArticles()
  }
  
  // Initialize data and watchers
  onMounted(async () => {
    await Promise.all([
      fetchCategories(),
      fetchArticles()
    ])
  })
  
  // Watch for filter changes
  watch([filter.searchTerm, filter.selectedCategory], fetchArticles)
  
  return {
    // Data
    articles: pagination.paginatedItems, // Return paginated articles
    allArticles, // Access to all articles if needed
    categories,
    loading,
    error,
    
    // Filter functionality
    selectedCategory: filter.selectedCategory,
    searchTerm: filter.localSearchTerm,
    updateSearchTerm: filter.updateSearchTerm,
    updateSelectedCategory: filter.updateSelectedCategory,
    clearAllFilters: filter.clearAllFilters,
    hasActiveFilters: filter.hasActiveFilters,
    getFilterSummary: filter.getFilterSummary,
    
    // Pagination functionality
    currentPage: pagination.currentPage,
    totalPages: pagination.totalPages,
    pageInfo: pagination.pageInfo,
    goToPage: pagination.goToPage,
    nextPage: pagination.nextPage,
    prevPage: pagination.prevPage,
    hasNextPage: pagination.hasNextPage,
    hasPrevPage: pagination.hasPrevPage,
    getPageNumbers: pagination.getPageNumbers,
    
    // Methods
    fetchArticles,
    refreshArticles,
    formatDate
  }
}

// 📝 Create: สำหรับสร้างบทความใหม่
export function useCreateArticle() {
  const router = useRouter()
  
  // Form state
  const formData = ref({
    title: '',
    slug: '',
    content: '',
    categoryName: '',
    tags: ''
  })
  
  const formState = ref({
    error: {} as Record<string, string>,
    success: false,
    loading: false
  })
  
  // Auto-generate slug from title
  const slugify = (text: string): string =>
    text
      .toLowerCase()
      .trim()
      .normalize('NFD')
      .replace(/[^\p{L}\p{N}\s-]/gu, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '')
  
  // Watch title changes to update slug
  watch(() => formData.value.title, (newTitle) => {
    formData.value.slug = slugify(newTitle)
  })
  
  // Form validation
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}
    
    if (!formData.value.title.trim()) {
      errors.title = 'กรุณาใส่ชื่อบทความ'
    }
    
    if (!formData.value.content.trim()) {
      errors.content = 'กรุณาใส่เนื้อหาบทความ'
    }
    
    if (!formData.value.categoryName.trim()) {
      errors.categoryName = 'กรุณาใส่ชื่อหมวดหมู่'
    }
    
    formState.value.error = errors
    return Object.keys(errors).length === 0
  }
  
  // Reset form
  const resetForm = () => {
    formData.value = {
      title: '',
      slug: '',
      content: '',
      categoryName: '',
      tags: ''
    }
    formState.value = {
      error: {},
      success: false,
      loading: false
    }
  }
  
  // Submit form
  const handleSubmit = async () => {
    formState.value.error = {}
    formState.value.success = false
    
    // Validate form
    if (!validateForm()) {
      return
    }
    
    // Check authentication
    const token = localStorage.getItem('token')
    if (!token) {
      formState.value.error.general = 'คุณต้องเข้าสู่ระบบก่อนสร้างบทความ'
      return
    }
    
    formState.value.loading = true
    
    try {
      const payload = {
        title: formData.value.title.trim(),
        slug: formData.value.slug.trim(),
        content: formData.value.content.trim(),
        category_name: formData.value.categoryName.trim(),
        tag_names: formData.value.tags 
          ? formData.value.tags.split(',').map(t => t.trim()).filter(Boolean) 
          : []
      }
      
      await $fetch('/api/articles', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: payload
      })
      
      formState.value.success = true
      
      // Reset form after success
      setTimeout(() => {
        resetForm()
        router.push('/articles')
      }, 1500)
      
    } catch (e: any) {
      if (e?.data?.errors) {
        formState.value.error = e.data.errors
      } else if (e?.data?.error) {
        formState.value.error.general = e.data.error
      } else {
        formState.value.error.general = '❌ เกิดข้อผิดพลาดในการสร้างบทความ'
      }
    } finally {
      formState.value.loading = false
    }
  }
  
  return {
    // Form data
    formData,
    
    // Individual form fields for easier v-model binding
    title: computed({
      get: () => formData.value.title,
      set: (value: string) => formData.value.title = value
    }),
    slug: computed({
      get: () => formData.value.slug,
      set: (value: string) => formData.value.slug = value
    }),
    content: computed({
      get: () => formData.value.content,
      set: (value: string) => formData.value.content = value
    }),
    categoryName: computed({
      get: () => formData.value.categoryName,
      set: (value: string) => formData.value.categoryName = value
    }),
    tags: computed({
      get: () => formData.value.tags,
      set: (value: string) => formData.value.tags = value
    }),
    
    // Form state
    error: computed(() => formState.value.error),
    success: computed(() => formState.value.success),
    loading: computed(() => formState.value.loading),
    
    // Methods
    handleSubmit,
    resetForm,
    validateForm
  }
}