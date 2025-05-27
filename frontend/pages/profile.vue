<script setup>
import { ref, onMounted } from "vue";

const user = ref({});
const form = ref({
  first_name: "",
  last_name: "",
  nickname: "",
  bio: "",
});
const selectedFile = ref(null);
const previewImage = ref(null);

const success = ref("");
const error = ref("");
const loading = ref(false);

// โหลดข้อมูลผู้ใช้
const fetchProfile = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("Missing token");

    const res = await $fetch("/api/user/", {
      headers: { Authorization: `Bearer ${token}` },
    });

    user.value = res.data;
    form.value = { ...res.data };
    previewImage.value = res.data.image || null;
  } catch (err) {
    console.error("❌ Error loading profile:", err);
    error.value = "❌ ไม่สามารถโหลดข้อมูลผู้ใช้ได้";
  } finally {
    loading.value = false;
  }
};

// เมื่อเลือกรูปภาพใหม่
const handleFileChange = (e) => {
  const file = e.target.files[0];
  if (!file) return;

  if (file.size > 10 * 1024 * 1024) {
    error.value = "❌ ขนาดไฟล์ต้องไม่เกิน 10MB";
    selectedFile.value = null;
    previewImage.value = user.value.image || null;
    return;
  }

  selectedFile.value = file;
  previewImage.value = URL.createObjectURL(file);
};

// กดบันทึก
const updateProfile = async () => {
  success.value = "";
  error.value = "";

  try {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("Missing token");

    const formData = new FormData();
    formData.append("first_name", form.value.first_name);
    formData.append("last_name", form.value.last_name);
    formData.append("nickname", form.value.nickname);
    formData.append("bio", form.value.bio || "");
    if (selectedFile.value) {
      formData.append("avatar", selectedFile.value);
    }

    const res = await fetch("/api/user", {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });

    const json = await res.json();
    if (res.ok) {
      success.value = "✅ แก้ไขข้อมูลสำเร็จ";
      user.value = json.data;
      form.value = { ...json.data };
      previewImage.value = json.data.image || null;
      selectedFile.value = null;
    } else {
      error.value = json.message || "❌ แก้ไขไม่สำเร็จ";
    }
  } catch (err) {
    console.error("❌ Error updating profile:", err);
    error.value = "❌ ไม่สามารถบันทึกข้อมูลได้";
  }
};

onMounted(fetchProfile);
</script>

<template>
  <div class="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-2xl mx-auto bg-white shadow-xl rounded-2xl p-8">
      <h2 class="text-3xl font-bold text-center text-blue-700 mb-8">แก้ไขโปรไฟล์</h2>

      <div v-if="success" class="bg-green-100 text-green-800 px-4 py-2 rounded mb-4 text-sm font-medium">
        {{ success }}
      </div>
      <div v-if="error" class="bg-red-100 text-red-800 px-4 py-2 rounded mb-4 text-sm font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="updateProfile" class="space-y-6">
        <!-- Avatar -->
        <div class="flex flex-col items-center">
          <div class="w-24 h-24 rounded-full border overflow-hidden shadow">
            <img :src="previewImage" v-if="previewImage" class="w-full h-full object-cover" />
            <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center text-gray-500 text-sm">No Image</div>
          </div>
          <label class="mt-4 text-sm font-medium text-gray-700">อัปโหลดรูปใหม่</label>
          <input type="file" accept="image/*" @change="handleFileChange" class="mt-1 text-sm text-gray-600" />
          <p class="text-xs text-gray-400 mt-1">ขนาดไม่เกิน 10MB</p>
        </div>

        <!-- First Name -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">ชื่อ</label>
          <input
            v-model="form.first_name"
            type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

        <!-- Last Name -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">นามสกุล</label>
          <input
            v-model="form.last_name"
            type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

        <!-- Nickname -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">ชื่อเล่น</label>
          <input
            v-model="form.nickname"
            type="text"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

        <!-- Bio -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-1">เกี่ยวกับฉัน</label>
          <textarea
            v-model="form.bio"
            rows="4"
            class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          ></textarea>
        </div>

        <!-- Submit -->
        <div class="pt-4 text-right">
          <button
            type="submit"
            class="inline-block bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-2 rounded-lg shadow-md transition"
          >
            💾 บันทึกข้อมูล
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
