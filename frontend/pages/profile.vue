<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const user = ref({})
const form = ref({
  first_name: '',
  last_name: '',
  nickname: '',
  bio: ''
})
const success = ref('')
const error = ref('')
const confirmDelete = ref(false)
const router = useRouter()

const fetchProfile = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await $fetch('/api/user/', {
      headers: { Authorization: `Bearer ${token}` }
    })
    user.value = res.data
    form.value = { ...res.data }
  } catch (err) {
    error.value = '❌ ไม่สามารถโหลดข้อมูลผู้ใช้ได้'
  }
}

const updateProfile = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await $fetch('/api/user', {
      method: 'PUT',
      body: form.value,
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    if (res.success) {
      success.value = '✅ แก้ไขข้อมูลสำเร็จ'
    } else {
      error.value = '❌ เกิดข้อผิดพลาดในการแก้ไข'
    }
  } catch (err) {
    error.value = '❌ ไม่สามารถบันทึกข้อมูลได้'
  }
}

const deleteAccount = async () => {
  try {
    const token = localStorage.getItem('token')
    await $fetch('/api/user', {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    localStorage.removeItem('token')  // เคลียร์ token
    router.push('/')                  // redirect กลับหน้าแรก
  } catch (err) {
    error.value = '❌ ไม่สามารถลบบัญชีได้'
  }
}


onMounted(fetchProfile)
</script>

<template>
  <div class="container mx-auto p-6 max-w-xl">
    <h1 class="text-2xl font-bold mb-6">โปรไฟล์ของฉัน</h1>

    <div v-if="success" class="text-green-600 mb-4">{{ success }}</div>
    <div v-if="error" class="text-red-600 mb-4">{{ error }}</div>

    <form @submit.prevent="updateProfile" class="space-y-4">
      <div>
        <label class="block mb-1 font-medium">ชื่อ</label>
        <input v-model="form.first_name" type="text" class="w-full p-2 border rounded" />
      </div>

      <div>
        <label class="block mb-1 font-medium">นามสกุล</label>
        <input v-model="form.last_name" type="text" class="w-full p-2 border rounded" />
      </div>

      <div>
        <label class="block mb-1 font-medium">ชื่อเล่น</label>
        <input v-model="form.nickname" type="text" class="w-full p-2 border rounded" />
      </div>

      <div>
        <label class="block mb-1 font-medium">เกี่ยวกับฉัน</label>
        <textarea v-model="form.bio" class="w-full p-2 border rounded"></textarea>
      </div>

      <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded">บันทึก</button>
    </form>

    <div class="mt-6">
      <button @click="confirmDelete = true" class="text-red-500">ลบบัญชีผู้ใช้</button>
    </div>

    <!-- Modal ยืนยันการลบ -->
    <div v-if="confirmDelete" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
        <p class="mb-4">คุณแน่ใจหรือไม่ว่าต้องการลบบัญชีของคุณ?</p>
        <div class="flex justify-end gap-3">
          <button @click="confirmDelete = false" class="text-gray-600">ยกเลิก</button>
          <button @click="deleteAccount" class="text-white bg-red-600 px-4 py-2 rounded">ยืนยัน</button>
        </div>
      </div>
    </div>
  </div>
</template>


