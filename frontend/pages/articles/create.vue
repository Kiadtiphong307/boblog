<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useRuntimeConfig } from '#app'

const config = useRuntimeConfig()
const router = useRouter()

const title = ref('')
const slug = ref('')
const content = ref('')
const categoryId = ref('')
const tags = ref('')
const error = ref('')
const success = ref(false)

const handleSubmit = async () => {
  error.value = ''
  success.value = false

  const token = localStorage.getItem('token')
  if (!token) {
    error.value = 'คุณต้องเข้าสู่ระบบก่อนสร้างบทความ'
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
        category_id: Number(categoryId.value),
        tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      },
    })

    success.value = true
    setTimeout(() => router.push('/articles'), 1200)
  } catch (e: any) {
    error.value = e?.data?.message || '❌ เกิดข้อผิดพลาดในการสร้างบทความ'
  }
}
</script>


<template>
    <div class="max-w-xl mx-auto py-8">
      <h1 class="text-2xl font-bold mb-4">📝 Create Article</h1>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <input v-model="title" type="text" placeholder="Title" class="w-full p-2 border rounded" required />
        <input v-model="slug" type="text" placeholder="Slug (unique)" class="w-full p-2 border rounded" required />
        <textarea v-model="content" placeholder="Content" class="w-full p-2 border rounded" rows="6" required />
        <input v-model="categoryId" type="number" placeholder="Category ID" class="w-full p-2 border rounded" required />
        <input v-model="tags" type="text" placeholder="Tags (comma separated)" class="w-full p-2 border rounded" />
        <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition">
          ➕ Create
        </button>
      </form>
  
      <div v-if="error" class="text-red-500 mt-2">⚠️ {{ error }}</div>
      <div v-if="success" class="text-green-500 mt-2">✅ Article created successfully!</div>
    </div>
  </template>
  

  