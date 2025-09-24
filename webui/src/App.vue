<template>
  <div class="app-container">
    <div class="chat-app">
      <div class="header">
        <h1 class="title">WASA Chat</h1>
        <div v-if="!isLoggedIn" class="login-section">
          <input
            v-model="username"
            placeholder="Enter username"
            class="input"
            @keyup.enter="login"
          />
          <button @click="login" class="btn-primary">Login</button>
        </div>
        <div v-else class="user-info">
          <span>Welcome, {{ currentUser?.username }}</span>
          <button @click="showProfileModal = true" class="btn-secondary">Profile</button>
          <button @click="logout" class="btn-secondary">Logout</button>
        </div>
      </div>

      <div v-if="isLoggedIn" class="main-content">
        <div class="sidebar">
          <div class="conversations-header">
            <h3>Conversations</h3>
            <button @click="showCreateConversationModal = true" class="btn-small">+ New</button>
          </div>
          <div class="conversations-list">
            <div
              v-for="conv in conversations"
              :key="conv.id"
              :class="['conversation-item', { active: selectedConversation?.id === conv.id }]"
              @click="selectConversation(conv)"
            >
              <div class="conv-name">{{ conv.name }}</div>
              <div class="conv-type">{{ conv.type }}</div>
            </div>
          </div>
        </div>

        <div class="chat-area">
          <div v-if="!selectedConversation" class="no-conversation">
            <p>Select a conversation to start chatting</p>
          </div>
          <div v-else class="conversation-view">
            <div class="conversation-header">
              <h3>{{ selectedConversation.name }}</h3>
              <div class="conversation-actions">
                <button @click="showGroupSettings = true" class="btn-small">Settings</button>
              </div>
            </div>
            <div class="messages-container" ref="messagesContainer">
              <div
                v-for="message in messages"
                :key="message.id"
                class="message"
                :class="{ 'own-message': message.author === currentUser?.username }"
              >
                <div class="message-content">
                  <div class="message-body">{{ message.body }}</div>
                  <div class="message-meta">
                    <span class="message-author">{{ message.author }}</span>
                    <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                  </div>
                  <div v-if="message.comments && message.comments.length > 0" class="message-comments">
                    <div v-for="comment in message.comments" :key="comment.id" class="comment">
                      <span class="comment-author">{{ comment.author }}:</span>
                      <span class="comment-body">{{ comment.body }}</span>
                    </div>
                  </div>
                  <div class="message-actions">
                    <button @click="forwardMessage(message)" class="btn-tiny">Forward</button>
                    <button @click="commentMessage(message)" class="btn-tiny">Comment</button>
                    <button v-if="message.author === currentUser?.username" @click="deleteMessage(message)" class="btn-tiny danger">Delete</button>
                  </div>
                </div>
              </div>
            </div>
            <div class="message-input">
              <input
                v-model="newMessage"
                placeholder="Type a message..."
                class="input"
                @keyup.enter="sendMessage"
              />
              <button @click="sendMessage" class="btn-primary">Send</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <div v-if="showCreateConversationModal" class="modal-overlay" @click="showCreateConversationModal = false">
      <div class="modal" @click.stop>
        <h3>Create Conversation</h3>
        <input v-model="newConversationName" placeholder="Conversation name" class="input" />
        <select v-model="newConversationType" class="input">
          <option value="group">Group</option>
          <option value="direct">Direct</option>
        </select>
        <div class="modal-actions">
          <button @click="createConversation" class="btn-primary">Create</button>
          <button @click="showCreateConversationModal = false" class="btn-secondary">Cancel</button>
        </div>
      </div>
    </div>

    <div v-if="showProfileModal" class="modal-overlay" @click="showProfileModal = false">
      <div class="modal" @click.stop>
        <h3>Profile Settings</h3>
        <div class="form-group">
          <label>Username:</label>
          <input v-model="profileUsername" placeholder="Username" class="input" />
          <button @click="updateUsername" class="btn-primary">Update</button>
        </div>
        <div class="form-group">
          <label>Profile Photo:</label>
          <input type="file" @change="handlePhotoUpload" accept="image/*" class="input" />
        </div>
        <div class="modal-actions">
          <button @click="showProfileModal = false" class="btn-secondary">Close</button>
        </div>
      </div>
    </div>

    <div v-if="showGroupSettings" class="modal-overlay" @click="showGroupSettings = false">
      <div class="modal" @click.stop>
        <h3>Group Settings</h3>
        <div class="form-group">
          <label>Group Name:</label>
          <input v-model="groupName" placeholder="Group name" class="input" />
          <button @click="updateGroupName" class="btn-primary">Update</button>
        </div>
        <div class="form-group">
          <label>Group Photo:</label>
          <input type="file" @change="handleGroupPhotoUpload" accept="image/*" class="input" />
        </div>
        <div class="form-group">
          <label>Add User:</label>
          <input v-model="userToAdd" placeholder="User ID" class="input" />
          <button @click="addUserToGroup" class="btn-primary">Add</button>
        </div>
        <div class="modal-actions">
          <button @click="leaveGroup" class="btn-danger">Leave Group</button>
          <button @click="showGroupSettings = false" class="btn-secondary">Close</button>
        </div>
      </div>
    </div>

    <div v-if="showCommentModal" class="modal-overlay" @click="showCommentModal = false">
      <div class="modal" @click.stop>
        <h3>Add Comment</h3>
        <textarea v-model="commentText" placeholder="Your comment..." class="input"></textarea>
        <div class="modal-actions">
          <button @click="submitComment" class="btn-primary">Submit</button>
          <button @click="showCommentModal = false" class="btn-secondary">Cancel</button>
        </div>
      </div>
    </div>

    <div v-if="showForwardModal" class="modal-overlay" @click="showForwardModal = false">
      <div class="modal" @click.stop>
        <h3>Forward Message</h3>
        <select v-model="forwardToConversation" class="input">
          <option value="">Select conversation</option>
          <option v-for="conv in conversations" :key="conv.id" :value="conv.id">
            {{ conv.name }}
          </option>
        </select>
        <div class="modal-actions">
          <button @click="submitForward" class="btn-primary">Forward</button>
          <button @click="showForwardModal = false" class="btn-secondary">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || (() => {
  const host = window.location.hostname
  if (host.endsWith('.app.github.dev')) {
    return `https://${host.replace(/-\d+\.app.github.dev$/, '-8080.app.github.dev')}`
  }
  return `${window.location.protocol}//${host}:8080`
})()

console.log('API_BASE computed →', API_BASE)

// State
const isLoggedIn = ref(false)
const currentUser = ref(null)
const token = ref('')
const username = ref('')
const conversations = ref([])
const selectedConversation = ref(null)
const messages = ref([])
const newMessage = ref('')

// Modal states
const showCreateConversationModal = ref(false)
const showProfileModal = ref(false)
const showGroupSettings = ref(false)
const showCommentModal = ref(false)
const showForwardModal = ref(false)

// Form data
const newConversationName = ref('')
const newConversationType = ref('group')
const profileUsername = ref('')
const groupName = ref('')
const userToAdd = ref('')
const commentText = ref('')
const forwardToConversation = ref('')
const messageToComment = ref(null)
const messageToForward = ref(null)

// Auth functions
async function login() {
  if (!username.value) return
  
  try {
    const { data } = await axios.post(`${API_BASE}/session`, {
      name: username.value
    })
    
    token.value = data.identifier
    isLoggedIn.value = true
    currentUser.value = { username: username.value }
    profileUsername.value = username.value
    
    // Save token to localStorage
    localStorage.setItem('token', token.value)
    
    // Load conversations
    await loadConversations()
  } catch (e) {
    console.error('Login error:', e)
    alert('Login failed')
  }
}

function logout() {
  isLoggedIn.value = false
  currentUser.value = null
  token.value = ''
  conversations.value = []
  selectedConversation.value = null
  messages.value = []
  localStorage.removeItem('token')
}

// Conversation functions
async function loadConversations() {
  try {
    const { data } = await axios.get(`${API_BASE}/conversations`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    conversations.value = data || []
  } catch (e) {
    console.error('Load conversations error:', e)
    conversations.value = []
  }
}

async function createConversation() {
  if (!newConversationName.value) return
  
  try {
    const { data } = await axios.post(`${API_BASE}/conversations`, {
      name: newConversationName.value,
      type: newConversationType.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    conversations.value.push(data)
    newConversationName.value = ''
    showCreateConversationModal.value = false
  } catch (e) {
    console.error('Create conversation error:', e)
    alert('Failed to create conversation')
  }
}

async function selectConversation(conv) {
  selectedConversation.value = conv
  groupName.value = conv.name
  await loadMessages(conv.id)
}

async function loadMessages(conversationId) {
  try {
    const { data } = await axios.get(`${API_BASE}/conversations/${conversationId}`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    messages.value = data.messages || []
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error('Load messages error:', e)
  }
}

// Message functions
async function sendMessage() {
  if (!newMessage.value || !selectedConversation.value) return
  
  try {
    const { data } = await axios.post(`${API_BASE}/messages`, {
      conversationId: selectedConversation.value.id,
      body: newMessage.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    messages.value.push(data)
    newMessage.value = ''
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error('Send message error:', e)
    alert('Failed to send message')
  }
}

async function forwardMessage(message) {
  messageToForward.value = message
  showForwardModal.value = true
}

async function submitForward() {
  if (!messageToForward.value || !forwardToConversation.value) return
  
  try {
    await axios.post(`${API_BASE}/messages/${messageToForward.value.id}/forward`, {
      conversationId: forwardToConversation.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    showForwardModal.value = false
    forwardToConversation.value = ''
    messageToForward.value = null
    alert('Message forwarded successfully')
  } catch (e) {
    console.error('Forward message error:', e)
    alert('Failed to forward message')
  }
}

async function commentMessage(message) {
  messageToComment.value = message
  showCommentModal.value = true
}

async function submitComment() {
  if (!messageToComment.value || !commentText.value) return
  
  try {
    await axios.post(`${API_BASE}/messages/${messageToComment.value.id}/comment`, {
      comment: commentText.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    // Reload messages to show the comment
    await loadMessages(selectedConversation.value.id)
    showCommentModal.value = false
    commentText.value = ''
    messageToComment.value = null
  } catch (e) {
    console.error('Comment message error:', e)
    alert('Failed to add comment')
  }
}

async function deleteMessage(message) {
  if (!confirm('Are you sure you want to delete this message?')) return
  
  try {
    await axios.delete(`${API_BASE}/messages/${message.id}`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    // Remove message from local state
    messages.value = messages.value.filter(m => m.id !== message.id)
  } catch (e) {
    console.error('Delete message error:', e)
    alert('Failed to delete message')
  }
}

// Profile functions
async function updateUsername() {
  if (!profileUsername.value) return
  
  try {
    await axios.put(`${API_BASE}/user/username`, {
      username: profileUsername.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    currentUser.value.username = profileUsername.value
    alert('Username updated successfully')
  } catch (e) {
    console.error('Update username error:', e)
    alert('Failed to update username')
  }
}

async function handlePhotoUpload(event) {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    await axios.put(`${API_BASE}/user/photo`, formData, {
      headers: { 
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'multipart/form-data'
      }
    })
    
    alert('Photo updated successfully')
  } catch (e) {
    console.error('Update photo error:', e)
    alert('Failed to update photo')
  }
}

// Group functions
async function updateGroupName() {
  if (!groupName.value || !selectedConversation.value) return
  
  try {
    await axios.put(`${API_BASE}/groups/${selectedConversation.value.id}/name`, {
      name: groupName.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    selectedConversation.value.name = groupName.value
    alert('Group name updated successfully')
  } catch (e) {
    console.error('Update group name error:', e)
    alert('Failed to update group name')
  }
}

async function handleGroupPhotoUpload(event) {
  const file = event.target.files[0]
  if (!file || !selectedConversation.value) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    await axios.put(`${API_BASE}/groups/${selectedConversation.value.id}/photo`, formData, {
      headers: { 
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'multipart/form-data'
      }
    })
    
    alert('Group photo updated successfully')
  } catch (e) {
    console.error('Update group photo error:', e)
    alert('Failed to update group photo')
  }
}

async function addUserToGroup() {
  if (!userToAdd.value || !selectedConversation.value) return
  
  try {
    await axios.post(`${API_BASE}/groups/${selectedConversation.value.id}/add`, {
      userId: userToAdd.value
    }, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    userToAdd.value = ''
    alert('User added to group successfully')
  } catch (e) {
    console.error('Add user to group error:', e)
    alert('Failed to add user to group')
  }
}

async function leaveGroup() {
  if (!selectedConversation.value || !confirm('Are you sure you want to leave this group?')) return
  
  try {
    await axios.post(`${API_BASE}/groups/${selectedConversation.value.id}/leave`, {}, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    
    // Remove conversation from local state
    conversations.value = conversations.value.filter(c => c.id !== selectedConversation.value.id)
    selectedConversation.value = null
    messages.value = []
    showGroupSettings.value = false
    alert('Left group successfully')
  } catch (e) {
    console.error('Leave group error:', e)
    alert('Failed to leave group')
  }
}

// Utility functions
function scrollToBottom() {
  const container = document.querySelector('.messages-container')
  if (container) {
    container.scrollTop = container.scrollHeight
  }
}

function formatTime(timestamp) {
  return new Date(timestamp).toLocaleTimeString()
}

// Initialize
onMounted(async () => {
  // Check if user is already logged in
  const savedToken = localStorage.getItem('token')
  if (savedToken) {
    token.value = savedToken
    isLoggedIn.value = true
    // We need to get user info, but for now just set a placeholder
    currentUser.value = { username: 'User' }
    await loadConversations()
  }
})
</script>

<style>
.app-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #121212;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  margin: 0;
}

.chat-app {
  background: #1e1e1e;
  border-radius: 1.5rem;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  width: 100%;
  max-width: 1200px;
  height: 80vh;
  color: #fff;
  display: flex;
  flex-direction: column;
}

.header {
  padding: 1.5rem;
  border-bottom: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 1.5rem;
  margin: 0;
}

.login-section {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.user-info {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 300px;
  border-right: 1px solid #333;
  display: flex;
  flex-direction: column;
}

.conversations-header {
  padding: 1rem;
  border-bottom: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conversations-header h3 {
  margin: 0;
  font-size: 1rem;
}

.conversations-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  padding: 1rem;
  cursor: pointer;
  border-bottom: 1px solid #2a2a2a;
  transition: background 0.2s;
}

.conversation-item:hover {
  background: #2a2a2a;
}

.conversation-item.active {
  background: #5e5ce6;
}

.conv-name {
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.conv-type {
  font-size: 0.8rem;
  color: #999;
  text-transform: capitalize;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.no-conversation {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}

.conversation-view {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.conversation-header {
  padding: 1rem;
  border-bottom: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conversation-header h3 {
  margin: 0;
}

.conversation-actions {
  display: flex;
  gap: 0.5rem;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message {
  display: flex;
  justify-content: flex-start;
}

.message.own-message {
  justify-content: flex-end;
}

.message-content {
  max-width: 70%;
  background: #2a2a2a;
  padding: 0.75rem 1rem;
  border-radius: 1rem;
}

.own-message .message-content {
  background: #5e5ce6;
}

.message-body {
  margin-bottom: 0.25rem;
}

.message-meta {
  font-size: 0.75rem;
  color: #999;
  display: flex;
  gap: 0.5rem;
}

.message-comments {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid #444;
}

.comment {
  font-size: 0.8rem;
  margin-bottom: 0.25rem;
}

.comment-author {
  font-weight: 500;
  color: #5e5ce6;
}

.message-actions {
  margin-top: 0.5rem;
  display: flex;
  gap: 0.5rem;
}

.btn-tiny {
  padding: 0.25rem 0.5rem;
  background: transparent;
  border: 1px solid #666;
  border-radius: 0.25rem;
  color: #999;
  font-size: 0.7rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-tiny:hover {
  background: #666;
  color: #fff;
}

.btn-tiny.danger {
  border-color: #ff4444;
  color: #ff4444;
}

.btn-tiny.danger:hover {
  background: #ff4444;
  color: #fff;
}

.message-input {
  padding: 1rem;
  border-top: 1px solid #333;
  display: flex;
  gap: 1rem;
}

.input {
  flex: 1;
  padding: 0.75rem 1rem;
  background: #2a2a2a;
  border: none;
  border-radius: 1rem;
  color: #fff;
  font-size: 1rem;
  outline: none;
}

.input:focus {
  box-shadow: 0 0 0 2px #5e5ce6;
}

.btn-primary {
  padding: 0.75rem 1.5rem;
  background: #5e5ce6;
  border: none;
  border-radius: 1rem;
  color: #fff;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #7a79f8;
}

.btn-secondary {
  padding: 0.5rem 1rem;
  background: transparent;
  border: 1px solid #5e5ce6;
  border-radius: 1rem;
  color: #5e5ce6;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #5e5ce6;
  color: #fff;
}

.btn-small {
  padding: 0.5rem;
  background: #5e5ce6;
  border: none;
  border-radius: 0.5rem;
  color: #fff;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-small:hover {
  background: #7a79f8;
}

.btn-danger {
  padding: 0.5rem 1rem;
  background: #ff4444;
  border: none;
  border-radius: 1rem;
  color: #fff;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-danger:hover {
  background: #ff6666;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #1e1e1e;
  border-radius: 1rem;
  padding: 2rem;
  min-width: 400px;
  max-width: 500px;
  color: #fff;
}

.modal h3 {
  margin: 0 0 1rem 0;
  font-size: 1.2rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

textarea.input {
  min-height: 100px;
  resize: vertical;
}
</style>