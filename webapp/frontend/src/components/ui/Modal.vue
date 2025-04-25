<template>
  <dialog class="modal place-content-center" :class="{ 'modal-open': openModal }">
    <div class="modal-box bg-base-100">
      <slot></slot>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps({
  open: Boolean,
})

const openModal = ref(false)

watch(
  () => props.open,
  (newVal) => {
    if (newVal !== undefined) {
      openModal.value = newVal
    }
  },
)
</script>
