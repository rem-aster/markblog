import { defineStore } from 'pinia'
import { ref } from 'vue'
import Client, { Local } from '../client'

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(false)
  const isLoading = ref(false)
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
    } catch (error) {
      isAuthenticated.value = false
    } finally {
      isLoading.value = false
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
      checkAuth()
      return true
    } catch (error) {
      console.error('Login failed:', error)
      // Handle error (show message to user, etc.)
    } finally {
      isLoading.value = false
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
      checkAuth()
      return true
    } catch (error) {
      console.error('Registration failed:', error)
      // Handle error
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    isLoading.value = true
    try {
      isAuthenticated.value = false
      const response = await client.webapp.Logout('GET')
      const data = await response.json()
      console.log('Login successful:', data)
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      isLoading.value = false
    }
  }

  return {
    isAuthenticated,
    checkAuth,
    login,
    register,
    logout,
  }
})
