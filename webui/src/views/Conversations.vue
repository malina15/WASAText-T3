<template>
  <div class="conversations-container">
    <div class="conversations-header">
      <h1>Conversations</h1>
      <div class="header-actions">
        <button @click="showCreateModal = true" class="btn-primary">
          <span class="icon">+</span>
          New Conversation
        </button>
        <button @click="goToProfile" class="btn-secondary">
          <span class="icon">👤</span>
          Profile
        </button>
        <button @click="logout" class="btn-secondary">
          <span class="icon">🚪</span>
          Logout
        </button>
      </div>
    </div>

    <div class="conversations-content">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading conversations...</p>
      </div>

      <div v-else-if="conversations.length === 0" class="empty-state">
        <div class="empty-icon">💬</div>
        <h3>No conversations yet</h3>
        <p>Start a new conversation to begin chatting!</p>
        <button @click="showCreateModal = true" class="btn-primary">
          Create Conversation
        </button>
      </div>

      <div v-else class="conversations-list">
        <div
          v-for="conversation in conversations"
          :key="conversation.id"
          class="conversation-item"
          @click="selectConversation(conversation)"
        >
          <div class="conversation-info">
            <h3 class="conversation-name">{{ conversation.name }}</h3>
            <span class="conversation-type">{{ conversation.type }}</span>
          </div>
          <div class="conversation-meta">
            <span class="conversation-date">{{ formatDate(conversation.createdAt) }}</span>
            <span class="conversation-arrow">→</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Conversation Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h3>Create New Conversation</h3>
        <form @submit.prevent="createConversation">
          <div class="form-group">
            <label for="conversation-name">Conversation Name</label>
            <input
              id="conversation-name"
              v-model="newConversationName"
              type="text"
              placeholder="Enter conversation name"
              class="form-input"
              required
              minlength="1"
              maxlength="100"
            />
          </div>
          
          <div class="form-group">
            <label for="conversation-type">Type</label>
            <select
              id="conversation-type"
              v-model="newConversationType"
              class="form-input"
            >
              <option value="group">Group</option>
              <option value="direct">Direct Message</option>
            </select>
          </div>

          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn-primary" :disabled="!newConversationName">
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()
const conversations = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const newConversationName = ref('')
const newConversationType = ref('group')

onMounted(async () => {
  await loadConversations()
})

async function loadConversations() {
  loading.value = true
  try {
    conversations.value = await apiService.getMyConversations()
  } catch (error) {
    console.error('Failed to load conversations:', error)
    alert('Failed to load conversations')
  } finally {
    loading.value = false
  }
}

async function createConversation() {
  if (!newConversationName.value) return
  
  try {
    const conversation = await apiService.createConversation(
      newConversationName.value,
      newConversationType.value
    )
    
    conversations.value.push(conversation)
    newConversationName.value = ''
    newConversationType.value = 'group'
    showCreateModal.value = false
    
    // Navigate to the new conversation
    router.push(`/conversations/${conversation.id}`)
  } catch (error) {
    console.error('Failed to create conversation:', error)
    alert('Failed to create conversation')
  }
}

function selectConversation(conversation) {
  router.push(`/conversations/${conversation.id}`)
}

function goToProfile() {
  router.push('/profile')
}

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/')
}

function formatDate(dateString) {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 1) return 'Today'
  if (diffDays === 2) return 'Yesterday'
  if (diffDays <= 7) return `${diffDays - 1} days ago`
  
  return date.toLocaleDateString()
}
</script>

<style scoped>
.conversations-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  min-height: 100vh;
  background: #f8fafc;
}

.conversations-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e2e8f0;
}

.conversations-header h1 {
  margin: 0;
  color: #2d3748;
  font-size: 2rem;
  font-weight: 700;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
}

.btn-secondary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #667eea;
  color: white;
}

.icon {
  font-size: 1.2rem;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: #718096;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #718096;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  margin: 0 0 0.5rem 0;
  color: #2d3748;
}

.empty-state p {
  margin: 0 0 2rem 0;
}

.conversations-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.conversation-item {
  background: white;
  border-radius: 0.75rem;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conversation-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.conversation-info {
  flex: 1;
}

.conversation-name {
  margin: 0 0 0.5rem 0;
  color: #2d3748;
  font-size: 1.25rem;
  font-weight: 600;
}

.conversation-type {
  background: #e2e8f0;
  color: #4a5568;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  text-transform: capitalize;
}

.conversation-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.5rem;
}

.conversation-date {
  color: #718096;
  font-size: 0.875rem;
}

.conversation-arrow {
  color: #cbd5e0;
  font-size: 1.5rem;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  width: 100%;
  max-width: 500px;
  margin: 1rem;
}

.modal h3 {
  margin: 0 0 1.5rem 0;
  color: #2d3748;
  font-size: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #2d3748;
}

.form-input {
  width: 100%;
  padding: 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}
</style>
