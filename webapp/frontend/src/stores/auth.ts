// stores/auth.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import Client from '../client'

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(false)
  const client = new Client('http://localhost:4000') // or your production URL

  // Check auth status on app load
  async function checkAuth() {
    try {
      // Replace with your actual auth check endpoint
      const response = await client.webapp.Login('GET')
      isAuthenticated.value = response.ok
    } catch (error) {
      isAuthenticated.value = false
    }
  }

  // Login function
  async function login(credentials: { username: string; password: string }) {
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
      return data
    } catch (error) {
      console.error('Login failed:', error)
      // Handle error (show message to user, etc.)
    }
  }

  // Register function
  async function register(credentials: { username: string; password: string }) {
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
      return data
    } catch (error) {
      console.error('Registration failed:', error)
      // Handle error
    }
  }

  // Logout function
  function logout() {
    isAuthenticated.value = false
    // Add any additional logout logic here
  }

  return {
    isAuthenticated,
    checkAuth,
    login,
    register,
    logout,
  }
})
