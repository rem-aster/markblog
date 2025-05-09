<template>
  <div>
    <div class="chat chat-start">
      <Profile :username="props.username" class="chat-image" />
      <div class="chat-bubble text-wrap truncate max-w-40">
        <p>{{ props.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Profile from './Profile.vue'

const props = defineProps<{
  content?: string
  username?: string
}>()
</script>
