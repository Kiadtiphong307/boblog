<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const user = ref({})
const form = ref({ first_name: '', last_name: '', nickname: '', bio: '' })
const router = useRouter()
const success = ref('')
const error = ref('')
const confirmDelete = ref(false)

const fetchProfile = async () => {
  const token = localStorage.getItem('token')
  const res = await fetch('/api/user', { headers: { Authorization: `Bearer ${token}` } })
  const json = await res.json()
  if (json.success) {
    user.value = json.data
    form.value = { ...json.data }
  }
}

const updateProfile = async () => {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/users/${user.value.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(form.value),
  })
  const json = await res.json()
  if (json.success) success.value = '✅ แก้ไขข้อมูลสำเร็จ'
  else error.value = json.error
}

const deleteAccount = async () => {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/users/${user.value.id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  const json = await res.json()
  if (json.success) {
    localStorage.removeItem('token')
    router.push('/login')
  } else {
    error.value = json.error
  }
}

onMounted(fetchProfile)
</script>

<template>
  <div class="max-w-lg mx-auto py-8">
    <h1 class="text-2xl font-bold mb-4">👤 โปรไฟล์ของคุณ</h1>

    <form @submit.prevent="updateProfile" class="space-y-4">
      <input v-model="form.first_name" placeholder="ชื่อจริง" class="input w-full" />
      <input v-model="form.last_name" placeholder="นามสกุล" class="input w-full" />
      <input v-model="form.nickname" placeholder="ชื่อเล่น" class="input w-full" />
      <textarea v-model="form.bio" placeholder="เกี่ยวกับคุณ..." class="textarea w-full"></textarea>

      <button type="submit" class="btn btn-primary w-full">💾 บันทึก</button>
    </form>

    <div class="text-green-600 mt-2" v-if="success">{{ success }}</div>
    <div class="text-red-600 mt-2" v-if="error">{{ error }}</div>

    <hr class="my-6" />

    <div class="text-center">
      <button v-if="!confirmDelete" @click="confirmDelete = true" class="btn btn-error">🗑️ ลบบัญชี</button>
      <div v-else>
        <p class="text-sm text-gray-600 mb-2">คุณแน่ใจหรือไม่ที่จะลบบัญชีนี้? การลบไม่สามารถย้อนคืนได้</p>
        <button @click="deleteAccount" class="btn btn-error mr-2">ยืนยันลบ</button>
        <button @click="confirmDelete = false" class="btn">ยกเลิก</button>
      </div>
    </div>
  </div>
</template>
