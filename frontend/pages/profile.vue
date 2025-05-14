<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
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

const handleFileChange = (e) => {
  const file = e.target.files[0];
  if (!file) return;

  // Optional: Check file size (max 10MB)
  if (file.size > 10 * 1024 * 1024) {
    error.value = "❌ ขนาดไฟล์ต้องไม่เกิน 10MB";
    selectedFile.value = null;
    previewImage.value = user.value.image || null;
    return;
  }

  selectedFile.value = file;
  previewImage.value = URL.createObjectURL(file); // show preview
};

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
  <div class="container mx-auto p-6 max-w-xl">
    <h1 class="text-2xl font-bold mb-6">โปรไฟล์ของฉัน</h1>

    <div v-if="success" class="text-green-600 mb-4">{{ success }}</div>
    <div v-if="error" class="text-red-600 mb-4">{{ error }}</div>

    <form @submit.prevent="updateProfile" class="space-y-4">
      <!-- รูปโปรไฟล์ -->
      <div v-if="previewImage" class="mb-4">
        <img :src="previewImage" alt="avatar" class="w-24 h-24 object-cover rounded-full border" />
      </div>

      <div>
        <label class="block mb-1 font-medium">เปลี่ยนรูปโปรไฟล์</label>
        <input type="file" accept="image/*" @change="handleFileChange" />
        <p class="text-xs text-gray-500 mt-1">ขนาดไม่เกิน 10MB</p>
      </div>

      <div>
        <label class="block mb-1 font-medium">ชื่อ</label>
        <input
          v-model="form.first_name"
          type="text"
          class="w-full p-2 border rounded"
        />
      </div>

      <div>
        <label class="block mb-1 font-medium">นามสกุล</label>
        <input
          v-model="form.last_name"
          type="text"
          class="w-full p-2 border rounded"
        />
      </div>

      <div>
        <label class="block mb-1 font-medium">ชื่อเล่น</label>
        <input
          v-model="form.nickname"
          type="text"
          class="w-full p-2 border rounded"
        />
      </div>

      <div>
        <label class="block mb-1 font-medium">เกี่ยวกับฉัน</label>
        <textarea
          v-model="form.bio"
          class="w-full p-2 border rounded"
        ></textarea>
      </div>

      <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded">
        บันทึก
      </button>
    </form>
  </div>
</template>
