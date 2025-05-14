<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useRuntimeConfig } from "#app";

const config = useRuntimeConfig();
const router = useRouter();

const title = ref("");
const slug = ref(""); // generate จาก title
const content = ref("");
const categoryName = ref(""); // 🟢 ใช้ชื่อหมวดหมู่แทน ID
const tags = ref("");
const error = ref("");
const success = ref(false);

// ฟังก์ชันสร้าง slug
const slugify = (text: string): string =>
  text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-ก-๙]/g, "")
    .replace(/[\s_-]+/g, "-")
    .replace(/^-+|-+$/g, "");

// generate slug อัตโนมัติ
watch(title, (newTitle) => {
  slug.value = slugify(newTitle);
});

// ฟอร์ม submit
const handleSubmit = async () => {
  error.value = "";
  success.value = false;

  const token = localStorage.getItem("token");
  if (!token) {
    error.value = "คุณต้องเข้าสู่ระบบก่อนสร้างบทความ";
    return;
  }

  try {
    await $fetch("/api/articles", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: {
        title: title.value,
        slug: slug.value,
        content: content.value,
        category_name: categoryName.value, // ✅ ส่งชื่อหมวดหมู่แทน ID
        tag_names: tags.value
          ? tags.value.split(",").map((t) => t.trim()).filter(Boolean)
          : [],
      },
    });

    success.value = true;
    setTimeout(() => router.push("/articles"), 1200);
  } catch (e: any) {
    error.value = e?.data?.message || "❌ เกิดข้อผิดพลาดในการสร้างบทความ";
  }
};
</script>

<template>
  <div class="max-w-xl mx-auto py-8">
    <h1 class="text-2xl font-bold mb-4">📝 สร้างบทความ</h1>

    <form @submit.prevent="handleSubmit" class="space-y-4">
      <input
        v-model="title"
        type="text"
        placeholder="ชื่อบทความ"
        class="w-full p-2 border rounded"
        required
      />
      <p class="text-sm text-gray-500">🔗 Slug ที่สร้าง: {{ slug }}</p>

      <textarea
        v-model="content"
        placeholder="เนื้อหา"
        class="w-full p-2 border rounded"
        rows="6"
        required
      />

      <input
        v-model="categoryName"
        type="text"
        placeholder="หมวดหมู่ (เช่น ข่าว, บทความ)"
        class="w-full p-2 border rounded"
        required
      />

      <input
        v-model="tags"
        type="text"
        placeholder="แท็ก (คั่นด้วย , เช่น go, fiber)"
        class="w-full p-2 border rounded"
      />

      <button
        type="submit"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
      >
        ➕ สร้างบทความ
      </button>
    </form>

    <div v-if="error" class="text-red-500 mt-2">⚠️ {{ error }}</div>
    <div v-if="success" class="text-green-500 mt-2">
      ✅ สร้างบทความเรียบร้อยแล้ว!
    </div>
  </div>
</template>
