<template>
  <div class="flex flex-row gap-2">
    <div class="flex flex-col gap-2 items-end">
      <slot> </slot>
    </div>
    <PostBody :content="content" />
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { computed } from 'vue'
import PostBody from './PostBody.vue'

const props = defineProps<{
  content?: string
}>()

const content = computed(() => {
  return DOMPurify.sanitize(marked.parse(props.content || '', { async: false, breaks: true }))
})
</script>
