import { ref, onMounted } from 'vue'
import { PlaceholdersImage } from '~/constants/Placeholders'
interface UserProfile {
  first_name: string
  last_name: string
  nickname: string
  bio?: string
  image?: string
}

export const useProfile = () => {
  const user = ref<UserProfile | null>(null)
  const form = ref<UserProfile>({
    first_name: '',
    last_name: '',
    nickname: '',
    bio: '',
  })
  const selectedFile = ref<File | null>(null)
  const previewImage = ref<string | null>(null)
  const success = ref('')
  const error = ref('')
  const loading = ref(false)

  const fetchProfile = async () => {
    loading.value = true
    try {
      const token = localStorage.getItem('token')
      if (!token) throw new Error('Missing token')

      const res = await $fetch<{ data: UserProfile }>('/api/user/', {
        headers: { Authorization: `Bearer ${token}` },
      })

      user.value = res.data
      form.value = { ...res.data }
      previewImage.value = res.data.image || null
    } catch (err) {
      console.error('❌ Error loading profile:', err)
      error.value = '❌ ไม่สามารถโหลดข้อมูลผู้ใช้ได้'
    } finally {
      loading.value = false
    }
  }

  const handleFileChange = (e: Event) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return

    if (file.size > 10 * 1024 * 1024) {
      error.value = PlaceholdersImage.imageNote
      selectedFile.value = null
      previewImage.value = user.value?.image || null
      return
    }

    selectedFile.value = file
    previewImage.value = URL.createObjectURL(file)
  }

  const updateProfile = async () => {
    success.value = ''
    error.value = ''

    try {
      const token = localStorage.getItem('token')
      if (!token) throw new Error('Missing token')

      const formData = new FormData()
      formData.append('first_name', form.value.first_name)
      formData.append('last_name', form.value.last_name)
      formData.append('nickname', form.value.nickname)
      formData.append('bio', form.value.bio || '')
      if (selectedFile.value) {
        formData.append('avatar', selectedFile.value)
      }

      const res = await fetch('/api/user', {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      })

      const json = await res.json()
      if (res.ok) {
        success.value = '✅ แก้ไขข้อมูลสำเร็จ'
        user.value = json.data
        form.value = { ...json.data }
        previewImage.value = json.data.image || null
        selectedFile.value = null
      } else {
        error.value = json.message || '❌ แก้ไขไม่สำเร็จ'
      }
    } catch (err) {
      console.error('❌ Error updating profile:', err)
      error.value = '❌ ไม่สามารถบันทึกข้อมูลได้'
    }
  }

  onMounted(fetchProfile)

  return {
    user,
    form,
    selectedFile,
    previewImage,
    success,
    error,
    loading,
    handleFileChange,
    updateProfile,
  }
}