<script setup lang="ts">
import type { Category } from '~/types/article'

interface Props {
  categories: Category[]
  searchTerm: string
  selectedCategory: string | number
  updateSearchTerm: (value: string) => void
  updateSelectedCategory: (value: string | number) => void
}

const props = defineProps<Props>()

const handleSearchInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  props.updateSearchTerm(target.value)
}

const handleCategoryChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  const value = target.value === '' ? '' : Number(target.value)
  props.updateSelectedCategory(value)
}
</script>


<template>
  <div class="flex flex-col md:flex-row md:items-center md:space-x-4 space-y-4 md:space-y-0 mb-8">
    <input
      :value="searchTerm"
      @input="handleSearchInput"
      type="text"
      placeholder="🔍 ค้นหาชื่อบทความหรือแท็ก"
      class="border px-4 py-2 rounded w-full md:w-2/3"
    />
    <select 
      :value="selectedCategory" 
      @change="handleCategoryChange"
      class="border px-4 py-2 rounded w-full md:w-1/3"
    >
      <option value="">📂 ทุกหมวดหมู่</option>
      <option v-for="cat in categories" :key="cat.id" :value="cat.id">
        {{ cat.name }}
      </option>
    </select>
  </div>
</template>