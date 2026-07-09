import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  // Seller State
  const sellerUnreadChatsCount = ref(0)
  const sellerNewOrdersCount = ref(0)

  // Buyer State
  const buyerUnreadChatsCount = ref(0)

  // Getters
  const hasSellerNotifications = computed(() => {
    return sellerUnreadChatsCount.value > 0 || sellerNewOrdersCount.value > 0
  })

  const hasBuyerNotifications = computed(() => {
    return buyerUnreadChatsCount.value > 0
  })

  // Seller Actions
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

  // Buyer Actions
  const incrementBuyerUnreadChats = () => {
    buyerUnreadChatsCount.value++
  }

  const setBuyerUnreadChats = (count) => {
    buyerUnreadChatsCount.value = count
  }

  const resetBuyerUnreadChats = () => {
    buyerUnreadChatsCount.value = 0
  }

  // Clear all on logout
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
