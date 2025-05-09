<template>
  <main>
    <Hero>
      <div class="flex flex-col gap-2 content-center">
        <Post v-for="post in posts" :key="post.id" :content="post.content">
          <Profile :username="post.username" @click="openModalActivity(post.username)" />
          <button
            class="btn btn-xs btn-square btn-ghost shadow-xl"
            @click="openModalComments(post.id)"
          >
            <Messages class="text-base-content" />
          </button>
        </Post>

        <button v-if="!loading && hasMorePosts" @click="fetchPosts" class="btn btn-neutral btn-sm">
          Load More
        </button>

        <div v-if="loading" class="text-center py-4">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-if="!hasMorePosts" class="text-center text-base-content py-4">
          You've reached the end of the feed
        </div>
      </div>
    </Hero>
  </main>

  <Modal :open="openedModal == 'comments'">
    <button class="btn btn-xs btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      <Close class="text-base-content" />
    </button>
    <div class="flex flex-col">
      <h1 class="text-xl font-semibold">Comments</h1>
      <form @submit.prevent="handleCommentSubmit" class="inline-flex gap-2 my-2">
        <input type="text" placeholder="Comment" class="input input-sm" v-model="inputContent" />
        <button
          class="btn btn-primary btn-sm btn-square"
          :class="{ 'btn-disabled': loading || !canComment }"
          :disabled="loading || !canComment"
        >
          <Right />
        </button>
      </form>

      <div class="overflow-y-auto max-h-76 flex flex-col">
        <Comment
          v-for="comment in comments"
          :key="comment.id"
          :username="comment.username"
          :content="comment.content"
        />

        <button
          v-if="!loading && hasMoreComments"
          @click="fetchComments"
          class="btn btn-neutral btn-sm mt-2"
        >
          Load More
        </button>

        <div v-if="loading" class="text-center py-4 mt-2">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div
          v-if="!hasMoreComments"
          class="text-center text-sm text-base-content font-semibold py-4 mt-2"
        >
          No more comments
        </div>
      </div>
    </div>
  </Modal>

  <Modal :open="openedModal == 'activity'">
    <button class="btn btn-xs btn-circle btn-ghost absolute right-2 top-2" @click="closeModal()">
      <Close class="text-base-content" />
    </button>
    <div class="flex flex-col">
      <h1 class="text-xl font-semibold">@{{ activityUsername }}'s Activity</h1>

      <div class="overflow-y-auto max-h-76 gap-2 flex flex-col">
        <ul class="list">
          <Activity
            v-for="act in activity"
            :key="act.postId"
            :content="act.content"
            :activityType="act.activityType"
            :activityTime="act.activityTime"
          />
        </ul>
        <button
          v-if="!loading && hasMoreActivity"
          @click="fetchActivity"
          class="btn btn-neutral btn-sm mt-2"
        >
          Load More
        </button>

        <div v-if="loading" class="text-center py-4 mt-2">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div
          v-if="!hasMoreActivity"
          class="text-center text-sm text-base-content font-semibold py-4 mt-2"
        >
          No more activity
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import Hero from '@/components/ui/Hero.vue'
import Post from '@/components/ui/Post.vue'
import Modal from '@/components/ui/Modal.vue'
import Close from '@/components/icons/Close.vue'
import Messages from '@/components/icons/Messages.vue'
import Right from '@/components/icons/Right.vue'
import { ref, onMounted, computed } from 'vue'
import Client, { Local } from '../client'
import Comment from '@/components/ui/Comment.vue'
import { useAlert } from '@/services/alert'
import { useAuthStore } from '@/stores/auth'
import Profile from '@/components/ui/Profile.vue'
import Activity from '@/components/ui/Activity.vue'

const authStore = useAuthStore()

interface PostType {
  id: string
  username: string
  content: string
  createdAt?: string
  ID?: string
  Content?: string
  Username?: string
  CreatedAt?: string
}

interface CommentType {
  id: string
  username: string
  content: string
  createdAt?: string
  ID?: string
  Content?: string
  Username?: string
  CreatedAt?: string
}

interface ActivityType {
  postId: string
  content: string
  actionTime?: string
  actionType: string
  PostID?: string
  Content?: string
  ActionTime?: string
  ActionType: string
}

const alert = useAlert()

const client = new Client(Local, {
  requestInit: {
    credentials: 'include',
  },
})

const openedModal = defineModel<string>('openedModal')

function closeModal() {
  commentsPost.value = ''
  activityUsername.value = ''
  comments.value = []
  activity.value = []
  commentOffset.value = 0
  activityOffset.value = 0
  hasMoreComments.value = true
  hasMoreActivity.value = true
  openedModal.value = ''
}

const inputContent = defineModel<string>('inputContent')

const canComment = computed(() => {
  return inputContent.value && inputContent.value.trim().length > 0
})

const commentsPost = ref('')
const activityUsername = ref('')
const posts = ref<Array<{ id: string; username: string; content: string }>>([])
const comments = ref<Array<{ id: string; username: string; content: string }>>([])
const activity = ref<
  Array<{ postId: string; activityTime: string; content: string; activityType: string }>
>([])
const loading = ref(false)
const hasMorePosts = ref(true)
const hasMoreComments = ref(true)
const hasMoreActivity = ref(true)
const postOffset = ref(0)
const commentOffset = ref(0)
const activityOffset = ref(0)
const limit = 5

async function openModalComments(postID: string) {
  commentsPost.value = postID
  await fetchComments()
  openedModal.value = 'comments'
}

async function openModalActivity(username: string) {
  activityUsername.value = username
  await fetchActivity()
  openedModal.value = 'activity'
}

async function fetchPosts() {
  if (loading.value || !hasMorePosts.value) return

  loading.value = true
  try {
    const response = await client.webapp.Feed(
      'POST',
      JSON.stringify({
        offset: postOffset.value,
        limit: limit,
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

    const data = await response.json()
    console.log('API Response:', data)
    let newPosts = []
    if (Array.isArray(data)) {
      newPosts = data
    } else if (data?.posts && Array.isArray(data.posts)) {
      newPosts = data.posts
    }
    const transformedPosts = newPosts.map((post: PostType) => ({
      id: post.ID || post.id,
      username: post.Username || post.username,
      content: post.Content || post.content,
    }))

    if (transformedPosts.length === 0) {
      hasMorePosts.value = false
    } else {
      posts.value = [...posts.value, ...transformedPosts]
      postOffset.value += transformedPosts.length
      if (transformedPosts.length < limit) {
        hasMorePosts.value = false
      }
    }
  } catch (error) {
    console.error('Error fetching posts:', error)
    alert.error('Fetching posts failed: ' + error, {
      position: 'bottom-right',
      closable: true,
    })
    hasMorePosts.value = false
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  if (loading.value || !hasMoreComments.value) return

  loading.value = true
  try {
    const response = await client.webapp.Discussion(
      'POST',
      JSON.stringify({
        id: commentsPost.value,
        offset: commentOffset.value,
        limit: limit,
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

    const data = await response.json()
    console.log('API Response:', data)
    let newComments = []
    if (Array.isArray(data)) {
      newComments = data
    } else if (data?.comments && Array.isArray(data.comments)) {
      newComments = data.comments
    }
    const transformedComments = newComments.map((comment: CommentType) => ({
      id: comment.ID || comment.id,
      username: comment.Username || comment.username,
      content: comment.Content || comment.content,
    }))

    if (transformedComments.length === 0) {
      hasMoreComments.value = false
    } else {
      comments.value = [...comments.value, ...transformedComments]
      commentOffset.value += transformedComments.length
      if (transformedComments.length < limit) {
        hasMoreComments.value = false
      }
    }
  } catch (error) {
    console.error('Error fetching comments:', error)
    alert.error('Fetching comments failed: ' + error, {
      position: 'bottom-right',
      closable: true,
    })
    hasMoreComments.value = false
  } finally {
    loading.value = false
  }
}

async function fetchActivity() {
  if (loading.value || !hasMoreActivity.value) return

  loading.value = true
  try {
    const response = await client.webapp.Activity(
      'POST',
      JSON.stringify({
        username: activityUsername.value,
        offset: activityOffset.value,
        limit: limit,
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

    const data = await response.json()
    console.log('API Response:', data)
    let newActivity = []
    if (Array.isArray(data)) {
      newActivity = data
    } else if (data?.activity && Array.isArray(data.activity)) {
      newActivity = data.activity
    }
    const transformedActivity = newActivity.map((activity: ActivityType) => ({
      postId: activity.PostID || activity.postId,
      actionType: activity.ActionType || activity.actionType,
      actionTime: activity.ActionTime || activity.actionTime,
      content: activity.Content || activity.content,
    }))

    if (transformedActivity.length === 0) {
      hasMoreActivity.value = false
    } else {
      activity.value = [...activity.value, ...transformedActivity]
      activityOffset.value += transformedActivity.length
      if (transformedActivity.length < limit) {
        hasMoreActivity.value = false
      }
    }
  } catch (error) {
    console.error('Error fetching activity:', error)
    alert.error('Fetching activity failed: ' + error, {
      position: 'bottom-right',
      closable: true,
    })
    hasMoreActivity.value = false
  } finally {
    loading.value = false
  }
}

async function handleCommentSubmit() {
  if (!canComment.value || loading.value || !inputContent.value) return

  loading.value = true
  try {
    const response = await client.webapp.Comment(
      'POST',
      JSON.stringify({
        content: inputContent.value,
        post_id: commentsPost.value,
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

    const data = await response.json()

    comments.value.unshift({
      id: data.ID || data.id,
      username: authStore.username || 'anonymous',
      content: inputContent.value,
    })

    inputContent.value = ''
    await fetchComments()
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

onMounted(() => {
  fetchPosts()
})
</script>
