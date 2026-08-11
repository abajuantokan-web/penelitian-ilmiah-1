import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  
  const sellerUnreadChatsCount = ref(0)
  const sellerNewOrdersCount = ref(0)

  
  const buyerUnreadChatsCount = ref(0)

  
  const hasSellerNotifications = computed(() => {
    return sellerUnreadChatsCount.value > 0 || sellerNewOrdersCount.value > 0
  })

  const hasBuyerNotifications = computed(() => {
    return buyerUnreadChatsCount.value > 0
  })

  
  const incrementSellerUnreadChats = () => {
    sellerUnreadChatsCount.value++
  }

  const setSellerUnreadChats = (count) => {
    sellerUnreadChatsCount.value = count
  }

  const resetSellerUnreadChats = () => {
    sellerUnreadChatsCount.value = 0
  }

  const incrementSellerNewOrders = () => {
    sellerNewOrdersCount.value++
  }

  const resetSellerNewOrders = () => {
    sellerNewOrdersCount.value = 0
  }

  
  const incrementBuyerUnreadChats = () => {
    buyerUnreadChatsCount.value++
  }

  const setBuyerUnreadChats = (count) => {
    buyerUnreadChatsCount.value = count
  }

  const resetBuyerUnreadChats = () => {
    buyerUnreadChatsCount.value = 0
  }

  
  const resetAll = () => {
    sellerUnreadChatsCount.value = 0
    sellerNewOrdersCount.value = 0
    buyerUnreadChatsCount.value = 0
  }

  return {
    sellerUnreadChatsCount,
    sellerNewOrdersCount,
    buyerUnreadChatsCount,
    hasSellerNotifications,
    hasBuyerNotifications,
    incrementSellerUnreadChats,
    setSellerUnreadChats,
    resetSellerUnreadChats,
    incrementSellerNewOrders,
    resetSellerNewOrders,
    incrementBuyerUnreadChats,
    setBuyerUnreadChats,
    resetBuyerUnreadChats,
    resetAll
  }
})
