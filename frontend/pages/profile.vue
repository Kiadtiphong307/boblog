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
  } catch (err) {
    console.error("❌ Error loading profile:", err);
    error.value = "❌ ไม่สามารถโหลดข้อมูลผู้ใช้ได้";
  } finally {
    loading.value = false;
  }
};

const updateProfile = async () => {
  success.value = "";
  error.value = "";

  try {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("Missing token");

    const res = await $fetch("/api/user", {
      method: "PUT",
      body: form.value,
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (res.success) {
      success.value = "✅ แก้ไขข้อมูลสำเร็จ";
    } else {
      error.value = res.message || "❌ เกิดข้อผิดพลาดในการแก้ไข";
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
