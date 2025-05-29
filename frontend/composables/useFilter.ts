import { ref, watch, onUnmounted } from 'vue'

export interface FilterOptions {
  debounceDelay?: number
}

export function useFilter(options: FilterOptions = {}) {
  const { debounceDelay = 500 } = options
  
  // Filter states
  const selectedCategory = ref<string | number>('')
  const searchTerm = ref('')
  const localSearchTerm = ref('') // for immediate v-model binding
  
  // Debounce timer
  let debounceTimer: NodeJS.Timeout | null = null
  
  // Debounced search function
  const debouncedSearch = (value: string) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    
    debounceTimer = setTimeout(() => {
      searchTerm.value = value.trim()
    }, debounceDelay)
  }
  
  // Update methods
  const updateSearchTerm = (value: string) => {
    localSearchTerm.value = value
    debouncedSearch(value)
  }
  
  const updateSelectedCategory = (value: string | number) => {
    selectedCategory.value = value
  }
  
  const clearSearch = () => {
    localSearchTerm.value = ''
    searchTerm.value = ''
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
  }
  
  const clearCategory = () => {
    selectedCategory.value = ''
  }
  
  const clearAllFilters = () => {
    clearSearch()
    clearCategory()
  }
  
  // Build query parameters for API calls
  const buildQueryParams = () => {
    const params = new URLSearchParams()
    
    if (searchTerm.value.trim()) {
      params.append('search', searchTerm.value.trim())
    }
    
    if (selectedCategory.value) {
      params.append('category_id', selectedCategory.value.toString())
    }
    
    return params
  }
  
  // Check if any filters are active
  const hasActiveFilters = computed(() => {
    return searchTerm.value.trim() !== '' || selectedCategory.value !== ''
  })
  
  // Get filter summary for display
  const getFilterSummary = () => {
    const filters: string[] = []
    
    if (searchTerm.value.trim()) {
      filters.push(`ค้นหา: "${searchTerm.value.trim()}"`)
    }
    
    if (selectedCategory.value) {
      filters.push(`หมวดหมู่: ${selectedCategory.value}`)
    }
    
    return filters
  }
  
  // Cleanup on unmount
  onUnmounted(() => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
  })
  
  return {
    // State
    selectedCategory,
    searchTerm,
    localSearchTerm,
    
    // Computed
    hasActiveFilters,
    
    // Methods
    updateSearchTerm,
    updateSelectedCategory,
    clearSearch,
    clearCategory,
    clearAllFilters,
    buildQueryParams,
    getFilterSummary
  }
}