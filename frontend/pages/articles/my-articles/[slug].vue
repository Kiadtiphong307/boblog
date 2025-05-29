<template>
  <div class="max-w-4xl mx-auto py-10 px-6">
    <!-- Article Content -->
    <div v-if="article" class="mb-8">
      <h1 class="text-3xl font-bold text-gray-800 mb-4">{{ article.title }}</h1>
      <div class="prose max-w-none">
        <div v-html="article.content"></div>
      </div>
    </div>

    <!-- Comments Section -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-bold text-gray-800 mb-6 flex items-center">
        💬 ความคิดเห็น ({{ comments.length }})
      </h2>

      <!-- Add New Comment Form -->
      <div class="mb-8 p-4 bg-gray-50 rounded-lg">
        <h3 class="font-semibold text-gray-700 mb-3">✍️ แสดงความคิดเห็น</h3>
        <textarea 
          v-model="newComment"
          rows="4"
          placeholder="แสดงความคิดเห็นของคุณ..."
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400 focus:border-transparent resize-none"
        ></textarea>
        <div class="flex justify-end mt-3">
          <button
            @click="submitComment"
            :disabled="!newComment.trim()"
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-lg transition"
          >
            📝 ส่งความคิดเห็น
          </button>
        </div>
      </div>

      <!-- Comments List -->
      <div v-if="comments.length > 0" class="space-y-4">
        <div 
          v-for="comment in comments" 
          :key="comment.id"
          class="border-l-4 border-blue-200 pl-4 py-4 bg-gray-50 rounded-r-lg"
        >
          <!-- Comment Header -->
          <div class="flex justify-between items-start mb-2">
            <div class="flex items-center space-x-2">
              <span class="font-semibold text-gray-800">
                {{ comment.user?.username }}
              </span>
              <span class="text-sm text-gray-500">
                {{ formatDate(comment.created_at) }}
              </span>
              <span v-if="comment.updated_at && comment.updated_at !== comment.created_at" 
                    class="text-xs text-gray-400 italic">
                (แก้ไขแล้ว)
              </span>
            </div>
            
            <!-- Action Buttons (แสดงเฉพาะเจ้าของคอมเมนต์) -->
            <div v-if="isCommentOwner(comment)" class="flex space-x-2">
              <button
                @click="startEditComment(comment)"
                class="text-blue-600 hover:text-blue-800 text-sm font-medium"
              >
                ✏️ แก้ไข
              </button>
              <button
                @click="deleteComment(comment.id)"
                class="text-red-600 hover:text-red-800 text-sm font-medium"
              >
                🗑️ ลบ
              </button>
            </div>
          </div>

          <!-- Comment Content -->
          <div v-if="editingCommentId !== comment.id">
            <p class="text-gray-700 leading-relaxed">{{ comment.content }}</p>
          </div>

          <!-- Edit Form -->
          <div v-else class="space-y-3">
            <textarea 
              v-model="editContent"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400 focus:border-transparent resize-none"
            ></textarea>
            <div class="flex space-x-2">
              <button
                @click="saveEditComment(comment.id)"
                :disabled="!editContent.trim()"
                class="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white text-sm font-medium py-1 px-3 rounded transition"
              >
                💾 บันทึก
              </button>
              <button
                @click="cancelEdit"
                class="bg-gray-600 hover:bg-gray-700 text-white text-sm font-medium py-1 px-3 rounded transition"
              >
                ❌ ยกเลิก
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- No Comments Message -->
      <div v-else class="text-center py-8 text-gray-500">
        <p class="text-lg">📝 ยังไม่มีความคิดเห็น</p>
        <p class="text-sm">เป็นคนแรกที่แสดงความคิดเห็นในบทความนี้!</p>
      </div>
    </div>
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
  fetchArticle,
  fetchComments,
  submitComment,
  startEditComment,
  cancelEdit,
  saveEditComment,
  deleteComment,
  isCommentOwner,
  formatDate,
} = useComment()
</script>