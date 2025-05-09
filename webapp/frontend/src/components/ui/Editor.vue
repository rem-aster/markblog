<template>
  <div class="grid grid-cols-1 gap-2 place-content-center">
    <h3 class="font-semibold text-lg">Preview</h3>
    <PostBody :content="editorContent"></PostBody>
    <br />
    <form @submit.prevent="handlePostSubmit">
      <fieldset class="fieldset">
        <label class="label">Content</label>
        <textarea
          v-model="editorContent"
          class="textarea"
          placeholder="# Start typing here"
        ></textarea>
        <button
          class="btn btn-neutral mt-4"
          :class="{ 'btn-disabled': loading || !canPost }"
          :disabled="loading || !canPost"
        >
          Post
        </button>
      </fieldset>
    </form>
  </div>
</template>

<script setup lang="ts">
import { defineModel } from 'vue'
import PostBody from './PostBody.vue'
import { ref, computed } from 'vue'
import Client, { Local } from '../../client'
import { useAlert } from '@/services/alert'

const client = new Client(Local, {
  requestInit: {
    credentials: 'include',
  },
})

const emit = defineEmits(['postSubmit'])

const loading = ref(false)
const alert = useAlert()

const editorContent = defineModel<string>()

const canPost = computed(() => {
  return editorContent.value && editorContent.value.trim().length > 0
})

async function handlePostSubmit() {
  if (!canPost.value || loading.value) return

  loading.value = true
  try {
    const response = await client.webapp.Post(
      'POST',
      JSON.stringify({
        content: editorContent.value,
      }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      },
    )

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    editorContent.value = ''
    emit('postSubmit')
  } catch (err) {
    console.error('Error creating post:', err)
    alert.error('Post creation failed: ' + err, {
      position: 'bottom-right',
      closable: true,
    })
  } finally {
    loading.value = false
  }
}
</script>
