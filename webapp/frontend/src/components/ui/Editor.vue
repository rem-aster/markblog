<template>
  <div class="grid grid-cols-1 gap-2">
    <Post :content="editorContent"></Post>
    <textarea v-model="editorContent" class="textarea" placeholder="# Hello"></textarea>
  </div>
</template>

<script setup lang="ts">
import { defineModel } from 'vue'
import Post from './Post.vue'

const editorContent = defineModel<string>()
</script>
