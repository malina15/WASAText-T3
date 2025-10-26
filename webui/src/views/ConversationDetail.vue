<template>
  <div class="conversation-detail-container">
    <div class="conversation-header">
      <button @click="goBack" class="back-button">
        <span class="icon">←</span>
        Back
      </button>
      <div class="conversation-info">
        <h1>{{ conversation?.name || 'Loading...' }}</h1>
        <span class="conversation-type">{{ conversation?.type }}</span>
      </div>
      <div class="header-actions">
        <button @click="showGroupSettings = true" class="btn-secondary">
          <span class="icon">⚙️</span>
          Settings
        </button>
      </div>
    </div>

    <div class="conversation-content">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading messages...</p>
      </div>

      <div v-else class="messages-container" ref="messagesContainer">
        <div
          v-for="message in messages"
          :key="message.id"
          class="message"
          :class="{ 'own-message': message.author === currentUsername }"
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
              <button @click="forwardMessage(message)" class="btn-tiny">
                <span class="icon">↗️</span>
                Forward
              </button>
              <button @click="commentMessage(message)" class="btn-tiny">
                <span class="icon">💬</span>
                Comment
              </button>
              <button 
                v-if="message.author === currentUsername" 
                @click="deleteMessage(message)" 
                class="btn-tiny danger"
              >
                <span class="icon">🗑️</span>
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="message-input-container">
        <div class="message-input">
          <input
            v-model="newMessage"
            placeholder="Type a message..."
            class="input"
            @keyup.enter="sendMessage"
            :disabled="loading"
          />
          <button 
            @click="sendMessage" 
            class="send-button"
            :disabled="!newMessage || loading"
          >
            <span class="icon">📤</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Group Settings Modal -->
    <div v-if="showGroupSettings" class="modal-overlay" @click="showGroupSettings = false">
      <div class="modal" @click.stop>
        <h3>Group Settings</h3>
        
        <div class="form-group">
          <label>Group Name</label>
          <input
            v-model="groupName"
            type="text"
            placeholder="Group name"
            class="form-input"
          />
          <button @click="updateGroupName" class="btn-primary">Update Name</button>
        </div>

        <div class="form-group">
          <label>Group Photo</label>
          <input
            type="file"
            @change="handleGroupPhotoUpload"
            accept="image/*"
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label>Add User</label>
          <input
            v-model="userToAdd"
            type="text"
            placeholder="User ID"
            class="form-input"
          />
          <button @click="addUserToGroup" class="btn-primary">Add User</button>
        </div>

        <div class="modal-actions">
          <button @click="leaveGroup" class="btn-danger">Leave Group</button>
          <button @click="showGroupSettings = false" class="btn-secondary">Close</button>
        </div>
      </div>
    </div>

    <!-- Comment Modal -->
    <div v-if="showCommentModal" class="modal-overlay" @click="showCommentModal = false">
      <div class="modal" @click.stop>
        <h3>Add Comment</h3>
        <textarea
          v-model="commentText"
          placeholder="Your comment..."
          class="form-input"
          rows="4"
        ></textarea>
        <div class="modal-actions">
          <button @click="showCommentModal = false" class="btn-secondary">Cancel</button>
          <button @click="submitComment" class="btn-primary">Submit</button>
        </div>
      </div>
    </div>

    <!-- Forward Modal -->
    <div v-if="showForwardModal" class="modal-overlay" @click="showForwardModal = false">
      <div class="modal" @click.stop>
        <h3>Forward Message</h3>
        <select v-model="forwardToConversation" class="form-input">
          <option value="">Select conversation</option>
          <option v-for="conv in allConversations" :key="conv.id" :value="conv.id">
            {{ conv.name }}
          </option>
        </select>
        <div class="modal-actions">
          <button @click="showForwardModal = false" class="btn-secondary">Cancel</button>
          <button @click="submitForward" class="btn-primary">Forward</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiService from '../services/api.js'

const route = useRoute()
const router = useRouter()

const conversation = ref(null)
const messages = ref([])
const allConversations = ref([])
const loading = ref(false)
const newMessage = ref('')
const currentUsername = ref(localStorage.getItem('username') || '')

// Modal states
const showGroupSettings = ref(false)
const showCommentModal = ref(false)
const showForwardModal = ref(false)

// Form data
const groupName = ref('')
const userToAdd = ref('')
const commentText = ref('')
const forwardToConversation = ref('')
const messageToComment = ref(null)
const messageToForward = ref(null)

const messagesContainer = ref(null)

onMounted(async () => {
  await loadConversation()
  await loadAllConversations()
})

watch(() => route.params.id, async () => {
  await loadConversation()
})

async function loadConversation() {
  const conversationId = route.params.id
  if (!conversationId) return

  loading.value = true
  try {
    conversation.value = await apiService.getConversation(conversationId)
    messages.value = conversation.value.messages || []
    groupName.value = conversation.value.name
    
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load conversation:', error)
    alert('Failed to load conversation')
  } finally {
    loading.value = false
  }
}

async function loadAllConversations() {
  try {
    allConversations.value = await apiService.getMyConversations()
  } catch (error) {
    console.error('Failed to load conversations:', error)
  }
}

async function sendMessage() {
  if (!newMessage.value || !conversation.value) return

  try {
    const message = await apiService.sendMessage(conversation.value.id, newMessage.value)
    messages.value.push(message)
    newMessage.value = ''
    
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('Failed to send message:', error)
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
    await apiService.forwardMessage(messageToForward.value.id, forwardToConversation.value)
    showForwardModal.value = false
    forwardToConversation.value = ''
    messageToForward.value = null
    alert('Message forwarded successfully')
  } catch (error) {
    console.error('Failed to forward message:', error)
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
    await apiService.commentMessage(messageToComment.value.id, commentText.value)
    await loadConversation() // Reload to show comment
    showCommentModal.value = false
    commentText.value = ''
    messageToComment.value = null
  } catch (error) {
    console.error('Failed to add comment:', error)
    alert('Failed to add comment')
  }
}

async function deleteMessage(message) {
  if (!confirm('Are you sure you want to delete this message?')) return

  try {
    await apiService.deleteMessage(message.id)
    messages.value = messages.value.filter(m => m.id !== message.id)
  } catch (error) {
    console.error('Failed to delete message:', error)
    alert('Failed to delete message')
  }
}

async function updateGroupName() {
  if (!groupName.value || !conversation.value) return

  try {
    await apiService.setGroupName(conversation.value.id, groupName.value)
    conversation.value.name = groupName.value
    alert('Group name updated successfully')
  } catch (error) {
    console.error('Failed to update group name:', error)
    alert('Failed to update group name')
  }
}

async function handleGroupPhotoUpload(event) {
  const file = event.target.files[0]
  if (!file || !conversation.value) return

  try {
    await apiService.setGroupPhoto(conversation.value.id, file)
    alert('Group photo updated successfully')
  } catch (error) {
    console.error('Failed to update group photo:', error)
    alert('Failed to update group photo')
  }
}

async function addUserToGroup() {
  if (!userToAdd.value || !conversation.value) return

  try {
    await apiService.addToGroup(conversation.value.id, userToAdd.value)
    userToAdd.value = ''
    alert('User added to group successfully')
  } catch (error) {
    console.error('Failed to add user to group:', error)
    alert('Failed to add user to group')
  }
}

async function leaveGroup() {
  if (!conversation.value || !confirm('Are you sure you want to leave this group?')) return

  try {
    await apiService.leaveGroup(conversation.value.id)
    router.push('/conversations')
  } catch (error) {
    console.error('Failed to leave group:', error)
    alert('Failed to leave group')
  }
}

function goBack() {
  router.push('/conversations')
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function formatTime(timestamp) {
  return new Date(timestamp).toLocaleTimeString()
}
</script>

<style scoped>
.conversation-detail-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f8fafc;
}

.conversation-header {
  display: flex;
  align-items: center;
  padding: 1rem 2rem;
  background: white;
  border-bottom: 2px solid #e2e8f0;
  gap: 1rem;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: transparent;
  color: #667eea;
  border: 2px solid #667eea;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.back-button:hover {
  background: #667eea;
  color: white;
}

.conversation-info {
  flex: 1;
}

.conversation-info h1 {
  margin: 0 0 0.25rem 0;
  color: #2d3748;
  font-size: 1.5rem;
}

.conversation-type {
  background: #e2e8f0;
  color: #4a5568;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  text-transform: capitalize;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.btn-secondary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #667eea;
  color: white;
}

.conversation-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
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

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 2rem;
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
  background: white;
  padding: 1rem 1.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.own-message .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.message-body {
  margin-bottom: 0.5rem;
  line-height: 1.5;
}

.message-meta {
  font-size: 0.875rem;
  opacity: 0.7;
  display: flex;
  gap: 1rem;
}

.message-comments {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.comment {
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  opacity: 0.8;
}

.comment-author {
  font-weight: 600;
}

.message-actions {
  margin-top: 0.75rem;
  display: flex;
  gap: 0.5rem;
}

.btn-tiny {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 0.25rem;
  color: inherit;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-tiny:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-tiny.danger {
  border-color: rgba(255, 68, 68, 0.5);
  color: #ff4444;
}

.btn-tiny.danger:hover {
  background: rgba(255, 68, 68, 0.1);
}

.message-input-container {
  padding: 1rem 2rem;
  background: white;
  border-top: 2px solid #e2e8f0;
}

.message-input {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.input {
  flex: 1;
  padding: 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
}

.input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.send-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.btn-primary {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 0.5rem;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-danger {
  padding: 0.75rem 1.5rem;
  background: #e53e3e;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-danger:hover {
  background: #c53030;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}
</style>
