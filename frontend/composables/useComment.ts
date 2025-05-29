import type { Article } from "~/types/article"
import type { Comment } from "~/types/article"

// 🔍 Read + 💬 Comment: แสดงรายละเอียดบทความและจัดการความคิดเห็น
export function useComment() {
    const route = useRoute()
    const article = ref<Article | null>(null)
    const comments = ref<Comment[]>([])
    const newComment = ref('')
    
    // สำหรับการแก้ไขคอมเมนต์
    const editingCommentId = ref<number | null>(null)
    const editContent = ref('')
  
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
        await $fetch(`/api/articles/${route.params.slug}/comments`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: {
            content: newComment.value,
          },
        })

        newComment.value = ''
        // รีเฟรช comments หลังจากโพสต์สำเร็จ
        await fetchComments()
      } catch (err) {
        console.error('❌ Failed to post comment:', err)
        alert('ไม่สามารถโพสต์ความคิดเห็นได้ กรุณาลองใหม่อีกครั้ง')
      }
    }

    // เริ่มแก้ไขคอมเมนต์
    const startEditComment = (comment: Comment) => {
      editingCommentId.value = comment.id
      editContent.value = comment.content
    }

    // ยกเลิกการแก้ไข
    const cancelEdit = () => {
      editingCommentId.value = null
      editContent.value = ''
    }

// บันทึกการแก้ไขคอมเมนต์
const saveEditComment = async (commentId: number) => {
  const token = localStorage.getItem('token')
  if (!token) {
    alert('กรุณาเข้าสู่ระบบ')
    return
  }
  if (!editContent.value.trim()) {
    alert('กรุณากรอกข้อความความคิดเห็น')
    return
  }

  try {
    // เปลี่ยน URL ให้ตรงกับ backend structure
    await $fetch(`/api/articles/${route.params.slug}/comments/${commentId}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: {
        content: editContent.value,
      },
    })

    // รีเซ็ตการแก้ไข
    editingCommentId.value = null
    editContent.value = ''
    
    // รีเฟรช comments
    await fetchComments()
    
    alert('แก้ไขความคิดเห็นสำเร็จ')
  } catch (err) {
    console.error('❌ Failed to update comment:', err)
    alert('ไม่สามารถแก้ไขความคิดเห็นได้ กรุณาลองใหม่อีกครั้ง')
  }
}

// ลบคอมเมนต์
const deleteComment = async (commentId: number) => {
  const token = localStorage.getItem('token')
  if (!token) {
    alert('กรุณาเข้าสู่ระบบ')
    return
  }

  // ยืนยันการลบ
  if (!confirm('คุณแน่ใจหรือไม่ที่จะลบความคิดเห็นนี้?')) {
    return
  }

  try {
    // เปลี่ยน URL ให้ตรงกับ backend structure
    await $fetch(`/api/articles/${route.params.slug}/comments/${commentId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    // รีเฟรช comments
    await fetchComments()
    
    alert('ลบความคิดเห็นสำเร็จ')
  } catch (err) {
    console.error('❌ Failed to delete comment:', err)
    alert('ไม่สามารถลบความคิดเห็นได้ กรุณาลองใหม่อีกครั้ง')
  }
}

    // ตรวจสอบว่าเป็นเจ้าของคอมเมนต์หรือไม่
// วิธีง่ายๆ - ใช้ username เปรียบเทียบ
const isCommentOwner = (comment: Comment) => {
  const userStr = localStorage.getItem('user')
  if (!userStr) return false
  
  try {
    const currentUser = JSON.parse(userStr)
    
    // เปรียบเทียบ username (จาก debug ข้างบนเห็นว่ามี username)
    const isOwner = currentUser.username === comment.user?.username
    
    console.log('Current Username:', currentUser.username)
    console.log('Comment Username:', comment.user?.username) 
    console.log('Is Owner?', isOwner)
    
    return isOwner
  } catch {
    return false
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
      editingCommentId,
      editContent,
      fetchArticle,
      fetchComments,
      submitComment,
      startEditComment,
      cancelEdit,
      saveEditComment,
      deleteComment,
      isCommentOwner,
      formatDate,
    }
}