<template>
  <li class="list-row">
    <div class="font-semibold">
      {{ activityType }}
    </div>
    <div>
      {{ props.content }}
    </div>
  </li>
</template>

<script setup lang="ts">
const props = defineProps<{
  content?: string
  activityType?: string
  activityTime?: string
}>()
</script>
