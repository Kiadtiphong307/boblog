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
      }
    } catch (err) {
      console.error('ไม่สามารถโหลดข้อมูลผู้ใช้:', err)
    }
  }

  const logout = () => {
    localStorage.removeItem('token')
    token.value = null
    nickname.value = null
    imageUrl.value = null
    router.push('/')
  }

  onMounted(() => {
    token.value = localStorage.getItem('token')
    if (token.value) fetchUser()
  })

  return { token, nickname, imageUrl, showDropdown, logout, fetchUser }
}