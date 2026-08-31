<template>
  <Teleport to="body">
    <div class="chatbox-wrapper" ref="wrapperRef" v-show="visible">
      
      <div class="chatbox glass-strong" :class="{ minimized: isMinimized }">
        
        <div class="chat-header" @click="toggleMinimize">
          <div class="chat-header-info">
            <div class="chat-avatar">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </div>
            <div class="chat-title-block">
              <h3 class="chat-title">OpenPeo Chat</h3>
              <span class="chat-status" :class="{ online: isConnected }">
                <span class="status-dot"></span>
                {{ isConnected ? 'Terhubung' : 'Offline' }}
              </span>
            </div>
          </div>

          <div class="chat-header-actions">
            <button class="chat-header-btn" @click.stop="toggleMinimize" :title="isMinimized ? 'Buka' : 'Kecilkan'">
              <svg v-if="!isMinimized" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="4 14 10 14 10 20"/>
                <polyline points="20 10 14 10 14 4"/>
                <line x1="14" y1="10" x2="21" y2="3"/>
                <line x1="3" y1="21" x2="10" y2="14"/>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 3 21 3 21 9"/>
                <polyline points="9 21 3 21 3 15"/>
                <line x1="21" y1="3" x2="14" y2="10"/>
                <line x1="3" y1="21" x2="10" y2="14"/>
              </svg>
            </button>
            <button class="chat-header-btn close-btn" @click.stop="$emit('close')" title="Tutup">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>

        
        <div v-if="!isMinimized" class="chat-user-bar">
          <label class="user-label">Chat dengan:</label>
          <select v-model="receiverId" class="input user-select" @change="onReceiverChange">
            <option v-for="c in contacts" :key="c.id" :value="c.id">
              {{ c.name }} ({{ c.role.toUpperCase() }}) {{ unreadCounts[c.id] ? `(💬 ${unreadCounts[c.id]} baru)` : '' }}
            </option>
          </select>
        </div>

        
        <div v-if="!isMinimized" class="chat-messages" ref="messagesRef">
          
          <div v-if="messages.length === 0" class="chat-welcome">
            <div class="welcome-icon">💬</div>
            <p>Mulai percakapan dengan vendor!</p>
            <span>Tanyakan detail produk, harga, atau pengiriman.</span>
          </div>

          
          <div
            v-for="msg in messages"
            :key="msg.id || msg._tempId"
            class="message-row"
            :class="{ 'is-me': msg.sender_id === senderId }"
          >
            <div class="message-bubble" :class="msg.sender_id === senderId ? 'bubble-me' : 'bubble-other'">
              <p class="message-text">{{ msg.content }}</p>
              <span class="message-time">{{ formatTime(msg.created_at) }}</span>
            </div>
          </div>
        </div>

        
        <div v-if="!isMinimized" class="chat-input-area">
          <input
            v-model="newMessage"
            type="text"
            class="input chat-input"
            placeholder="Ketik pesan..."
            @keydown.enter.prevent="sendMessage"
            :disabled="!isConnected"
          />
          <button
            class="btn btn-primary btn-send"
            @click="sendMessage"
            :disabled="!newMessage.trim() || !isConnected"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="22" y1="2" x2="11" y2="13"/>
              <polygon points="22 2 15 22 11 13 2 9 22 2"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { animate } from 'animejs'
import axios from '../axios'
import { useAuthStore } from '../stores/auth'
import { useWebsocketStore } from '../stores/websocket'
import { useChatStore } from '../stores/chat'

const props = defineProps({
  visible: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close', 'update-unread'])

const authStore = useAuthStore()
const websocketStore = useWebsocketStore()
const chatStore = useChatStore()

const contacts = ref([])
const newMessage = ref('')
const unreadCounts = ref({})
const isMinimized = ref(false)
const messagesRef = ref(null)
const wrapperRef = ref(null)

const senderId = computed(() => authStore.user?.id)
const receiverId = ref(null)
const messages = computed(() => websocketStore.messages)
const isConnected = computed(() => websocketStore.isConnected)

const totalUnread = computed(() => {
  return Object.values(unreadCounts.value).reduce((sum, count) => sum + count, 0)
})

watch(totalUnread, (newVal) => {
  emit('update-unread', newVal)
}, { immediate: true })

// Re-fetch contacts when a new message arrives
watch(() => websocketStore.messages.length, (newVal, oldVal) => {
  if (newVal > oldVal) {
    loadContacts()
  }
})

async function loadContacts() {
  if (!authStore.user?.id) return
  try {
    const res = await axios.get('/api/chat/contacts', {
      params: { user_id: authStore.user.id, role: authStore.user.role }
    })
    if (res.data.success && res.data.data) {
      contacts.value = res.data.data
      contacts.value.forEach(c => {
        unreadCounts.value[c.id] = c.unread_count || 0
      })
      if (contacts.value.length > 0 && !receiverId.value) {
        receiverId.value = contacts.value[0].id
        chatStore.currentReceiverId = receiverId.value
      }
    }
  } catch (error) {
    console.warn('Could not load chat contacts:', error)
  }
}

async function loadChatHistory() {
  if (!senderId.value || !receiverId.value) return
  try {
    const res = await axios.get('/api/messages', {
      params: { sender_id: senderId.value, receiver_id: receiverId.value, limit: 50 }
    })
    if (res.data.success && res.data.data) {
      websocketStore.setMessages(res.data.data)
      scrollToBottom()
    }
  } catch (error) {
    console.warn('Could not load chat history:', error)
  }
}

function sendMessage() {
  if (!newMessage.value.trim() || !isConnected.value || !receiverId.value) return
  websocketStore.sendMessage({
    sender_id: senderId.value,
    receiver_id: receiverId.value,
    content: newMessage.value.trim()
  })
  newMessage.value = ''
  scrollToBottom()
}

function onReceiverChange() {
  chatStore.currentReceiverId = receiverId.value
  websocketStore.setMessages([])
  unreadCounts.value[receiverId.value] = 0
  loadChatHistory()
}

function toggleMinimize() {
  isMinimized.value = !isMinimized.value
  if (!isMinimized.value) {
    scrollToBottom()
  }
}

async function scrollToBottom() {
  await nextTick()
  if (messagesRef.value) {
    messagesRef.value.scrollTop = messagesRef.value.scrollHeight
  }
}

watch(() => messages.value.length, async () => {
  await scrollToBottom()
})

function formatTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

onMounted(async () => {
  if (wrapperRef.value) {
    animate(wrapperRef.value, {
      opacity: [0, 1],
      translateY: [40, 0],
      scale: [0.9, 1],
      duration: 500,
      ease: 'outExpo'
    })
  }

  await loadContacts()
  if (receiverId.value) {
    await loadChatHistory()
  }
})

onUnmounted(() => {
  chatStore.currentReceiverId = null
  websocketStore.setMessages([])
})
</script>

<style scoped>
.chatbox-wrapper {
  position: fixed;
  bottom: var(--space-xl);
  right: var(--space-xl);
  z-index: 1500;
  max-height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
}

.chatbox {
  width: 380px;
  max-height: calc(100vh - 60px);
  border-radius: var(--radius-2xl);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-lg);
  transition: all var(--transition-base);
}

.chatbox.minimized {
  width: 320px;
}


.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-md) var(--space-lg);
  background: linear-gradient(135deg, var(--color-indigo-deep), var(--color-indigo));
  cursor: pointer;
  user-select: none;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.chat-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: var(--gradient-amber);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0b0a12;
  flex-shrink: 0;
}

.chat-title-block {
  display: flex;
  flex-direction: column;
}

.chat-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.chat-status {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.7rem;
  color: var(--color-text-muted);
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-text-muted);
  transition: background var(--transition-fast);
}

.chat-status.online .status-dot {
  background: var(--color-success);
  box-shadow: 0 0 6px var(--color-success);
}

.chat-status.online {
  color: var(--color-success);
}

.chat-header-actions {
  display: flex;
  gap: var(--space-xs);
}

.chat-header-btn {
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.chat-header-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  color: var(--color-text-primary);
}

.close-btn:hover {
  background: rgba(231, 76, 60, 0.3);
  color: var(--color-error);
}


.chat-user-bar {
  padding: var(--space-sm) var(--space-md);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.user-label {
  font-size: 0.72rem;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.user-select {
  padding: 0.4rem 0.6rem;
  font-size: 0.78rem;
  border-radius: var(--radius-sm);
  flex: 1;
}


.chat-messages {
  flex: 1;
  min-height: 150px;
  max-height: 380px;
  overflow-y: auto;
  padding: var(--space-md);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  scroll-behavior: smooth;
}

.chat-welcome {
  text-align: center;
  padding: var(--space-3xl) var(--space-md);
  color: var(--color-text-muted);
}

.welcome-icon {
  font-size: 3rem;
  margin-bottom: var(--space-md);
  animation: float 3s ease-in-out infinite;
}

.chat-welcome p {
  font-size: 0.95rem;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-xs);
}

.chat-welcome span {
  font-size: 0.8rem;
}


.message-row {
  display: flex;
  animation: fadeInUp 0.25s ease-out;
}

.message-row.is-me {
  justify-content: flex-end;
}


.message-bubble {
  max-width: 80%;
  padding: 0.6rem 0.9rem;
  border-radius: var(--radius-lg);
  position: relative;
}

.bubble-me {
  background: linear-gradient(135deg, var(--color-burnt-orange), var(--color-amber));
  color: #0b0a12;
  border-bottom-right-radius: var(--space-xs);
}

.bubble-other {
  background: var(--color-surface);
  color: var(--color-text-primary);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-bottom-left-radius: var(--space-xs);
}

.message-text {
  font-size: 0.85rem;
  line-height: 1.5;
  word-break: break-word;
}

.message-time {
  display: block;
  font-size: 0.65rem;
  margin-top: 0.2rem;
  opacity: 0.6;
  text-align: right;
}


.chat-input-area {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.chat-input {
  flex: 1;
  border-radius: var(--radius-full);
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
}

.btn-send {
  width: 2.6rem;
  height: 2.6rem;
  padding: 0;
  border-radius: var(--radius-full);
  flex-shrink: 0;
}

.btn-send:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}


@media (max-width: 480px) {
  .chatbox-wrapper {
    bottom: 0;
    right: 0;
    left: 0;
  }

  .chatbox {
    width: 100%;
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  }

  .chatbox.minimized {
    width: 100%;
  }
}
</style>
