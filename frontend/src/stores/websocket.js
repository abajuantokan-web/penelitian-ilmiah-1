import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useAuthStore } from './auth'
import { useChatStore } from './chat'
import { useNotificationStore } from './notification'
import { useDashboardStore } from './dashboard'
import { BASE_URL } from '../axios'

export const useWebsocketStore = defineStore('websocket', () => {
  const ws = ref(null)
  const isConnected = ref(false)
  const isReconnecting = ref(false)
  const messages = ref([])
  
  let reconnectInterval = null
  let reconnectAttempts = 0

  const authStore = useAuthStore()
  const chatStore = useChatStore()
  const notificationStore = useNotificationStore()
  const dashboardStore = useDashboardStore()

  const playNotificationSound = () => {
    try {
      const AudioContext = window.AudioContext || window.webkitAudioContext;
      const ctx = new AudioContext();
      const osc = ctx.createOscillator();
      const gainNode = ctx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(900, ctx.currentTime);
      osc.frequency.exponentialRampToValueAtTime(1400, ctx.currentTime + 0.15);
      gainNode.gain.setValueAtTime(0.3, ctx.currentTime);
      gainNode.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.15);
      osc.connect(gainNode);
      gainNode.connect(ctx.destination);
      osc.start();
      osc.stop(ctx.currentTime + 0.15);
    } catch (e) {
      console.log("Audio not supported", e);
    }
  }

  const connectWebSocket = () => {
    if (!authStore.user?.id) return
    if (ws.value && (ws.value.readyState === WebSocket.OPEN || ws.value.readyState === WebSocket.CONNECTING)) {
      return
    }

    isReconnecting.value = true
    const token = authStore.token || localStorage.getItem('openpeo_token')
    
    // Konversi URL HTTP/HTTPS ke WS/WSS dengan benar
    // https://... → wss://... | http://... → ws://...
    const wsUrl = BASE_URL.replace(/^https:\/\//, 'wss://').replace(/^http:\/\//, 'ws://')
    // Hapus trailing slash jika ada agar path tidak double-slash
    const wsBase = wsUrl.replace(/\/$/, '')
    ws.value = new WebSocket(`${wsBase}/ws/chat?token=${token}`)

    ws.value.onopen = () => {
      console.log('🔗 WebSocket Connected (Singleton)')
      isConnected.value = true
      isReconnecting.value = false
      reconnectAttempts = 0
    }

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        // Handle system events (orders, etc)
        if (data.type) {
          if (data.type === 'NEW_ORDER_CREATED') {
            dashboardStore.handleNewOrder(data.data)
            playNotificationSound()
          } else if (data.type === 'ORDER_STATUS_UPDATED') {
            dashboardStore.handleOrderStatusUpdated(data.data)
          }
          return
        }

        // Determine if this message belongs to the active chat conversation
        const myId = authStore.user?.id
        const partnerId = chatStore.currentReceiverId

        const isActiveConversation = partnerId && (
          (data.sender_id === myId && data.receiver_id === partnerId) ||
          (data.sender_id === partnerId && data.receiver_id === myId)
        )

        if (isActiveConversation) {
          // For my own sent message: replace the optimistic temp message with the confirmed one
          if (data.sender_id === myId) {
            const tempIndex = messages.value.findIndex(m =>
              !m.id && m.sender_id === myId && m.content === data.content
            )
            if (tempIndex !== -1) {
              messages.value[tempIndex] = data
            } else {
              const exists = messages.value.some(m => m.id === data.id)
              if (!exists) messages.value.push(data)
            }
          } else {
            // Incoming message from partner
            const exists = messages.value.some(m => m.id === data.id)
            if (!exists) messages.value.push(data)
          }
        }

        // Notification logic (only for messages from others)
        if (data.sender_id === myId) return

        const isFromOther = !partnerId || data.sender_id !== partnerId
        const isBuyer = authStore.user?.role !== 'seller'
        const isSeller = authStore.user?.role === 'seller'

        if (isBuyer) {
          if (!chatStore.isOpen || isFromOther) {
            notificationStore.incrementBuyerUnreadChats()
            playNotificationSound()
          }
        } else if (isSeller) {
          if (!chatStore.isSellerChatOpen || isFromOther) {
            notificationStore.incrementSellerUnreadChats()
            playNotificationSound()
          }
        }

      } catch (err) {
        console.error('Failed to parse incoming message:', err)
      }
    }

    const handleDisconnect = () => {
      isConnected.value = false
      isReconnecting.value = true
      
      const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 10000)
      console.log(`WebSocket Terputus. Menghubungkan ulang dalam ${delay}ms...`)
      
      clearTimeout(reconnectInterval)
      reconnectInterval = setTimeout(() => {
        reconnectAttempts++
        connectWebSocket()
      }, delay)
    }

    ws.value.onclose = handleDisconnect
    ws.value.onerror = handleDisconnect
  }

  const disconnectWebSocket = () => {
    clearTimeout(reconnectInterval)
    if (ws.value) {
      ws.value.onclose = null
      ws.value.onerror = null
      ws.value.onmessage = null
      ws.value.close()
      ws.value = null
    }
    isConnected.value = false
    isReconnecting.value = false
  }

  const sendMessage = (payload) => {
    if (ws.value && isConnected.value) {
      ws.value.send(JSON.stringify(payload))

      // Optimistic temp message — id is null so server echo can replace it
      const tempMessage = {
        id: null,
        sender_id: payload.sender_id,
        receiver_id: payload.receiver_id,
        content: payload.content,
        created_at: new Date().toISOString()
      }

      messages.value.push(tempMessage)
    }
  }

  const setMessages = (newMessages) => {
    messages.value = newMessages
  }

  
  watch(() => authStore.isAuthenticated, (isAuth) => {
    if (isAuth) {
      connectWebSocket()
    } else {
      disconnectWebSocket()
    }
  }, { immediate: true })

  return {
    ws,
    isConnected,
    isReconnecting,
    messages,
    connectWebSocket,
    disconnectWebSocket,
    sendMessage,
    setMessages
  }
})
