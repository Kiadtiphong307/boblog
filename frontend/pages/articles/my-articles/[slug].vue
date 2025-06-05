<template>
  <div class="max-w-3xl mx-auto py-8 px-4">
    <!-- Title -->
    <h1 class="text-3xl font-bold mb-4">{{ article?.title }}</h1>
    
    <!-- Author, Category, Date -->
    <div class="text-sm text-gray-500 mb-2">
      👤 {{ article?.author?.username }} |
      📂 {{ article?.category?.name }} |
      🕒 {{ formatDate(article?.created_at || '') }}
    </div>
    
    <div class="mb-4">
      <span
        v-for="tag in article?.tags || []"
        :key="tag.id"
        class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded mr-2"
      >
        #{{ tag.name }}
      </span>
    </div>
    
    <div class="prose max-w-none" v-html="article?.content"></div>

    <!-- Comments -->
    <hr class="my-6" />
    <h2 class="text-lg font-semibold mb-2">💬 Comments ({{ comments.length }})</h2>
    
    <ul class="space-y-3 mb-6">
      <li
        v-for="c in comments"
        :key="c.id"
        class="bg-gray-100 p-3 rounded text-sm text-gray-800"
      >
        <!-- Comment Header with Action Buttons -->
        <div class="flex justify-between items-start mb-1">
          <div>
            <span class="font-semibold text-blue-700">
              👤 {{ c.user?.username }}
            </span>
            <span class="text-gray-500 text-xs ml-2">
              🕒 {{ formatDate(c.created_at) }}
            </span>
            <span v-if="c.updated_at && c.updated_at !== c.created_at" 
                  class="text-xs text-gray-400 italic ml-1">
              (แก้ไขแล้ว)
            </span>
          </div>
          
          <!-- Action Buttons (แสดงเฉพาะเจ้าของคอมเมนต์) -->
          <div v-if="isCommentOwner(c)" class="flex space-x-1">
            <button
              @click="startEditComment(c)"
              class="text-blue-600 hover:text-blue-800 text-xs font-medium px-1"
            >
              แก้ไข
            </button>
            <span class="text-gray-400">|</span>
            <button
              @click="deleteComment(c.id)"
              class="text-red-600 hover:text-red-800 text-xs font-medium px-1"
            >
              ลบ
            </button>
          </div>
        </div>

        <!-- Comment Content -->
        <div v-if="editingCommentId !== c.id">
          <div>{{ c.content }}</div>
        </div>

        <!-- Edit Form -->
        <div v-else class="space-y-2 mt-2">
          <textarea 
            v-model="editContent"
            rows="3"
            class="w-full px-2 py-1 border border-gray-300 rounded text-sm resize-none"
          ></textarea>
          <div class="flex space-x-2">
            <button
              @click="saveEditComment(c.id)"
              :disabled="!editContent.trim()"
              class="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white text-xs font-medium py-1 px-2 rounded"
            >
              บันทึก
            </button>
            <button
              @click="cancelEdit"
              class="bg-gray-600 hover:bg-gray-700 text-white text-xs font-medium py-1 px-2 rounded"
            >
              ยกเลิก
            </button>
          </div>
        </div>
      </li>
    </ul>

    <!-- Comment Form -->
    <form @submit.prevent="submitComment" class="space-y-2">
      <textarea
        v-model="newComment"
        rows="3"
        class="w-full border rounded p-2"
        placeholder="แสดงความคิดเห็น..."
      ></textarea>
      <button
        type="submit"
        :disabled="!newComment.trim()"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
      >
        💬 แสดงความคิดเห็น
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useComment } from '~/composables/useComment'

const {
  article,
  comments,
  newComment,
  editingCommentId,
  editContent,
  submitComment,
  startEditComment,
  cancelEdit,
  saveEditComment,
  deleteComment,
  isCommentOwner,
  formatDate,
} = useComment()
</script>