<script setup lang="ts">
definePageMeta({ layout: false });

import { ref } from "vue";

const form = ref({
  username: "",
  email: "",
  password: "",
  confirm_password: "",
  first_name: "",
  last_name: "",
  nickname: "",
});

const error = ref<Record<string, string>>({});
const success = ref("");

const register = async () => {
  error.value = {};
  success.value = "";

  if (form.value.password !== form.value.confirm_password) {
    error.value.confirm_password = "❌ รหัสผ่านไม่ตรงกัน";
    return;
  }

  try {
    const res = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(form.value),
    });

    const json = await res.json();

    if (res.ok) {
      success.value = "✅ สมัครสมาชิกสำเร็จ";
      form.value = {
        username: "",
        email: "",
        password: "",
        confirm_password: "",
        first_name: "",
        last_name: "",
        nickname: "",
      };
    } else {
      // รับ errors แบบ object เช่น { email: "อีเมลนี้มีผู้ใช้งานแล้ว" }
      if (json.errors && typeof json.errors === "object") {
        error.value = json.errors;
      } else {
        error.value.general = json.error || "❌ เกิดข้อผิดพลาดในการสมัคร";
      }
    }
  } catch (err) {
    error.value.general = "❌ ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์ได้";
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-xl p-8 space-y-6">
      <h1 class="text-2xl font-bold text-center text-gray-800">
        📝 สมัครสมาชิก
      </h1>

      <form @submit.prevent="register" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.username" type="text" placeholder="ชื่อผู้ใช้งาน"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.username" class="text-sm text-red-500 mt-1">
              {{ error.username }}
            </p>
          </div>
          <div>
            <input v-model="form.nickname" type="text" placeholder="ชื่อเล่น"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.nickname" class="text-sm text-red-500 mt-1">
              {{ error.nickname }}
            </p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.first_name" type="text" placeholder="ชื่อ"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.first_name" class="text-sm text-red-500 mt-1">
              {{ error.first_name }}
            </p>
          </div>
          <div>
            <input v-model="form.last_name" type="text" placeholder="นามสกุล"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.last_name" class="text-sm text-red-500 mt-1">
              {{ error.last_name }}
            </p>
          </div>
        </div>

        <div>
          <input v-model="form.email" type="email" placeholder="อีเมล"
            class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
          <p v-if="error.email" class="text-sm text-red-500 mt-1">
            {{ error.email }}
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <input v-model="form.password" type="password" placeholder="รหัสผ่าน"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.password" class="text-sm text-red-500 mt-1">
              {{ error.password }}
            </p>
          </div>
          <div>
            <input v-model="form.confirm_password" type="password" placeholder="ยืนยันรหัสผ่าน"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200" />
            <p v-if="error.confirm_password" class="text-sm text-red-500 mt-1">
              {{ error.confirm_password }}
            </p>
          </div>
        </div>

        <button type="submit"
          class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition">
          ➕ สมัครสมาชิก
        </button>
      </form>

      <p v-if="success" class="text-green-600 text-center font-medium">
        {{ success }}
      </p>
      <p v-if="error.general" class="text-red-600 text-center font-medium">
        {{ error.general }}
      </p>
    </div>
  </div>
</template>
