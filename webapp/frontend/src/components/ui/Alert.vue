<template>
  <!-- Template remains the same -->
  <transition name="fade">
    <div
      v-if="show"
      role="alert"
      :class="['alert fixed z-50 transition-all duration-300', alertClass, positionClass]"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6 shrink-0 stroke-current"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path :d="iconPath" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
      </svg>
      <span>{{ message }}</span>
      <button v-if="closable" @click="close" class="btn btn-sm btn-ghost">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

type AlertType = 'info' | 'success' | 'warning' | 'error'
type AlertPosition = 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left'

const props = defineProps({
  type: {
    type: String as () => AlertType,
    default: 'info',
    validator: (value: string): value is AlertType =>
      ['info', 'success', 'warning', 'error'].includes(value),
  },
  message: {
    type: String,
    required: true,
  },
  duration: {
    type: Number,
    default: 5000,
  },
  position: {
    type: String as () => AlertPosition,
    default: 'top-right',
    validator: (value: string): value is AlertPosition =>
      ['top-right', 'top-left', 'bottom-right', 'bottom-left'].includes(value),
  },
  closable: {
    type: Boolean,
    default: true,
  },
})

const show = ref(false)

interface AlertIcons {
  error: string
  success: string
  warning: string
  info: string
  [key: string]: string // Index signature
}

const icons: AlertIcons = {
  error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
  success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
  warning:
    'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
  info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
}

const alertClass = computed(() => ({
  'alert-info': props.type === 'info',
  'alert-success': props.type === 'success',
  'alert-warning': props.type === 'warning',
  'alert-error': props.type === 'error',
}))

const positionClass = computed(() => ({
  'top-4 right-4': props.position === 'top-right',
  'top-4 left-4': props.position === 'top-left',
  'bottom-4 right-4': props.position === 'bottom-right',
  'bottom-4 left-4': props.position === 'bottom-left',
}))

const iconPath = computed(() => icons[props.type])

let timeout: ReturnType<typeof setTimeout>

function close() {
  show.value = false
}

onMounted(() => {
  show.value = true
  if (props.duration > 0) {
    timeout = setTimeout(() => {
      close()
    }, props.duration)
  }
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
