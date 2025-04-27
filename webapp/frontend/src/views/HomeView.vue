<template>
  <main>
    <Hero>
      <div>
        <Post :content="'# Cool\n**Nice**'" :username="authStore.username"></Post>
      </div>
    </Hero>
  </main>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import Hero from '@/components/ui/Hero.vue'
import Post from '@/components/ui/Post.vue'

const authStore = useAuthStore()
</script>
