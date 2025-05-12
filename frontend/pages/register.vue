<script setup lang="ts">
definePageMeta({
  layout: false
})

import { ref } from 'vue'

const form = ref({
    username: '',
    email: '',
    password: '',
    confirm_password: '',
})

const error = ref('')
const success = ref('')

const register = async () => {
    error.value = ''
    success.value = ''

    // ✅ ตรวจสอบรหัสผ่านตรงกันก่อนส่ง
    if (form.value.password !== form.value.confirm_password) {
        error.value = '❌ รหัสผ่านไม่ตรงกัน'
        console.error('❌ รหัสผ่านไม่ตรงกัน')
        return
    }

    try {
        const res = await fetch('/api/auth/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                username: form.value.username,
                email: form.value.email,
                password: form.value.password,
                confirm_password: form.value.confirm_password,
            }),
        })

        const json = await res.json()

        if (res.ok) {
            success.value = '✅ สมัครสมาชิกสำเร็จ'
            console.log('✅ สมัครสมาชิกสำเร็จ:', json)
            // ✅ ล้างฟอร์ม
            form.value = { username: '', email: '', password: '', confirm_password: '' }
        } else {
            error.value = json.error || '❌ เกิดข้อผิดพลาดในการสมัคร'
            console.warn('⚠️ สมัครไม่สำเร็จ:', json)
        }
    } catch (err) {
        error.value = '❌ ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์ได้'
        console.error('❌ ไม่สามารถเชื่อมต่อ API:', err)
    }
}
</script>

<template>
    <div class="max-w-md mx-auto mt-10 p-6 border rounded shadow space-y-4">
        <h1 class="text-xl font-bold text-center">📝 สมัครสมาชิก</h1>

        <form @submit.prevent="register" class="space-y-4">
            <input v-model="form.username" type="text" placeholder="Username" class="input input-bordered w-full"
                required />
            <input v-model="form.email" type="email" placeholder="Email" class="input input-bordered w-full" required />
            <input v-model="form.password" type="password" placeholder="Password" class="input input-bordered w-full"
                required />
            <input v-model="form.confirm_password" type="password" placeholder="Confirm Password"
                class="input input-bordered w-full" required />
            <button type="submit" class="btn btn-primary w-full">สมัครสมาชิก</button>
        </form>

        <p v-if="success" class="text-green-600 text-center">{{ success }}</p>
        <p v-if="error" class="text-red-600 text-center">{{ error }}</p>
    </div>
</template>
