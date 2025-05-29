import { ref } from 'vue'
import { useRouter } from 'vue-router'

export const useLogin = () => {
  const emailOrUsername = ref('')
  const password = ref('')
  const error = ref('')
  const router = useRouter()

  const login = async () => {
    error.value = ''

    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: emailOrUsername.value,
          password: password.value,
        }),
      })

      const json = await res.json()

      if (res.ok && json.data.token) {
        // ✅ เก็บ token
        localStorage.setItem('token', json.data.token)
        
        // ✅ เก็บข้อมูล user (จุดสำคัญที่ขาดไป!)
        if (json.data.user) {
          localStorage.setItem('user', JSON.stringify(json.data.user))
        }
        
        router.push('/')
      } else {
        error.value = json.error || 'เข้าสู่ระบบล้มเหลว'
      }
    } catch (err) {
      error.value = 'ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์'
    }
  }

  return {
    emailOrUsername,
    password,
    error,
    login,
  }
}