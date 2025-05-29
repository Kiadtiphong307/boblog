import type { Article } from "~/types/article"
import type { Comment } from "~/types/article"

// 🔍 Read + 💬 Comment: แสดงรายละเอียดบทความและจัดการความคิดเห็น
export function useComment() {
    const route = useRoute()
    const article = ref<Article | null>(null)
    const comments = ref<Comment[]>([])
    const newComment = ref('')
  
    // ดึงรายละเอียดบทความตาม slug
    const fetchArticle = async () => {
      const token = localStorage.getItem('token')
      try {
        const res = await $fetch<{data: Article}>(`/api/articles/${route.params.slug}`, {
          headers: { Authorization: `Bearer ${token}` },
        })
        article.value = (res as any).data
      } catch (err) {
        console.error('❌ Failed to load article:', err)
      }
    }
  
  // ดึงความคิดเห็นของบทความนั้น ๆ
  const fetchComments = async () => {
    const token = localStorage.getItem('token')
    try {
      const res = await $fetch(`/api/articles/${route.params.slug}/comments`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      
      // ตรวจสอบโครงสร้างข้อมูลที่ได้รับ
      console.log('Response from API:', res)
      
      // ปรับการดึงข้อมูลให้ถูกต้อง
      if (res && typeof res === 'object') {
        // ถ้า response มี property data
        if ('data' in res) {
          comments.value = (res as any).data || []
        } 
        // ถ้า response เป็น array โดยตรง
        else if (Array.isArray(res)) {
          comments.value = res
        }
        // ถ้า response มี property comments
        else if ('comments' in res) {
          comments.value = (res as any).comments || []
        }
        // ถ้าไม่ตรงเงื่อนไขใดเลย
        else {
          comments.value = []
          console.warn('Unexpected response structure:', res)
        }
      } else {
        comments.value = []
      }
      
      console.log('Comments loaded:', comments.value)
    } catch (err) {
      console.error('❌ Failed to load comments:', err)
      comments.value = []
    }
  }
  
  // ส่งความคิดเห็นใหม่
  const submitComment = async () => {
    const token = localStorage.getItem('token')
    if (!token) {
      alert('กรุณาเข้าสู่ระบบก่อนแสดงความคิดเห็น')
      return
    }
    if (!newComment.value.trim()) return
  
    try {
      const response = await $fetch(`/api/articles/${route.params.slug}/comments`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: {
          content: newComment.value,
        },
      })
      
      console.log('Comment posted:', response)
      
      newComment.value = ''
      // รีเฟรช comments หลังจากโพสต์สำเร็จ
      await fetchComments()
    } catch (err) {
      console.error('❌ Failed to post comment:', err)
      alert('ไม่สามารถโพสต์ความคิดเห็นได้ กรุณาลองใหม่อีกครั้ง')
    }
  }
  
    const formatDate = (iso: string) =>
      new Date(iso).toLocaleString('th-TH', {
        dateStyle: 'medium',
        timeStyle: 'short',
      })
  
    onMounted(() => {
      fetchArticle()
      fetchComments()
    })
  
    return {
      article,
      comments,
      newComment,
      fetchArticle,
      fetchComments,
      submitComment,
      formatDate,
    }
  }