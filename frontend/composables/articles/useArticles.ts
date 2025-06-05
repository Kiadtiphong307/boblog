import { ref, watch, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Article, Category } from '~/types/article'
import { usePagination } from '~/composables/usePagination'
import { useFilter } from '~/composables/useFilter'


// Format date to Thai locale string
export const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// Convert text to URL-friendly slug
export const slugify = (text: string): string =>
  text
    .toLowerCase()
    .trim()
    .normalize('NFD')
    .replace(/[^\p{L}\p{N}\s-]/gu, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')

// Fetch categories from API
export const fetchCategoriesAPI = async () => {
  try {
    const res = await $fetch<{ data: Category[] }>('/api/categories')
    return { data: res.data || [], error: null }
  } catch (err) {
    console.error('❌ Failed to load categories:', err)
    return { data: [], error: 'ไม่สามารถโหลดหมวดหมู่ได้' }
  }
}

// Fetch articles from API with query parameters
export const fetchArticlesAPI = async (queryParams: URLSearchParams) => {
  try {
    const url = `/api/articles${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
    const res = await $fetch<{ data: Article[] }>(url)
    return { data: res.data || [], error: null }
  } catch (err) {
    console.error('❌ Error loading articles:', err)
    return { data: [], error: 'ไม่สามารถโหลดบทความได้' }
  }
}

// Create new article via API
export const createArticleAPI = async (payload: any, token: string) => {
  try {
    const response = await $fetch('/api/articles', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: payload
    })
    return { success: true, error: null }
  } catch (e: any) {
    let error = '❌ เกิดข้อผิดพลาดในการสร้างบทความ'
    
    if (e?.data?.errors) {
      return { success: false, error: e.data.errors }
    } else if (e?.data?.error) {
      error = e.data.error
    }
    
    return { success: false, error: { general: error } }
  }
}

// Article List Management
export function useArticleListState() {
  const allArticles = ref<Article[]>([])
  const categories = ref<Category[]>([])
  const loading = ref(true)
  const error = ref<string>('')
  
  const filter = useFilter({ debounceDelay: 500 })
  const pagination = usePagination(allArticles, { perPage: 4 })
  
  // Load categories
  const loadCategories = async () => {
    const { data, error: apiError } = await fetchCategoriesAPI()
    categories.value = data
    if (apiError) error.value = apiError
  }
  
  // Load articles with current filters
  const loadArticles = async () => {
    loading.value = true
    error.value = ''
    pagination.resetToFirstPage()
    
    const queryParams = filter.buildQueryParams()
    const { data, error: apiError } = await fetchArticlesAPI(queryParams)
    
    allArticles.value = data
    if (apiError) error.value = apiError
    loading.value = false
  }
  
  return {
    // State
    allArticles,
    categories,
    loading,
    error,
    
    // Filter
    filter,
    
    // Pagination
    pagination,
    
    // Methods
    loadCategories,
    loadArticles
  }
}

// Read Articles
export function useArticles() {
  const {
    allArticles,
    categories,
    loading,
    error,
    filter,
    pagination,
    loadCategories,
    loadArticles
  } = useArticleListState()
  
  // Initialize data on mount
  onMounted(async () => {
    await Promise.all([
      loadCategories(),
      loadArticles()
    ])
  })
  
  // Watch for filter changes
  watch([filter.searchTerm, filter.selectedCategory], loadArticles)
  
  return {
    // Data
    articles: pagination.paginatedItems,
    allArticles,
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
    fetchArticles: loadArticles,
    refreshArticles: loadArticles,
    formatDate
  }
}

// Article Form Management
export function useArticleForm() {
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
  watch(() => formData.value.title, (newTitle) => {
    formData.value.slug = slugify(newTitle)
  })
  
  // Validation rules
  const validationRules = {
    title: (value: string) => !value.trim() ? 'กรุณาใส่ชื่อบทความ' : null,
    content: (value: string) => !value.trim() ? 'กรุณาใส่เนื้อหาบทความ' : null,
    categoryName: (value: string) => !value.trim() ? 'กรุณาใส่ชื่อหมวดหมู่' : null
  }
  
  // Validate form
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}
    
    Object.entries(validationRules).forEach(([field, validator]) => {
      const error = validator(formData.value[field as keyof typeof formData.value])
      if (error) errors[field] = error
    })
    
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
  
  // Prepare payload for API
  const preparePayload = () => ({
    title: formData.value.title.trim(),
    slug: formData.value.slug.trim(),
    content: formData.value.content.trim(),
    category_name: formData.value.categoryName.trim(),
    tag_names: formData.value.tags 
      ? formData.value.tags.split(',').map(t => t.trim()).filter(Boolean) 
      : []
  })
  
  return {
    formData,
    formState,
    validateForm,
    resetForm,
    preparePayload
  }
}

// Check if user is authenticated
export function useAuthCheck() {
  const checkAuth = () => {
    const token = localStorage.getItem('token')
    return {
      isAuthenticated: !!token,
      token,
      error: !token ? 'คุณต้องเข้าสู่ระบบก่อนสร้างบทความ' : null
    }
  }
  
  return { checkAuth }
}


// Create Article
export function useCreateArticle() {
  const router = useRouter()
  const { formData, formState, validateForm, resetForm, preparePayload } = useArticleForm()
  const { checkAuth } = useAuthCheck()
  
  // Handle form submission
  const handleSubmit = async () => {
    formState.value.error = {}
    formState.value.success = false
    
    // Validate form
    if (!validateForm()) {
      return
    }
    
    // Check authentication
    const { isAuthenticated, token, error: authError } = checkAuth()
    if (!isAuthenticated) {
      formState.value.error.general = authError!
      return
    }
    
    formState.value.loading = true
    
    // Submit to API
    const payload = preparePayload()
    const { success, error } = await createArticleAPI(payload, token!)
    
    if (success) {
      formState.value.success = true
      
      // Navigate after success
      setTimeout(() => {
        resetForm()
        router.push('/articles')
      }, 1500)
    } else {
      formState.value.error = error as Record<string, string>
    }
    
    formState.value.loading = false
  }
  
  // Computed properties for v-model binding
  const computedFields = {
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
    })
  }
  
  return {
    // Form data
    formData,
    ...computedFields,
    
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