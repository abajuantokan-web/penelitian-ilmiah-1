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
    const token = authStore.token || localStorage.getItem('token')
    
    const wsUrl = BASE_URL.replace(/^http/, 'ws')
    ws.value = new WebSocket(`${wsUrl}/ws/chat?token=${token}`)

    ws.value.onopen = () => {
      console.log('🔗 WebSocket Connected (Singleton)')
      isConnected.value = true
      isReconnecting.value = false
      reconnectAttempts = 0
    }

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        
        
        if (data.type) {
          if (data.type === 'NEW_ORDER_CREATED') {
            dashboardStore.handleNewOrder(data.data)
            playNotificationSound()
          } else if (data.type === 'ORDER_STATUS_UPDATED') {
            dashboardStore.handleOrderStatusUpdated(data.data)
          }
          return 
        }
        
        
        if (data.sender_id === authStore.user?.id) {
          
          
          return
        }
        
        
        if (chatStore.currentReceiverId && (data.sender_id === chatStore.currentReceiverId || data.receiver_id === chatStore.currentReceiverId)) {
          const exists = messages.value.some(m => m.id === data.id)
          if (!exists) {
            messages.value.push(data)
          }
        }
        
        
        
        const isFromOther = !chatStore.currentReceiverId || (data.sender_id !== chatStore.currentReceiverId && data.receiver_id !== chatStore.currentReceiverId)
        
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
      
      
      const tempMessage = {
        id: Date.now(), 
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
