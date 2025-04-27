import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { update } from 'jdenticon'

const pinia = createPinia()
const app = createApp(App)

app.use(router)
app.use(pinia)
app.directive('identicon', {
  mounted(el, binding) {
    update(el, binding.value)
  },
  updated(el, binding) {
    if (binding.value !== binding.oldValue) {
      update(el, binding.value)
    }
  },
})

app.mount('#app')
