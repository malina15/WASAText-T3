<template>
  <div class="profile-container">
    <div class="profile-header">
      <button @click="goBack" class="back-button">
        <span class="icon">←</span>
        Back
      </button>
      <h1>Profile Settings</h1>
    </div>

    <div class="profile-content">
      <div class="profile-card">
        <div class="profile-section">
          <h3>Personal Information</h3>
          
          <div class="form-group">
            <label for="username">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              placeholder="Enter your username"
              class="form-input"
              :class="{ error: usernameError }"
              minlength="3"
              maxlength="16"
              pattern="^[a-zA-Z0-9_]{3,16}$"
            />
            <div v-if="usernameError" class="error-message">{{ usernameError }}</div>
            <button 
              @click="updateUsername" 
              class="btn-primary"
              :disabled="!username || username === originalUsername || updatingUsername"
            >
              <span v-if="updatingUsername">Updating...</span>
              <span v-else>Update Username</span>
            </button>
          </div>

          <div class="form-group">
            <label for="photo">Profile Photo</label>
            <div class="photo-upload">
              <input
                id="photo"
                type="file"
                @change="handlePhotoUpload"
                accept="image/*"
                class="file-input"
              />
              <label for="photo" class="file-label">
                <span class="icon">📷</span>
                Choose Photo
              </label>
              <div v-if="photoUploading" class="upload-status">
                <div class="spinner"></div>
                <span>Uploading...</span>
              </div>
            </div>
            <p class="help-text">Upload a JPEG, PNG, or GIF image (max 10MB)</p>
          </div>
        </div>

        <div class="profile-section">
          <h3>Account Information</h3>
          <div class="info-item">
            <label>Current Username:</label>
            <span>{{ originalUsername }}</span>
          </div>
          <div class="info-item">
            <label>Account Status:</label>
            <span class="status-active">Active</span>
          </div>
        </div>

        <div class="profile-section">
          <h3>Actions</h3>
          <div class="action-buttons">
            <button @click="logout" class="btn-danger">
              <span class="icon">🚪</span>
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()

const username = ref('')
const originalUsername = ref('')
const usernameError = ref('')
const updatingUsername = ref(false)
const photoUploading = ref(false)

onMounted(() => {
  originalUsername.value = localStorage.getItem('username') || ''
  username.value = originalUsername.value
})

async function updateUsername() {
  if (!username.value || username.value === originalUsername.value) return
  
  updatingUsername.value = true
  usernameError.value = ''
  
  try {
    await apiService.setMyUserName(username.value)
    originalUsername.value = username.value
    localStorage.setItem('username', username.value)
    alert('Username updated successfully!')
  } catch (error) {
    console.error('Failed to update username:', error)
    usernameError.value = error.response?.data?.error || 'Failed to update username'
  } finally {
    updatingUsername.value = false
  }
}

async function handlePhotoUpload(event) {
  const file = event.target.files[0]
  if (!file) return
  
  // Validate file size (10MB max)
  if (file.size > 10 * 1024 * 1024) {
    alert('File size must be less than 10MB')
    return
  }
  
  // Validate file type
  const validTypes = ['image/jpeg', 'image/png', 'image/gif']
  if (!validTypes.includes(file.type)) {
    alert('Please upload a JPEG, PNG, or GIF image')
    return
  }
  
  photoUploading.value = true
  
  try {
    await apiService.setMyPhoto(file)
    alert('Profile photo updated successfully!')
  } catch (error) {
    console.error('Failed to update photo:', error)
    alert('Failed to update profile photo')
  } finally {
    photoUploading.value = false
    // Reset file input
    event.target.value = ''
  }
}

function goBack() {
  router.push('/conversations')
}

function logout() {
  if (confirm('Are you sure you want to logout?')) {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/')
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  min-height: 100vh;
  background: #f8fafc;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e2e8f0;
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

.profile-header h1 {
  margin: 0;
  color: #2d3748;
  font-size: 2rem;
  font-weight: 700;
}

.profile-content {
  display: flex;
  justify-content: center;
}

.profile-card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 600px;
}

.profile-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #e2e8f0;
}

.profile-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.profile-section h3 {
  margin: 0 0 1.5rem 0;
  color: #2d3748;
  font-size: 1.25rem;
  font-weight: 600;
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

.form-input.error {
  border-color: #e53e3e;
}

.error-message {
  color: #e53e3e;
  font-size: 0.875rem;
  margin-top: 0.25rem;
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

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.photo-upload {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.file-input {
  display: none;
}

.file-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.875rem 1rem;
  background: #f7fafc;
  border: 2px dashed #cbd5e0;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  justify-content: center;
}

.file-label:hover {
  background: #edf2f7;
  border-color: #667eea;
}

.upload-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #667eea;
  font-size: 0.875rem;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e2e8f0;
  border-top: 2px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.help-text {
  color: #718096;
  font-size: 0.875rem;
  margin: 0.5rem 0 0 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f7fafc;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: 500;
  color: #4a5568;
}

.info-item span {
  color: #2d3748;
}

.status-active {
  background: #c6f6d5;
  color: #22543d;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.action-buttons {
  display: flex;
  gap: 1rem;
}

.btn-danger {
  display: flex;
  align-items: center;
  gap: 0.5rem;
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
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(229, 62, 62, 0.3);
}

.icon {
  font-size: 1.2rem;
}
</style>
