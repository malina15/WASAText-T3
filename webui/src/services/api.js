import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || (() => {
  const host = window.location.hostname
  if (host.endsWith('.app.github.dev')) {
    return `https://${host.replace(/-\d+\.app.github.dev$/, '-8080.app.github.dev')}`
  }
  return `${window.location.protocol}//${host}:8080`
})()

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE,
  timeout: 10000,
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/'
    }
    return Promise.reject(error)
  }
)

// API Service class with all endpoints
class ApiService {
  // Authentication
  async login(username) {
    const response = await api.post('/session', { name: username })
    return response.data
  }

  // User management
  async setMyUserName(username) {
    const response = await api.put('/user/username', { username })
    return response.data
  }

  async setMyPhoto(file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.put('/user/photo', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  }

  // Conversations
  async getMyConversations() {
    const response = await api.get('/conversations')
    return response.data
  }

  async createConversation(name, type = 'group') {
    const response = await api.post('/conversations', { name, type })
    return response.data
  }

  async getConversation(id) {
    const response = await api.get(`/conversations/${id}`)
    return response.data
  }

  // Messages
  async sendMessage(conversationId, body) {
    const response = await api.post('/messages', { conversationId, body })
    return response.data
  }

  async forwardMessage(messageId, conversationId) {
    const response = await api.post(`/messages/${messageId}/forward`, { conversationId })
    return response.data
  }

  async commentMessage(messageId, comment) {
    const response = await api.post(`/messages/${messageId}/comment`, { comment })
    return response.data
  }

  async uncommentMessage(messageId) {
    const response = await api.post(`/messages/${messageId}/uncomment`)
    return response.data
  }

  async deleteMessage(messageId) {
    const response = await api.delete(`/messages/${messageId}`)
    return response.data
  }

  // Group management
  async addToGroup(groupId, userId) {
    const response = await api.post(`/groups/${groupId}/add`, { userId })
    return response.data
  }

  async leaveGroup(groupId) {
    const response = await api.post(`/groups/${groupId}/leave`)
    return response.data
  }

  async setGroupName(groupId, name) {
    const response = await api.put(`/groups/${groupId}/name`, { name })
    return response.data
  }

  async setGroupPhoto(groupId, file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.put(`/groups/${groupId}/photo`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  }
}

// Export singleton instance
export default new ApiService()
