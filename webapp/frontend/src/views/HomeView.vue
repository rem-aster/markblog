<template>
  <main>
    <Hero>
      <div>
        <button
          @click="authStore.checkAuth()"
          class="btn"
          :class="{ 'btn-primary': authStore.isAuthenticated }"
        >
          Lol
        </button>
      </div>
    </Hero>
  </main>
</template>

<script setup lang="ts">
import Hero from '@/components/ui/Hero.vue'
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
</script>
