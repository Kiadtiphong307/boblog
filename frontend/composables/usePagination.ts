import { ref, computed } from 'vue'

export interface PaginationOptions {
  perPage?: number
  initialPage?: number
}

export function usePagination<T>(
  items: Ref<T[]>,
  options: PaginationOptions = {}
) {
  const { perPage = 4, initialPage = 1 } = options
  
  const currentPage = ref(initialPage)
  const itemsPerPage = ref(perPage)
  
  // Computed values
  const totalItems = computed(() => items.value.length)
  const totalPages = computed(() => 
    Math.ceil(totalItems.value / itemsPerPage.value) || 1
  )
  
  const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value
    const end = start + itemsPerPage.value
    return items.value.slice(start, end)
  })
  
  const hasNextPage = computed(() => currentPage.value < totalPages.value)
  const hasPrevPage = computed(() => currentPage.value > 1)
  
  const pageInfo = computed(() => ({
    current: currentPage.value,
    total: totalPages.value,
    showing: {
      from: totalItems.value === 0 ? 0 : (currentPage.value - 1) * itemsPerPage.value + 1,
      to: Math.min(currentPage.value * itemsPerPage.value, totalItems.value),
      total: totalItems.value
    }
  }))
  
  // Methods
  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page
    }
  }
  
  const nextPage = () => {
    if (hasNextPage.value) {
      currentPage.value++
    }
  }
  
  const prevPage = () => {
    if (hasPrevPage.value) {
      currentPage.value--
    }
  }
  
  const resetToFirstPage = () => {
    currentPage.value = 1
  }
  
  const setItemsPerPage = (newPerPage: number) => {
    itemsPerPage.value = newPerPage
    resetToFirstPage()
  }
  
  // Get page numbers for pagination UI
  const getPageNumbers = (maxVisible = 5) => {
    const pages: number[] = []
    const total = totalPages.value
    const current = currentPage.value
    
    if (total <= maxVisible) {
      for (let i = 1; i <= total; i++) {
        pages.push(i)
      }
    } else {
      const start = Math.max(1, current - Math.floor(maxVisible / 2))
      const end = Math.min(total, start + maxVisible - 1)
      
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
    }
    
    return pages
  }
  
  return {
    // State
    currentPage,
    itemsPerPage,
    
    // Computed
    totalItems,
    totalPages,
    paginatedItems,
    hasNextPage,
    hasPrevPage,
    pageInfo,
    
    // Methods
    goToPage,
    nextPage,
    prevPage,
    resetToFirstPage,
    setItemsPerPage,
    getPageNumbers
  }
}