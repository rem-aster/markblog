import { createApp, h } from 'vue'
import Alert from '@/components/ui/Alert.vue'

type AlertType = 'info' | 'success' | 'warning' | 'error'

interface AlertOptions {
  type?: AlertType
  message: string
  duration?: number
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left'
  closable?: boolean
}

export const useAlert = () => {
  const showAlert = (options: AlertOptions) => {
    const alertDiv = document.createElement('div')
    document.body.appendChild(alertDiv)

    const app = createApp({
      render() {
        return h(Alert, {
          type: options.type,
          message: options.message,
          duration: options.duration,
          position: options.position,
          closable: options.closable,
          onClose: () => {
            app.unmount()
            document.body.removeChild(alertDiv)
          },
        })
      },
    })

    app.mount(alertDiv)
  }

  return {
    info: (message: string, options?: Omit<AlertOptions, 'type' | 'message'>) =>
      showAlert({ type: 'info', message, ...options }),
    success: (message: string, options?: Omit<AlertOptions, 'type' | 'message'>) =>
      showAlert({ type: 'success', message, ...options }),
    warning: (message: string, options?: Omit<AlertOptions, 'type' | 'message'>) =>
      showAlert({ type: 'warning', message, ...options }),
    error: (message: string, options?: Omit<AlertOptions, 'type' | 'message'>) =>
      showAlert({ type: 'error', message, ...options }),
  }
}
