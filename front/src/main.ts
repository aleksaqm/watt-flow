import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import { useUserStore } from './stores/user'
import { getRoleFromToken, getUserIdFromToken } from './utils/jwtDecoder'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.config.globalProperties.$axios = axios
axios.defaults.baseURL = 'http://localhost:8080'
axios.interceptors.request.use((config) => {
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // Handle unauthorized access, e.g., redirect to login
      console.error('Unauthorized access - redirecting to login')
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
const initializeUserStore = () => {
  const userStore = useUserStore();
  userStore.setRole(getRoleFromToken())
  userStore.setId(getUserIdFromToken())
};

initializeUserStore();

app.mount('#app')
