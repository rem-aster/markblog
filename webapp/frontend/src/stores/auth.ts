import { defineStore } from 'pinia'
import { ref } from 'vue'
import Client, { Local } from '../client'
import { useAlert } from '@/services/alert'

const alert = useAlert()

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(false)
  const isLoading = ref(false)
  const username = ref('')
  const client = new Client(Local, {
    requestInit: {
      credentials: 'include',
    },
  })

  async function checkAuth() {
    isLoading.value = true
    try {
      const response = await client.webapp.CheckAuth('GET')
      const data = await response.json()
      isAuthenticated.value = data['authenticated']
      username.value = data['user']['username']
    } catch (error) {
      console.error('Auth check failed:', error)
      isAuthenticated.value = false
    } finally {
      isLoading.value = false
      return isAuthenticated
    }
  }

  async function login(credentials: { username: string; password: string }) {
    isLoading.value = true
    try {
      const response = await client.webapp.Login(
        'POST',
        JSON.stringify({
          username: credentials.username,
          password: credentials.password,
        }),
        {
          headers: {
            'Content-Type': 'application/json',
          },
        },
      )

      const data = await response.json()
      console.log('Login successful:', data)
    } catch (error) {
      console.error('Login failed:', error)
      alert.error('Login failed: ' + error, {
        position: 'bottom-right',
        closable: true,
      })
    } finally {
      isLoading.value = false
      return checkAuth()
    }
  }

  async function register(credentials: { username: string; password: string }) {
    isLoading.value = true
    try {
      const response = await client.webapp.Register(
        'POST',
        JSON.stringify({
          username: credentials.username,
          password: credentials.password,
        }),
        {
          headers: {
            'Content-Type': 'application/json',
          },
        },
      )

      const data = await response.json()
      console.log('Registration successful:', data)
    } catch (error) {
      console.error('Registration failed:', error)
      alert.error('Registration failed: ' + error, {
        position: 'bottom-right',
        closable: true,
      })
    } finally {
      isLoading.value = false
      return checkAuth()
    }
  }

  async function logout() {
    isLoading.value = true
    try {
      isAuthenticated.value = false
      const response = await client.webapp.Logout('GET')
      const data = await response.json()
      console.log('Logout successful:', data)
    } catch (error) {
      console.error('Logout failed:', error)
      alert.error('Logout failed: ' + error, {
        position: 'bottom-right',
        closable: true,
      })
    } finally {
      isLoading.value = false
    }
  }

  return {
    isAuthenticated,
    username,
    isLoading,
    checkAuth,
    login,
    register,
    logout,
  }
})
