import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
  const isOpen = ref(false)
  const isSellerChatOpen = ref(false)
  const currentReceiverId = ref(null)
  const currentReceiverName = ref('Chat')

  const openChat = (receiverId, receiverName) => {
    currentReceiverId.value = receiverId
    currentReceiverName.value = receiverName || 'Chat'
    isOpen.value = true
  }

  const closeChat = () => {
    isOpen.value = false
  }

  const toggleChat = () => {
    isOpen.value = !isOpen.value
  }

  const clearReceiver = () => {
    currentReceiverId.value = null
    currentReceiverName.value = 'Chat'
  }

  return {
    isOpen,
    isSellerChatOpen,
    currentReceiverId,
    currentReceiverName,
    openChat,
    closeChat,
    toggleChat,
    clearReceiver
  }
})
