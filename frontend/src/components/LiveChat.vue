<template>
  <div v-if="authStore.isAuthenticated && authStore.user?.role !== 'seller'" class="live-chat-wrapper">
    
    
    <button 
      v-if="!chatStore.isOpen" 
      @click="chatStore.toggleChat()" 
      class="chat-fab"
      aria-label="Buka Chat"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
      <span v-if="notificationStore.hasBuyerNotifications" class="absolute -top-1 -right-1 flex h-3 w-3 rounded-full bg-red-500 ring-2 ring-white" style="position: absolute; top: -4px; right: -4px; background-color: #ef4444; width: 14px; height: 14px; border-radius: 50%; border: 2px solid white;"></span>
    </button>

    
    <div v-if="chatStore.isOpen" class="chat-window shadow-xl border border-gray-200">
      
      
      <template v-if="currentView === 'history'">
        <header class="chat-header">
          <div class="chat-header-info">
            <h3 class="chat-title">Riwayat Chat</h3>
            <span class="chat-status">
              <span class="status-dot" :class="{ online: websocketStore.isConnected, reconnecting: websocketStore.isReconnecting }"></span>
              {{ websocketStore.isConnected ? 'Terhubung' : (websocketStore.isReconnecting ? 'Menghubungkan Ulang...' : 'Terputus') }}
            </span>
          </div>
          <button @click="chatStore.closeChat()" class="close-btn" title="Tutup Chat">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </header>

        <div class="contacts-area">
          <div v-if="isLoadingContacts" class="chat-loading">
            <div class="spinner"></div>
            <p>Memuat kontak...</p>
          </div>
          <div v-else-if="contacts.length === 0" class="chat-empty-state">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-gray-400 mb-2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
            <p>Belum ada riwayat chat.</p>
          </div>
          <ul v-else class="contacts-list">
            <li 
              v-for="contact in contacts" 
              :key="contact.id"
              class="contact-item"
              @click="openContactChat(contact)"
            >
              <div class="contact-avatar">
                {{ getInitials(contact.name) }}
              </div>
              <div class="contact-info">
                <span class="contact-name">{{ contact.name }}</span>
                <span v-if="contact.unread_count > 0" class="contact-unread">{{ contact.unread_count }}</span>
              </div>
            </li>
          </ul>
        </div>
      </template>

      
      <template v-else-if="currentView === 'chat'">
        <header class="chat-header">
          <div class="chat-header-info">
            <button @click="goBackToList" class="back-btn" title="Kembali" aria-label="Kembali ke list">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="19" y1="12" x2="5" y2="12"></line>
                <polyline points="12 19 5 12 12 5"></polyline>
              </svg>
            </button>
            <div class="chat-avatar">
              {{ getInitials(chatStore.currentReceiverName) }}
            </div>
            <div>
              <h3 class="chat-title">{{ chatStore.currentReceiverName }}</h3>
              <span class="chat-status">
                <span class="status-dot" :class="{ online: websocketStore.isConnected, reconnecting: websocketStore.isReconnecting }"></span>
                {{ websocketStore.isConnected ? 'Terhubung' : (websocketStore.isReconnecting ? 'Menghubungkan Ulang...' : 'Terputus') }}
              </span>
            </div>
          </div>
          <button @click="chatStore.closeChat()" class="close-btn" title="Tutup Chat">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </header>

        <div v-if="isLoadingHistory" class="chat-loading">
          <div class="spinner"></div>
          <p>Memuat riwayat chat...</p>
        </div>

        <div v-else class="chat-messages" ref="messagesContainer">
          <div 
            v-for="(msg, index) in websocketStore.messages" 
            :key="msg.id || index"
            class="message-wrapper"
            :class="[msg.sender_id === authStore.user?.id ? 'message-out' : 'message-in']"
          >
            <div class="message-bubble">
              {{ msg.content }}
            </div>
            <div class="message-time">
              {{ formatTime(msg.created_at) }}
              <span v-if="msg.sender_id === authStore.user?.id && msg.is_read" class="read-receipt">✓✓</span>
            </div>
          </div>
          
          <div v-if="websocketStore.messages.length === 0" class="no-messages">
            Mulai percakapan dengan {{ chatStore.currentReceiverName }}
          </div>
        </div>

        
        <footer class="chat-footer">
          <form @submit.prevent="sendMessage" class="chat-form">
            <input 
              type="text" 
              v-model="newMessage" 
              class="chat-input" 
              placeholder="Ketik pesan..." 
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

    
    <div class="buyer-toast-container" :class="{ 'toast-visible': showBuyerToast }">
      <div class="buyer-toast-content">
        <div class="buyer-toast-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <div class="buyer-toast-text" @click="chatStore.toggleChat()">
          <strong>Pesan Baru</strong>
          <p>{{ toastSenderName }}</p>
        </div>
        <button class="buyer-toast-close" @click="showBuyerToast = false">&times;</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useChatStore } from '../stores/chat'
import { useNotificationStore } from '../stores/notification'
import { useWebsocketStore } from '../stores/websocket'
import axios from 'axios'

const authStore = useAuthStore()
const chatStore = useChatStore()
const notificationStore = useNotificationStore()
const websocketStore = useWebsocketStore()


const currentView = ref('history') 
const contacts = ref([])
const newMessage = ref('')
const messagesContainer = ref(null)


const showBuyerToast = ref(false)
const toastSenderName = ref('')
const isLoadingContacts = ref(false)
const isLoadingHistory = ref(false)

let toastTimeout = null

const globalUnreadCount = computed(() => {
  return contacts.value.reduce((sum, c) => sum + (c.unread_count || 0), 0)
})

const getInitials = (name) => {
  if (!name) return 'CH'
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
  if (currentView.value === 'chat') {
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
      const totalUnread = contacts.value.reduce((sum, c) => sum + (c.unread_count || 0), 0)
      notificationStore.setBuyerUnreadChats(totalUnread)
    }
  } catch (err) {
    console.error('Failed to fetch contacts:', err)
  } finally {
    isLoadingContacts.value = false
  }
}


const fetchChatHistory = async () => {
  if (!authStore.user?.id || !chatStore.currentReceiverId) return
  
  isLoadingHistory.value = true
  try {
    const res = await axios.get('http://localhost:8081/api/messages', {
      params: {
        sender_id: authStore.user.id,
        receiver_id: chatStore.currentReceiverId,
        limit: 100
      }
    })
    
    if (res.data.success) {
      websocketStore.setMessages(res.data.data || [])
      
      const contact = contacts.value.find(c => c.id === chatStore.currentReceiverId)
      if (contact) {
        contact.unread_count = 0
        const totalUnread = contacts.value.reduce((sum, c) => sum + (c.unread_count || 0), 0)
        notificationStore.setBuyerUnreadChats(totalUnread)
      }
      
      await scrollToBottom()
    }
  } catch (err) {
    console.error('Failed to fetch chat history:', err)
  } finally {
    isLoadingHistory.value = false
  }
}


const sendMessage = () => {
  if (!newMessage.value.trim() || !websocketStore.isConnected || !chatStore.currentReceiverId) return
  
  const payload = {
    sender_id: authStore.user.id,
    receiver_id: chatStore.currentReceiverId,
    content: newMessage.value.trim()
  }
  
  websocketStore.sendMessage(payload)
  newMessage.value = ''
}


const openContactChat = (contact) => {
  chatStore.openChat(contact.id, contact.name)
}

const goBackToList = () => {
  chatStore.clearReceiver()
  currentView.value = 'history'
  fetchContacts() 
}


watch(() => chatStore.currentReceiverId, async (newId) => {
  if (newId) {
    currentView.value = 'chat'
    websocketStore.setMessages([])
    await fetchChatHistory()
  } else {
    currentView.value = 'history'
    websocketStore.setMessages([])
  }
})

watch(() => chatStore.isOpen, (isOpen) => {
  if (isOpen) {
    if (!chatStore.currentReceiverId) {
      fetchContacts() 
    } else {
      scrollToBottom() 
    }
  }
})

watch(() => authStore.isAuthenticated, (isAuth) => {
  if (isAuth && authStore.user?.role !== 'seller') {
    fetchContacts()
  }
})


onMounted(() => {
  
  if (authStore.isAuthenticated && authStore.user?.role !== 'seller') {
    fetchContacts()
    
    if (chatStore.currentReceiverId) {
      fetchChatHistory()
    }
  }
})

onUnmounted(() => {
  if (toastTimeout) clearTimeout(toastTimeout)
})
</script>

<style scoped>
.live-chat-wrapper {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 16px;
}

.chat-fab {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #1a1a1a;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.chat-fab:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background-color: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
}

.chat-window {
  width: 360px;
  height: 520px;
  background-color: white;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 10px 40px rgba(0,0,0,0.15);
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}


.chat-header {
  background: #ffffff;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e5e7eb;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  background: none;
  border: none;
  color: #374151;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.back-btn:hover {
  color: #111827;
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f3f4f6;
  color: #374151;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
}

.chat-title {
  margin: 0 0 2px;
  font-size: 0.95rem;
  font-weight: 600;
  color: #111827;
}

.chat-status {
  font-size: 0.75rem;
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

.close-btn {
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  color: #111827;
  background-color: #f3f4f6;
}


.contacts-area {
  flex: 1;
  overflow-y: auto;
  background-color: #fff;
}

.contacts-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 16px;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}

.contact-item:hover {
  background: #f9fafb;
}

.contact-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #e5e7eb;
  color: #374151;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.95rem;
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
}

.contact-unread {
  background-color: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 12px;
}


.chat-empty-state, .chat-loading {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  font-size: 0.9rem;
  text-align: center;
  padding: 20px;
  height: 100%;
}

.chat-messages {
  flex: 1;
  background-color: #f9fafb;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.no-messages {
  text-align: center;
  margin-top: auto;
  margin-bottom: auto;
  font-size: 0.85rem;
  color: #9ca3af;
  background: white;
  padding: 8px 16px;
  border-radius: 20px;
  align-self: center;
}

.message-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 80%;
}

.message-in {
  align-self: flex-start;
}

.message-out {
  align-self: flex-end;
  align-items: flex-end;
}

.message-bubble {
  padding: 10px 14px;
  border-radius: 16px;
  font-size: 0.9rem;
  line-height: 1.4;
  word-wrap: break-word;
}

.message-in .message-bubble {
  background-color: white;
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
  font-size: 0.7rem;
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
  padding: 16px;
  background: white;
  border-top: 1px solid #e5e7eb;
}

.chat-form {
  display: flex;
  gap: 8px;
}

.chat-input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 24px;
  padding: 10px 16px;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s;
}

.chat-input:focus {
  border-color: #111827;
}

.chat-input:disabled {
  background-color: #f3f4f6;
  cursor: not-allowed;
}

.chat-send-btn {
  width: 40px;
  height: 40px;
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
  width: 24px;
  height: 24px;
  border: 3px solid #e5e7eb;
  border-top-color: #111827;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}


.buyer-toast-container {
  position: fixed;
  bottom: 100px;
  right: -350px;
  transition: right 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  z-index: 9999;
}

.buyer-toast-container.toast-visible {
  right: 24px;
}

.buyer-toast-content {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-left: 4px solid #111827;
  min-width: 280px;
  cursor: pointer;
}

.buyer-toast-icon {
  background: #f3f4f6;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #111827;
}

.buyer-toast-text {
  flex: 1;
}

.buyer-toast-text strong {
  display: block;
  font-size: 0.95rem;
  color: #111827;
  margin-bottom: 2px;
}

.buyer-toast-text p {
  margin: 0;
  font-size: 0.85rem;
  color: #6b7280;
}

.buyer-toast-close {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  margin-left: 8px;
}

.buyer-toast-close:hover {
  color: #111827;
}

@media (max-width: 480px) {
  .chat-window {
    width: calc(100vw - 32px);
    height: 60vh;
    max-height: 500px;
  }
}
</style>
