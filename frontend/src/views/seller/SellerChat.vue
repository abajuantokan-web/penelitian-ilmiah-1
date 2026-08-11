<template>
  <div class="seller-chat-container">
    <div class="chat-sidebar">
      <div class="chat-sidebar-header">
        <h2>Pesan Masuk</h2>
      </div>
      <div v-if="isLoadingContacts" class="contacts-loading">
        <div class="spinner"></div>
      </div>
      <div v-else-if="contacts.length === 0" class="contacts-empty">
        Belum ada pesan dari pembeli.
      </div>
      <ul v-else class="contacts-list">
        <li 
          v-for="contact in contacts" 
          :key="contact.id"
          :class="['contact-item', { active: activeContact?.id === contact.id }]"
          @click="selectContact(contact)"
        >
          <div class="contact-avatar">
            {{ getInitials(contact.name) }}
          </div>
          <div class="contact-info">
            <span class="contact-name">{{ contact.name }}</span>
            <span v-if="contact.unread_count > 0" class="unread-badge">{{ contact.unread_count }}</span>
          </div>
        </li>
      </ul>
    </div>

    <div class="chat-main">
      <div v-if="!activeContact" class="chat-empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-gray-400 mb-2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
        <p>Pilih pembeli untuk melihat pesan</p>
      </div>
      
      <template v-else>
        <header class="chat-header">
          <div class="chat-header-info">
            <div class="chat-avatar">{{ getInitials(activeContact.name) }}</div>
            <div>
              <h3 class="chat-title">{{ activeContact.name }}</h3>
              <span class="chat-status">
                <span class="status-dot" :class="{ online: websocketStore.isConnected, reconnecting: websocketStore.isReconnecting }"></span>
                {{ websocketStore.isConnected ? 'Terhubung' : (websocketStore.isReconnecting ? 'Menghubungkan Ulang...' : 'Terputus') }}
              </span>
            </div>
          </div>
        </header>

        <div class="chat-messages" ref="messagesContainer">
          <div v-if="isLoadingHistory" class="chat-loading">
            <div class="spinner"></div>
          </div>
          <div v-else-if="websocketStore.messages.length === 0" class="no-messages">
            Mulai percakapan dengan {{ activeContact.name }}
          </div>
          <template v-else>
            <div 
              v-for="(msg, index) in websocketStore.messages" 
              :key="msg.id || index"
              class="message-wrapper"
              :class="[msg.sender_id === authStore.user?.id ? 'message-out' : 'message-in']"
            >
              <div class="message-bubble">{{ msg.content }}</div>
              <div class="message-time">
                {{ formatTime(msg.created_at) }}
                <span v-if="msg.sender_id === authStore.user?.id && msg.is_read" class="read-receipt">✓✓</span>
              </div>
            </div>
          </template>
        </div>

        <footer class="chat-footer">
          <form @submit.prevent="sendMessage" class="chat-form">
            <input 
              type="text" 
              v-model="newMessage" 
              class="chat-input" 
              placeholder="Ketik balasan Anda..." 
              :disabled="!websocketStore.isConnected"
            />
            <button 
              type="submit" 
              class="chat-send-btn"
              :disabled="!newMessage.trim() || !websocketStore.isConnected"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="22" y1="2" x2="11" y2="13"></line>
                <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
              </svg>
            </button>
          </form>
        </footer>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import { useWebsocketStore } from '../../stores/websocket'
import { useChatStore } from '../../stores/chat'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const websocketStore = useWebsocketStore()
const chatStore = useChatStore()


const contacts = ref([])
const activeContact = ref(null)
const newMessage = ref('')
const messagesContainer = ref(null)

const isLoadingContacts = ref(true)
const isLoadingHistory = ref(false)


const emit = defineEmits(['new-message-received', 'update-unread'])

watch(contacts, (newContacts) => {
  const total = newContacts.reduce((sum, c) => sum + (c.unread_count || 0), 0)
  emit('update-unread', total)
}, { deep: true })

const getInitials = (name) => {
  if (!name) return 'U'
  return name.substring(0, 2).toUpperCase()
}

const formatTime = (isoString) => {
  if (!isoString) return new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  const date = new Date(isoString)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

watch(() => websocketStore.messages.length, async () => {
  if (activeContact.value) {
    await scrollToBottom()
  }
})

const fetchContacts = async () => {
  if (!authStore.user?.id) return
  isLoadingContacts.value = true
  try {
    const res = await axios.get('http://localhost:8081/api/chat/contacts', {
      params: { user_id: authStore.user.id, role: authStore.user.role }
    })
    if (res.data.success) {
      contacts.value = res.data.data || []
    }
  } catch (err) {
    console.error('Gagal mengambil kontak:', err)
  } finally {
    isLoadingContacts.value = false
  }
}

const fetchChatHistory = async () => {
  if (!activeContact.value) return
  isLoadingHistory.value = true
  try {
    const res = await axios.get('http://localhost:8081/api/messages', {
      params: {
        sender_id: authStore.user.id,
        receiver_id: activeContact.value.id,
        limit: 100
      }
    })
    if (res.data.success) {
      websocketStore.setMessages(res.data.data || [])
      
      const contact = contacts.value.find(c => c.id === activeContact.value.id)
      if (contact) contact.unread_count = 0
      await scrollToBottom()
    }
  } catch (err) {
    console.error('Failed to fetch chat history:', err)
  } finally {
    isLoadingHistory.value = false
  }
}


const sendMessage = () => {
  if (!newMessage.value.trim() || !websocketStore.isConnected || !activeContact.value) return
  
  const payload = {
    sender_id: authStore.user.id,
    receiver_id: activeContact.value.id,
    content: newMessage.value.trim()
  }
  
  websocketStore.sendMessage(payload)
  newMessage.value = ''
}

const selectContact = (contact) => {
  activeContact.value = contact
  chatStore.currentReceiverId = contact.id
  websocketStore.setMessages([])
  fetchChatHistory()
}

onMounted(() => {
  chatStore.currentReceiverId = null
  notificationStore.resetSellerUnreadChats()
  fetchContacts()
})

onUnmounted(() => {
})

defineExpose({
  fetchContacts
})
</script>

<style scoped>
.seller-chat-container {
  display: flex;
  height: calc(100vh - 200px);
  min-height: 500px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.chat-sidebar {
  width: 280px;
  border-right: 1px solid #e5e7eb;
  background: #f9fafb;
  display: flex;
  flex-direction: column;
}

.chat-sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #fff;
}

.chat-sidebar-header h2 {
  margin: 0;
  font-size: 1.1rem;
  font-family: 'Playfair Display', serif;
}

.contacts-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}

.contact-item:hover {
  background: #f3f4f6;
}

.contact-item.active {
  background: #fff;
  border-left: 4px solid #111827;
}

.contact-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e5e7eb;
  color: #374151;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
}

.contact-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.contact-name {
  font-weight: 500;
  color: #111827;
  font-size: 0.95rem;
}

.unread-badge {
  background-color: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 12px;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.chat-header {
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chat-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #374151;
}

.chat-title {
  margin: 0 0 2px;
  font-size: 1rem;
  font-weight: 600;
}

.chat-status {
  font-size: 0.8rem;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #9ca3af;
}

.status-dot.online {
  background-color: #10b981;
}

.status-dot.reconnecting {
  background-color: #eab308;
}

.chat-messages {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #fafafa;
}

.message-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 70%;
}

.message-in {
  align-self: flex-start;
}

.message-out {
  align-self: flex-end;
  align-items: flex-end;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 0.95rem;
  line-height: 1.4;
  word-wrap: break-word;
}

.message-in .message-bubble {
  background-color: #fff;
  color: #111827;
  border-bottom-left-radius: 4px;
  border: 1px solid #e5e7eb;
}

.message-out .message-bubble {
  background-color: #111827;
  color: white;
  border-bottom-right-radius: 4px;
}

.message-time {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.read-receipt {
  color: #3b82f6;
}

.chat-footer {
  padding: 20px 24px;
  background: #fff;
  border-top: 1px solid #e5e7eb;
}

.chat-form {
  display: flex;
  gap: 12px;
}

.chat-input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 24px;
  padding: 12px 20px;
  font-size: 0.95rem;
  outline: none;
  transition: border-color 0.2s;
}

.chat-input:focus {
  border-color: #111827;
}

.chat-send-btn {
  width: 46px;
  height: 46px;
  border-radius: 50%;
  background-color: #111827;
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.chat-send-btn:hover:not(:disabled) {
  background-color: #1f2937;
}

.chat-send-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e5e7eb;
  border-top-color: #111827;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: auto;
}

.contacts-loading, .contacts-empty {
  padding: 40px 20px;
  text-align: center;
  color: #6b7280;
  font-size: 0.9rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
