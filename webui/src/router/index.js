import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Conversations from '../views/Conversations.vue'
import ConversationDetail from '../views/ConversationDetail.vue'
import Profile from '../views/Profile.vue'

const routes = [
  {
    path: '/',
    name: 'Login',
    component: Login
  },
  {
    path: '/conversations',
    name: 'Conversations',
    component: Conversations,
    meta: { requiresAuth: true }
  },
  {
    path: '/conversations/:id',
    name: 'ConversationDetail',
    component: ConversationDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    next('/')
  } else if (to.path === '/' && token) {
    next('/conversations')
  } else {
    next()
  }
})

export default router
