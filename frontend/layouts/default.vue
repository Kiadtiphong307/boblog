<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const token = ref<string | null>(null)
const username = ref<string | null>(null)
const showDropdown = ref(false)
const router = useRouter()

const fetchUser = async () => {
  try {
    const res = await fetch('/api/user', {
      headers: {
        Authorization: `Bearer ${token.value}`
      }
    })
    const json = await res.json()
    console.log("USER API RESPONSE:", json)
    if (res.ok) {
      username.value = json.data.username 
    }
  } catch (err) {
    console.error('ไม่สามารถโหลดข้อมูลผู้ใช้', err)
  }
}


const logout = () => {
  localStorage.removeItem('token')
  token.value = null
  username.value = null
  router.push('/')
}

onMounted(() => {
  token.value = localStorage.getItem('token')
  if (token.value) {
    fetchUser()
  }
})
</script>

<template>
  <div>
    <nav class="flex items-center justify-between px-24 py-6 bg-white shadow-md border-b border-blue-200">
      <div class="flex items-center space-x-2">
        <span class="text-blue-500 font-bold text-4xl">BOBLOG</span>
        <span class="flex items-center ml-8 text-xl text-black cursor-pointer">
          <svg class="w-5 h-5 mr-1 text-black" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 20h9" />
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M16.5 3.5a2.121 2.121 0 113 3L7 19.5 3 21l1.5-4L16.5 3.5z" />
          </svg>
          <div class="text-black text-2xl">เขียนบทความ</div>
        </span>
      </div>

      <!-- ✅ ถ้ายังไม่ login -->
      <NuxtLink v-if="!token" to="/login">
        <button class="bg-blue-500 hover:bg-blue-600 text-white font-bold text-xl py-2 px-6 rounded transition">
          เข้าสู่ระบบ
        </button>
      </NuxtLink>

      <!-- ✅ ถ้า login แล้ว -->
      <div v-else class="relative flex items-center space-x-2">
        <!-- ไอคอนโปรไฟล์ -->
        <div @click="showDropdown = !showDropdown"
          class="bg-gray-100 hover:bg-gray-200 text-purple-800 p-3 rounded cursor-pointer transition">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"
            stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M5.121 17.804A13.937 13.937 0 0112 15c2.219 0 4.29.538 6.121 1.487M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>

        <!-- ✅ ข้อความแสดงชื่อผู้ใช้ -->
        <div class="text-gray-800 font-medium">
          สวัสดี, {{ Username }}
        </div>

        <!-- dropdown -->
        <div v-show="showDropdown" class="absolute right-0 mt-12 w-40 bg-white border rounded shadow-md z-50">
          <NuxtLink to="/profile" class="block px-4 py-2 text-gray-700 hover:bg-gray-100">โปรไฟล์</NuxtLink>
          <button @click="logout"
            class="block w-full text-left px-4 py-2 text-red-600 hover:bg-red-100">ออกจากระบบ</button>
        </div>
      </div>
    </nav>

    <div>
      <NuxtPage />
    </div>
  </div>
</template>
