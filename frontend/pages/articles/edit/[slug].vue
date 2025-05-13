<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const slug = route.params.slug;

// บทความ
const title = ref("");
const content = ref("");
const selectedCategory = ref(null);
const selectedTags = ref([]);

// สถานะ
const loading = ref(false);
const error = ref(null);

// ตัวเลือก
const categories = ref([]);
const tags = ref([]);

// โหลดบทความ
const fetchArticle = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem("token");
    const res = await $fetch(`/api/articles/${slug}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const article = res.data;
    title.value = article.title;
    content.value = article.content;
    selectedCategory.value = article.category?.id || null;
    selectedTags.value = article.tags?.map((tag) => tag.id) || [];
  } catch (err) {
    error.value = "❌ ไม่สามารถโหลดบทความได้";
    console.error(err);
  } finally {
    loading.value = false;
  }
};

// โหลดหมวดหมู่ + แท็ก
const fetchOptions = async () => {
  try {
    const token = localStorage.getItem("token");

    const resCat = await $fetch("/api/categories", {
      headers: { Authorization: `Bearer ${token}` },
    });
    categories.value = resCat.data;

    const resTags = await $fetch("/api/tags", {
      headers: { Authorization: `Bearer ${token}` },
    });

    console.log("✅ TAG RESPONSE:", resTags);
    tags.value = (resTags.data || []).map((tag) => ({
      id: tag.ID,
      name: tag.Name,
    }));

    console.log("📦 TAGS IN COMPONENT:", tags.value);
  } catch (err) {
    console.error("❌ โหลด options ล้มเหลว", err);
  }
};

// ส่งข้อมูลอัปเดต
const updateArticle = async () => {
  try {
    const token = localStorage.getItem("token");
    await $fetch(`/api/articles/${slug}`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: {
        title: title.value,
        content: content.value,
        category_id: selectedCategory.value,
        tag_ids: selectedTags.value,
      },
    });

    alert("✅ แก้ไขบทความเรียบร้อยแล้ว");
    router.push("/articles/my-articles");
  } catch (err) {
    alert("❌ ไม่สามารถแก้ไขบทความได้");
    console.error(err);
  }
};

// เริ่มโหลดเมื่อเข้าหน้า
onMounted(async () => {
  await fetchOptions();
  await fetchArticle();
});
</script>

<template>
  <div class="container mx-auto p-4 max-w-3xl">
    <h1 class="text-2xl font-bold mb-6">✏️ แก้ไขบทความ</h1>

    <div v-if="loading">⏳ กำลังโหลด...</div>
    <div v-if="error" class="text-red-500">{{ error }}</div>

    <form @submit.prevent="updateArticle" v-if="!loading">
      <div class="mb-4">
        <label class="block mb-1 font-semibold">หัวข้อบทความ</label>
        <input
          v-model="title"
          type="text"
          class="w-full border p-2 rounded"
          required
        />
      </div>

      <div class="mb-4">
        <label class="block mb-1 font-semibold">เนื้อหา</label>
        <textarea
          v-model="content"
          class="w-full border p-2 rounded"
          rows="8"
          required
        ></textarea>
      </div>

      <div class="mb-4">
        <label class="block mb-1 font-semibold">หมวดหมู่</label>
        <select
          v-model="selectedCategory"
          class="w-full border p-2 rounded"
          required
        >
          <option value="" disabled>-- เลือกหมวดหมู่ --</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>

      <div class="mb-4">
        <label class="block mb-1 font-semibold">แท็ก</label>

        <div v-if="tags.length === 0" class="text-gray-500">
          🔄 กำลังโหลดแท็ก... หรือไม่มีแท็กให้เลือก
        </div>

        <div class="flex flex-wrap gap-2" v-if="tags.length > 0">
          <label
            v-for="tag in tags"
            :key="tag.id"
            class="flex items-center gap-1"
          >
            <input type="checkbox" :value="tag.id" v-model="selectedTags" />
            {{ tag.name }}
          </label>
        </div>
      </div>

      <button
        type="submit"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
      >
        💾 บันทึกการแก้ไข
      </button>
    </form>
  </div>
</template>
