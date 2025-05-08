<template>
  <div class="test-auth-container">
    <div class="card">
      <h1>🔐 Test Auth/JWT Flow</h1>
      <form @submit.prevent="register" class="form">
        <h3>Register</h3>
        <input v-model="regUser" placeholder="Username" required />
        <input v-model="regPass" type="password" placeholder="Password" required />
        <button type="submit" class="btn btn-primary">Register</button>
      </form>
      <form @submit.prevent="login" class="form">
        <h3>Login</h3>
        <input v-model="loginUser" placeholder="Username" required />
        <input v-model="loginPass" type="password" placeholder="Password" required />
        <button type="submit" class="btn btn-primary">Login</button>
      </form>
      <div class="btn-group">
        <button @click="getProfile" class="btn">Get Profile</button>
        <button @click="createPost" class="btn">Create Post</button>
        <button @click="getPosts" class="btn">Get Posts</button>
        <button @click="updatePost" class="btn">Update Post</button>
        <button @click="deletePost" class="btn">Delete Post</button>
        <button @click="logout" class="btn btn-danger">Logout</button>
      </div>
      <div v-if="jwt" class="jwt-card">
        <strong>JWT:</strong>
        <pre>{{ jwt }}</pre>
      </div>
      <div v-if="result" class="result-card">
        <pre>{{ result }}</pre>
      </div>
      <div v-if="error" class="error-msg">{{ error }}</div>
    </div>
    <footer class="footer">
      <small>Boblog Test UI &copy; 2025</small>
    </footer>
  </div>
</template>

<style scoped>
.test-auth-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #e0e7ff 0%, #f0fdfa 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 16px 0 rgba(0,0,0,0.08);
  padding: 2rem 1.5rem;
  width: 100%;
  max-width: 430px;
  margin: 2rem auto 1rem auto;
}
h1 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #312e81;
  font-size: 1.7rem;
}
.form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1.2rem;
}
.form input {
  padding: 0.5rem;
  border: 1px solid #c7d2fe;
  border-radius: 6px;
  font-size: 1rem;
  background: #f8fafc;
  outline: none;
  transition: border 0.2s;
}
.form input:focus {
  border: 1.5px solid #6366f1;
}
.btn-group {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 1.3rem;
}
.btn {
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 6px;
  background: #6366f1;
  color: #fff;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}
.btn:hover {
  background: #4338ca;
}
.btn-danger {
  background: #f43f5e;
}
.btn-danger:hover {
  background: #be123c;
}
.btn-primary {
  background: #10b981;
}
.btn-primary:hover {
  background: #047857;
}
.jwt-card {
  background: #f1f5f9;
  border-left: 4px solid #6366f1;
  padding: 0.7rem 1rem;
  margin-bottom: 1rem;
  border-radius: 6px;
  word-break: break-all;
  font-size: 0.98rem;
}
.result-card {
  background: #e0f2fe;
  border-left: 4px solid #0ea5e9;
  padding: 1rem 1rem;
  margin-bottom: 1rem;
  border-radius: 6px;
  font-size: 0.97rem;
  color: #075985;
}
.error-msg {
  color: #dc2626;
  background: #fef2f2;
  padding: 0.7rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  border-left: 4px solid #dc2626;
}
.footer {
  text-align: center;
  color: #64748b;
  margin-top: 1.5rem;
}
@media (max-width: 600px) {
  .card {
    padding: 1rem 0.3rem;
    max-width: 98vw;
  }
}
</style>

<script setup>
import { ref } from 'vue'

const regUser = ref('')
const regPass = ref('')
const loginUser = ref('')
const loginPass = ref('')
const result = ref('')
const error = ref('')
const jwt = ref('')
const postId = ref('')

// ปรับ endpoint ให้ตรงกับ backend จริง (เช่น /api/register, /api/login, /api/profile)
const API_BASE = 'http://localhost:8080' // กรณีเดียวกับ nuxt proxy, ถ้าไม่ใช้ proxy ให้ใส่ URL backend

async function register() {
  error.value = ''
  result.value = ''
  try {
    const res = await fetch(`${API_BASE}/api/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: regUser.value, password: regPass.value })
    })
    if (!res.ok) throw new Error('Register failed')
    result.value = 'Register success!'
  } catch (e) {
    error.value = e.message
  }
}

async function login() {
  error.value = ''
  result.value = ''
  try {
    const res = await fetch(`${API_BASE}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: loginUser.value, password: loginPass.value })
    })
    if (!res.ok) throw new Error('Login failed')
    const data = await res.json()
    jwt.value = data.token
    localStorage.setItem('jwt', data.token)
    result.value = 'Login success! JWT saved.'
  } catch (e) {
    error.value = e.message
  }
}

async function getProfile() {
  error.value = ''
  result.value = ''
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('No JWT token, กรุณา Login ก่อน')
    const res = await fetch(`${API_BASE}/api/profile`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Get profile failed (JWT อาจหมดอายุหรือผิด)')
    const data = await res.json()
    result.value = JSON.stringify(data, null, 2)
  } catch (e) {
    error.value = e.message
  }
}

async function createPost() {
  error.value = ''
  result.value = ''
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('No JWT token, กรุณา Login ก่อน')
    const res = await fetch(`${API_BASE}/api/posts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ title: 'Test Post', content: 'This is a test post' })
    })
    if (!res.ok) throw new Error('Create post failed')
    const data = await res.json()
    postId.value = data.id
    result.value = 'Post created successfully!'
  } catch (e) {
    error.value = e.message
  }
}

async function getPosts() {
  error.value = ''
  result.value = ''
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('No JWT token, กรุณา Login ก่อน')
    const res = await fetch(`${API_BASE}/api/posts`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Get posts failed')
    const data = await res.json()
    result.value = JSON.stringify(data, null, 2)
  } catch (e) {
    error.value = e.message
  }
}

async function updatePost() {
  error.value = ''
  result.value = ''
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('No JWT token, กรุณา Login ก่อน')
    const res = await fetch(`${API_BASE}/api/posts/${postId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ title: 'Updated Test Post', content: 'This is an updated test post' })
    })
    if (!res.ok) throw new Error('Update post failed')
    result.value = 'Post updated successfully!'
  } catch (e) {
    error.value = e.message
  }
}

async function deletePost() {
  error.value = ''
  result.value = ''
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('No JWT token, กรุณา Login ก่อน')
    const res = await fetch(`${API_BASE}/api/posts/${postId.value}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Delete post failed')
    result.value = 'Post deleted successfully!'
  } catch (e) {
    error.value = e.message
  }
}

async function logout() {
  error.value = ''
  result.value = ''
  try {
    localStorage.removeItem('jwt')
    jwt.value = ''
    result.value = 'Logged out successfully!'
  } catch (e) {
    error.value = e.message
  }
}
</script>