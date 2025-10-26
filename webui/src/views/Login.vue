<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="app-title">WASA Chat</h1>
      <p class="app-subtitle">Connect and communicate with your team</p>
      
      <form @submit.prevent="login" class="login-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Enter your username"
            class="form-input"
            :class="{ error: error }"
            required
            minlength="3"
            maxlength="16"
            pattern="^[a-zA-Z0-9_]{3,16}$"
          />
          <div v-if="error" class="error-message">{{ error }}</div>
        </div>
        
        <button 
          type="submit" 
          class="login-button"
          :disabled="loading || !username"
        >
          <span v-if="loading">Signing in...</span>
          <span v-else>Sign In</span>
        </button>
      </form>
      
      <div class="login-info">
        <p>Don't have an account? No problem!</p>
        <p>Just enter a username and we'll create one for you.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()
const username = ref('')
const loading = ref(false)
const error = ref('')

async function login() {
  if (!username.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await apiService.login(username.value)
    
    // Store token
    localStorage.setItem('token', response.identifier)
    localStorage.setItem('username', username.value)
    
    // Redirect to conversations
    router.push('/conversations')
  } catch (err) {
    console.error('Login error:', err)
    error.value = err.response?.data?.error || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.login-card {
  background: white;
  border-radius: 1rem;
  padding: 3rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.app-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 0.5rem 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.app-subtitle {
  color: #718096;
  margin: 0 0 2rem 0;
  font-size: 1.1rem;
}

.login-form {
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
  text-align: left;
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

.login-button {
  width: 100%;
  padding: 0.875rem 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
}

.login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-info {
  color: #718096;
  font-size: 0.875rem;
  line-height: 1.5;
}

.login-info p {
  margin: 0.5rem 0;
}
</style>
