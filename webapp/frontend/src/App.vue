<template>
  <Navbar>
    <template #start>
      <a class="text-lg font-black underline decoration-3" href="/">markblog</a>
    </template>
    <template #end>
      <div class="inline-flex gap-2" v-if="isAuthenticated">
        <div class="tooltip tooltip-bottom tooltip-primary">
          <div class="tooltip-content">
            <div class="font-black">New Post</div>
          </div>
          <a class="btn btn-square btn-primary">
            <Pen class="size-[24px] text-primary-content" />
          </a>
        </div>
        <div class="tooltip tooltip-bottom">
          <div class="tooltip-content">
            <div class="font-black">Logout</div>
          </div>
          <a class="btn btn-ghost font-black btn-square" @click="authStore.logout()">
            <Exit class="size-[24px] text-neutral" />
          </a>
        </div>
      </div>
      <div class="inline-flex gap-2" v-else>
        <a class="btn btn-primary transition ease-in-out hover:scale-110" @click="openModalLogin()"
          >Login</a
        >
        <a
          class="btn btn-neutral transition ease-in-out hover:scale-110"
          @click="openModalRegister()"
          >Register</a
        >
      </div>
    </template>
  </Navbar>

  <Modal :open="openedModal == 'login'">
    <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      ✕
    </button>
    <h1 class="font-black text-primary text-xl mb-4">Login</h1>
    <form @submit.prevent="handleLogin">
      <fieldset class="fieldset">
        <label class="label">Username</label>
        <input
          v-model="username"
          type="input"
          class="input validator"
          pattern="[A-Za-z][A-Za-z0-9\-]*"
          minlength="3"
          maxlength="30"
          title="Only letters, numbers or dash"
          required
          placeholder="Aboba123"
        />
        <label class="label">Password</label>
        <input
          v-model="password"
          type="password"
          class="input validator"
          required
          minlength="8"
          placeholder="Qwerty123"
        />

        <button class="btn btn-primary mt-4">Login</button>

        <a class="link link-hover place-end-end" @click="openModalRegister()">Register instead</a>
      </fieldset>
    </form>
  </Modal>

  <Modal :open="openedModal == 'register'">
    <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      ✕
    </button>
    <h1 class="font-black text-xl mb-4">Register</h1>
    <form @submit.prevent="handleRegister">
      <fieldset class="fieldset">
        <label class="label">Username</label>
        <input
          v-model="username"
          type="input"
          class="input validator"
          pattern="[A-Za-z][A-Za-z0-9\-]*"
          minlength="3"
          maxlength="30"
          title="Only letters, numbers or dash"
          required
          placeholder="Aboba123"
        />
        <label class="label">Password</label>
        <input
          v-model="password"
          type="password"
          class="input validator"
          required
          minlength="8"
          placeholder="Qwerty123"
        />

        <button class="btn btn-neutral mt-4">Register</button>

        <a class="link link-hover place-end-end" @click="openModalLogin()">Login instead</a>
      </fieldset>
    </form>
  </Modal>

  <RouterView />
</template>

<script setup lang="ts">
import Navbar from './components/ui/Navbar.vue'
import Modal from './components/ui/Modal.vue'
import { RouterView } from 'vue-router'
import Pen from './components/icons/Pen.vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { ref, onMounted } from 'vue'
import Exit from './components/icons/Exit.vue'

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)
const { logout } = authStore

const openedModal = ref('')

function openModalRegister() {
  openedModal.value = 'register'
}

function openModalLogin() {
  openedModal.value = 'login'
}

function closeModal() {
  openedModal.value = ''
}

const username = ref('')
const password = ref('')

async function handleLogin() {
  if (await authStore.login({ username: username.value, password: password.value })) {
  }
}

async function handleRegister() {
  if (await authStore.register({ username: username.value, password: password.value })) {
  }
}

onMounted(() => {
  authStore.checkAuth()
})
</script>
