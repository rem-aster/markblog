<template>
  <div class="inline-flex gap-2">
    <div class="tooltip">
      <div class="tooltip-content">
        <a class="font-black">{{ displayUsername }}</a>
      </div>
      <div class="btn btn-circle btn-ghost shadow-xl">
        <svg class="size-8" v-identicon="props.username"></svg>
      </div>
    </div>
    <div class="card card-border bg-base-100 shadow-xl">
      <div class="card-body text-wrap prose" v-html="content"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { computed } from 'vue'

const props = defineProps<{
  content?: string
  username?: string
}>()

const content = computed(() => {
  return DOMPurify.sanitize(marked.parse(props.content || '', { async: false }))
})

const displayUsername = computed(() => '@' + (props.username || 'anonymous'))
</script>
