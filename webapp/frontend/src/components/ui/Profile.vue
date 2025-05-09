<template>
  <div class="tooltip">
    <div class="tooltip-content z-50">
      <a class="font-black">{{ displayUsername }}</a>
    </div>
    <div class="btn btn-circle btn-ghost shadow-xl">
      <svg class="size-8" v-identicon="props.username"></svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  username?: string
}>()

const displayUsername = computed(() => '@' + (props.username || 'anonymous'))
</script>
