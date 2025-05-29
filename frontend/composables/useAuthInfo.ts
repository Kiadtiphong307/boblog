import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const token = ref<string | null>(null)
const nickname = ref<string | null>(null)
const imageUrl = ref<string | null>(null)
const showDropdown = ref(false)

export function useAuthInfo() {
  const router = useRouter()

  const fetchUser = async () => {
    if (!token.value) return

    try {
      const res = await fetch('/api/user', {
        headers: { Authorization: `Bearer ${token.value}` },
      })

      if (res.status === 401) {
        console.warn('🔒 Token หมดอายุ หรือไม่มีสิทธิ์')
        logout()
        return
      }

      const json = await res.json()
      if (res.ok) {
        nickname.value = json.data.nickname
        imageUrl.value = json.data.image || null
        
        // ✅ อัพเดตข้อมูล user ใน localStorage ด้วย
        const userData = {
          id: json.data.id,
          username: json.data.username,
          nickname: json.data.nickname,
          email: json.data.email,
          image: json.data.image
        }
        localStorage.setItem('user', JSON.stringify(userData))
        console.log('💾 User data updated:', userData)
      }
    } catch (err) {
      console.error('ไม่สามารถโหลดข้อมูลผู้ใช้:', err)
    }
  }

  const logout = () => {
    // ✅ ลบทั้ง token และ user data
    localStorage.removeItem('token')
    localStorage.removeItem('user')  // ← เพิ่มบรรทัดนี้!
    
    // หรือใช้ localStorage.clear() เพื่อลบทั้งหมด
    // localStorage.clear()
    
    token.value = null
    nickname.value = null
    imageUrl.value = null
    
    console.log('🔓 Logged out - localStorage cleared')
    router.push('/')
  }

  onMounted(() => {
    token.value = localStorage.getItem('token')
    if (token.value) fetchUser()
  })

  return { token, nickname, imageUrl, showDropdown, logout, fetchUser }
}