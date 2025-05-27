import { ref } from 'vue'

export const useRegister = () => {
  const form = ref({
    username: '',
    email: '',
    password: '',
    confirm_password: '',
    first_name: '',
    last_name: '',
    nickname: '',
  })

  const error = ref<Record<string, string>>({})
  const success = ref('')

  const resetForm = () => {
    form.value = {
      username: '',
      email: '',
      password: '',
      confirm_password: '',
      first_name: '',
      last_name: '',
      nickname: '',
    }
  }

  const register = async () => {
    error.value = {}
    success.value = ''

    if (form.value.password !== form.value.confirm_password) {
      error.value.confirm_password = '❌ รหัสผ่านไม่ตรงกัน'
      return
    }

    try {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form.value),
      })

      const json = await res.json()

      if (res.ok) {
        success.value = '✅ สมัครสมาชิกสำเร็จ'
        resetForm()
      } else {
        if (json.errors && typeof json.errors === 'object') {
          error.value = json.errors
        } else {
          error.value.general = json.error || '❌ เกิดข้อผิดพลาดในการสมัคร'
        }
      }
    } catch (err) {
      error.value.general = '❌ ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์ได้'
    }
  }

  return {
    form,
    error,
    success,
    register,
  }
}
