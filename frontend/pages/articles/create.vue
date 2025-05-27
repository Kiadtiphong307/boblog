 <template>
  <div class="max-w-2xl mx-auto py-12 px-6">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">📝 สร้างบทความใหม่</h1>

    <form @submit.prevent="handleSubmit" class="space-y-6 bg-white p-8 rounded-2xl shadow-xl">
      <!-- ชื่อบทความ -->
      <div>
        <label class="block text-gray-700 font-medium mb-1">ชื่อบทความ</label>
        <input
          v-model="title"
          type="text"
          placeholder="ชื่อบทความ"
          class="w-full border border-gray-300 rounded-xl p-3 focus:ring focus:ring-blue-200"
          required
        />
        <p class="text-sm text-gray-500 mt-1">🔗 Slug ที่สร้าง: <span class="font-mono">{{ slug }}</span></p>
        <p v-if="error.title" class="text-sm text-red-500 mt-1">{{ error.title }}</p>
        <p v-if="error.slug" class="text-sm text-red-500 mt-1">{{ error.slug }}</p>
      </div>

      <!-- เนื้อหา -->
      <div>
        <label class="block text-gray-700 font-medium mb-1">เนื้อหา</label>
        <textarea
          v-model="content"
          placeholder="พิมพ์เนื้อหา..."
          rows="8"
          class="w-full border border-gray-300 rounded-xl p-3 focus:ring focus:ring-blue-200"
          required
        ></textarea>
        <p v-if="error.content" class="text-sm text-red-500 mt-1">{{ error.content }}</p>
      </div>

      <!-- หมวดหมู่ -->
      <div>
        <label class="block text-gray-700 font-medium mb-1">หมวดหมู่</label>
        <input
          v-model="categoryName"
          type="text"
          placeholder="เช่น ข่าว, บทความ"
          class="w-full border border-gray-300 rounded-xl p-3 focus:ring focus:ring-blue-200"
          required
        />
        <p v-if="error.category_name" class="text-sm text-red-500 mt-1">{{ error.category_name }}</p>
      </div>

      <!-- แท็ก -->
      <div>
        <label class="block text-gray-700 font-medium mb-1">แท็ก</label>
        <input
          v-model="tags"
          type="text"
          placeholder="คั่นด้วย , เช่น go, fiber"
          class="w-full border border-gray-300 rounded-xl p-3 focus:ring focus:ring-blue-200"
        />
        <p v-if="error.tag_names" class="text-sm text-red-500 mt-1">{{ error.tag_names }}</p>
      </div>

      <!-- ปุ่ม Submit -->
      <div>
        <button
          type="submit"
          class="w-full bg-blue-600 text-white py-3 rounded-xl font-semibold hover:bg-blue-700 transition"
        >
          ➕ สร้างบทความ
        </button>
      </div>

      <!-- ข้อความแจ้งเตือน -->
      <p v-if="error.general" class="text-red-600 font-medium text-center">
        {{ error.general }}
      </p>
      <p v-if="success" class="text-green-600 font-medium text-center">
        ✅ สร้างบทความเรียบร้อยแล้ว!
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useProductForm } from '~/composables/articles/useArticles'

const {
  title,
  slug,
  content,
  categoryName,
  tags,
  error,
  success,
  handleSubmit,
} = useProductForm()
</script>
