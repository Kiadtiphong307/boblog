import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export function useEditArticles() {
  const route = useRoute()
  const router = useRouter()
  const slug = route.params.slug as string

  // Fields
  const title = ref('')
  const content = ref('')
  const selectedCategory = ref<number | null>(null)
  const selectedTags = ref<{ id?: number; name: string }[]>([])
  const tagInput = ref('')
  const showSuggestions = ref(false)

  // Options
  const categories = ref<{ id: number; name: string }[]>([])
  const tags = ref<{ id: number; name: string }[]>([])

  // State
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Fetch article for edit
  const fetchArticle = async () => {
    loading.value = true
    try {
      const token = localStorage.getItem('token')
      const res = await $fetch<{ data: any }>(`/api/articles/${slug}`, {  
        headers: { Authorization: `Bearer ${token}` },
      })
      const article = res.data
      title.value = article.title
      content.value = article.content
      selectedCategory.value = article.category?.id || null
      selectedTags.value = article.tags?.map((tag: any) => ({ id: tag.id, name: tag.name })) || []
    } catch (err) {
      error.value = '❌ ไม่สามารถโหลดบทความได้'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Fetch categories and tags 
  const fetchOptions = async () => {
    try {
      const token = localStorage.getItem('token')
      const resCat = await $fetch<{ data: any[] }>('/api/categories', {
        headers: { Authorization: `Bearer ${token}` },
      })
      categories.value = resCat.data

      const resTags = await $fetch<{ data: any[] }>('/api/tags', {
        headers: { Authorization: `Bearer ${token}` },
      })
      tags.value = resTags.data.map(tag => ({ id: tag.id, name: tag.name }))
    } catch (err) {
      console.error('❌ โหลด options ล้มเหลว', err)
    }
  }

  // Filter tag suggestions
  const filteredTagSuggestions = computed(() => {
    const selectedNames = selectedTags.value.map(t => t.name.toLowerCase())
    return tags.value.filter(tag =>
      !selectedNames.includes(tag.name.toLowerCase()) &&
      tag.name.toLowerCase().includes(tagInput.value.toLowerCase())
    )
  })

  const handleBlur = () => {
    setTimeout(() => {
      showSuggestions.value = false
    }, 200)
  }

  const handleTagInput = () => {
    const name = tagInput.value.trim()
    if (!name) return
    const existing = tags.value.find(t => t.name.toLowerCase() === name.toLowerCase())
    const alreadySelected = selectedTags.value.find(t => t.name.toLowerCase() === name.toLowerCase())
    if (!alreadySelected) {
      if (existing) {
        selectedTags.value.push(existing)
      } else {
        selectedTags.value.push({ name })
      }
    }
    tagInput.value = ''
    showSuggestions.value = false
  }

  const selectTag = (tag: { id: number; name: string }) => {
    selectedTags.value.push(tag)
    tagInput.value = ''
    showSuggestions.value = false
  }

  const removeTag = (tag: { id?: number; name: string }) => {
    selectedTags.value = selectedTags.value.filter(t => t !== tag)
  }

  const updateArticle = async () => {
    try {
      const token = localStorage.getItem('token')
      await $fetch(`/api/articles/${slug}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: {
          title: title.value,
          content: content.value,
          category_id: selectedCategory.value,
          tag_ids: selectedTags.value.filter(tag => tag.id).map(tag => tag.id),
          new_tags: selectedTags.value.filter(tag => !tag.id).map(tag => tag.name),
        },
      })
      alert('✅ แก้ไขบทความเรียบร้อยแล้ว')
      router.push('/articles/my-articles')
    } catch (err) {
      alert('❌ ไม่สามารถแก้ไขบทความได้')
      console.error(err)
    }
  }

  onMounted(async () => {
    await fetchOptions()
    await fetchArticle()
  })

  return {
    title,
    content,
    selectedCategory,
    selectedTags,
    tagInput,
    showSuggestions,
    categories,
    tags,
    loading,
    error,
    filteredTagSuggestions,
    fetchOptions,
    handleBlur,
    handleTagInput,
    selectTag,
    removeTag,
    updateArticle,
  }
}
