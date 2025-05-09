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
          <a class="btn btn-square btn-primary" @click="openModalPost()">
            <Pen class="text-primary-content" />
          </a>
        </div>
        <div class="tooltip tooltip-bottom">
          <div class="tooltip-content">
            <div class="font-black">Settings</div>
          </div>
          <a class="btn btn-neutral font-black btn-square">
            <Cog class="text-neutral-content" />
          </a>
        </div>
        <div class="tooltip tooltip-bottom">
          <div class="tooltip-content">
            <div class="font-black">Logout</div>
          </div>
          <a class="btn btn-neutral font-black btn-square" @click="authStore.logout()">
            <Exit class="text-neutral-content" />
          </a>
        </div>
      </div>

      <div class="inline-flex gap-2" v-else>
        <a class="btn btn-primary font-black" @click="openModalLogin()">Login</a>
        <a class="btn btn-neutral font-black" @click="openModalRegister()">Register</a>
      </div>
    </template>
  </Navbar>

  <Modal :open="openedModal == 'login'">
    <button class="btn btn-xs btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      <Close class="text-base-content" />
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

        <button class="btn btn-primary mt-4" :class="{ 'btn-disabled': authStore.isLoading }">
          Login
        </button>

        <a class="link link-hover place-end-end" @click="openModalRegister()">Register instead</a>
      </fieldset>
    </form>
  </Modal>

  <Modal :open="openedModal == 'register'">
    <button class="btn btn-xs btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      <Close class="text-base-content" />
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

        <button class="btn btn-neutral mt-4" :class="{ 'btn-disabled': authStore.isLoading }">
          Register
        </button>

        <a class="link link-hover place-end-end" @click="openModalLogin()">Login instead</a>
      </fieldset>
    </form>
  </Modal>

  <Modal :open="openedModal == 'post'" class="">
    <button class="btn btn-xs btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      <Close class="text-base-content" />
    </button>
    <Editor @post-submit="handlePostSubmit" />
  </Modal>

  <HomeView v-model="openedModal" @comment-submit="handleCommentSubmit" />
</template>

<script setup lang="ts">
import Navbar from './components/ui/Navbar.vue'
import Modal from './components/ui/Modal.vue'
import Pen from './components/icons/Pen.vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { ref, onMounted } from 'vue'
import Exit from './components/icons/Exit.vue'
import Close from './components/icons/Close.vue'
import Editor from './components/ui/Editor.vue'
import Cog from './components/icons/Cog.vue'
import HomeView from './views/HomeView.vue'

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const openedModal = ref('')

function openModalRegister() {
  openedModal.value = 'register'
}

function openModalLogin() {
  openedModal.value = 'login'
}

function openModalPost() {
  openedModal.value = 'post'
}

function closeModal() {
  openedModal.value = ''
}

const username = ref('')
const password = ref('')

function handlePostSubmit() {
  closeModal()
  window.location.reload()
}

function handleCommentSubmit() {
  window.location.reload()
}

async function handleLogin() {
  if (
    (await authStore.login({ username: username.value, password: password.value })) &&
    authStore.isAuthenticated
  ) {
    closeModal()
  }
}

async function handleRegister() {
  if (
    (await authStore.register({ username: username.value, password: password.value })) &&
    authStore.isAuthenticated
  ) {
    closeModal()
  }
}

onMounted(() => {
  authStore.checkAuth()
})
</script>
