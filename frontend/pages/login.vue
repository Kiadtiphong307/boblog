<script setup>
definePageMeta({
  layout: false
})

import { ref } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()

const login = async () => {
  error.value = ''

  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, password: password.value }),
    })

    const json = await res.json()

    if (res.ok && json.data.token) {
      localStorage.setItem('token', json.data.token)
      router.push('/') // กลับไปหน้าแรกหลัง login สำเร็จ
    } else {
      error.value = json.error || 'เข้าสู่ระบบล้มเหลว'
    }
  } catch (err) {
    error.value = 'ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded shadow-md w-full max-w-md">
      <h1 class="text-2xl font-bold mb-6 text-center">🔐 เข้าสู่ระบบ</h1>

      <form @submit.prevent="login" class="space-y-4">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          required
          class="input input-bordered w-full"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          required
          class="input input-bordered w-full"
        />
        <button
          type="submit"
          class="btn btn-primary w-full"
        >
          เข้าสู่ระบบ
        </button>
      </form>

      <p v-if="error" class="text-red-600 mt-3 text-center">{{ error }}</p>

      <p class="mt-6 text-center text-sm text-gray-600">
        ยังไม่มีบัญชี?
        <NuxtLink to="/register" class="text-blue-600 hover:underline">
          สมัครสมาชิก
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
