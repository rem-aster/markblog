<template>
  <div class="card bg-base-100 shadow-xl w-50">
    <div
      class="card-body text-wrap prose break-words whitespace-normal"
      :class="{ 'text-base-300': content === '' }"
      v-html="content === '' ? 'Empty' : content"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { computed } from 'vue'

const props = defineProps<{
  content?: string
}>()

const content = computed(() => {
  return DOMPurify.sanitize(marked.parse(props.content || '', { async: false, breaks: true }))
})
</script>
