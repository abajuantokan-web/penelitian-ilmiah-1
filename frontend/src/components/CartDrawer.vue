<template>
  <Teleport to="body">
    
    <Transition name="backdrop">
      <div v-if="cartStore.isOpen" class="cart-backdrop" @click="cartStore.closeDrawer"></div>
    </Transition>

    
    <Transition name="drawer">
      <aside v-if="cartStore.isOpen" class="cart-drawer" role="dialog" aria-label="Keranjang Belanja">
        
        <div class="cart-drawer__header">
          <h2 class="cart-drawer__title">Keranjang Anda</h2>
          <span class="cart-drawer__count">{{ cartStore.totalItems }} item</span>
          <button class="cart-drawer__close" @click="cartStore.closeDrawer" aria-label="Tutup keranjang">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        
        <div class="cart-drawer__body">
          
          <div v-if="cartStore.isLoading" class="cart-drawer__loading">
            <div class="cart-spinner"></div>
            <p>Memuat keranjang...</p>
          </div>

          
          <div v-else-if="!cartStore.items || cartStore.items.length === 0" class="cart-drawer__empty">
            <div class="cart-empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
                <line x1="3" y1="6" x2="21" y2="6"/>
                <path d="M16 10a4 4 0 0 1-8 0"/>
              </svg>
            </div>
            <p class="cart-empty-title">Keranjang Anda kosong</p>
            <p class="cart-empty-desc">Jelajahi koleksi premium kami dan temukan produk autentik NTT.</p>
            <button class="cart-empty-btn" @click="cartStore.closeDrawer">Belanja Sekarang</button>
          </div>

          
          <div v-else class="cart-drawer__items">
            <TransitionGroup name="cart-item">
              <div v-for="item in cartStore.items" :key="item.id" class="cart-item">
                <div class="cart-item__image">
                  <img :src="$getImageUrl(item.product?.image_url)" :alt="item.product?.name" loading="lazy">
                </div>
                <div class="cart-item__info">
                  <p class="cart-item__name">{{ item.product?.name }}</p>
                  <p class="cart-item__meta">{{ item.product?.region }} · Qty: {{ item.quantity }}</p>
                  <p class="cart-item__price">{{ formatPrice(item.product?.price * item.quantity) }}</p>
                </div>
                <button 
                  class="cart-item__remove" 
                  @click="handleRemove(item.id)"
                  :disabled="removingId === item.id"
                  aria-label="Hapus item"
                >
                  <svg v-if="removingId !== item.id" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                  <span v-else class="cart-item__removing-spinner"></span>
                </button>
              </div>
            </TransitionGroup>
          </div>
        </div>

        
        <div v-if="cartStore.items && cartStore.items.length > 0" class="cart-drawer__footer">
          <div class="cart-drawer__total">
            <span>Total</span>
            <span class="cart-drawer__total-price">{{ formatPrice(cartStore.totalPrice) }}</span>
          </div>
          <div class="cart-drawer__preorder-note">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            <span>Produk ini bersifat pre-order. Waktu pembuatan bervariasi.</span>
          </div>
          
          <div v-if="checkoutMsg" class="checkout-msg" :class="{'success': checkoutMsg === 'Pesanan berhasil dibuat!'}">
            {{ checkoutMsg }}
          </div>
          
          <button 
            class="cart-drawer__checkout-btn" 
            @click="handleCheckout"
            :disabled="isCheckingOut"
          >
            {{ isCheckingOut ? 'Memproses...' : 'Lanjut ke Pembayaran' }}
          </button>
        </div>
      </aside>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'


const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const removingId = ref(null)
const isCheckingOut = ref(false)
const checkoutMsg = ref('')


const formatPrice = (price) => {
  if (!price) return 'Rp0'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

const handleRemove = async (cartItemId) => {
  removingId.value = cartItemId
  await cartStore.removeFromCart(cartItemId)
  removingId.value = null
}

const handleCheckout = () => {
  cartStore.closeDrawer()
  router.push({ path: '/checkout', query: { source: 'cart' } })
}
</script>

<style scoped>

.cart-backdrop {
  position: fixed;
  inset: 0;
  z-index: 10001;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(3px);
}

.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.3s ease;
}
.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}


.cart-drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  z-index: 10002;
  width: 420px;
  max-width: 92vw;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  box-shadow: -8px 0 32px rgba(0, 0, 0, 0.12);
}

.drawer-enter-active,
.drawer-leave-active {
  transition: transform 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}
.drawer-enter-from,
.drawer-leave-to {
  transform: translateX(100%);
}


.cart-drawer__header {
  display: flex;
  align-items: center;
  padding: 24px 28px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
}

.cart-drawer__title {
  font-family: 'Playfair Display', serif;
  font-size: 1.3rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0;
}

.cart-drawer__count {
  margin-left: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: #999;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.cart-drawer__close {
  margin-left: auto;
  background: none;
  border: none;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}

.cart-drawer__close:hover {
  background-color: #f5f5f5;
  color: #1a1a1a;
}


.cart-drawer__body {
  flex: 1;
  overflow-y: auto;
  padding: 0;
}


.cart-drawer__loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  gap: 16px;
  color: #999;
  font-size: 0.85rem;
}

.cart-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #eee;
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}


.cart-drawer__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  text-align: center;
}

.cart-empty-icon {
  color: #ccc;
  margin-bottom: 24px;
}

.cart-empty-title {
  font-family: 'Playfair Display', serif;
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px;
}

.cart-empty-desc {
  font-size: 0.85rem;
  color: #999;
  margin: 0 0 28px;
  line-height: 1.6;
  max-width: 260px;
}

.cart-empty-btn {
  padding: 12px 32px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.8rem;
  font-weight: 500;
  letter-spacing: 1px;
  text-transform: uppercase;
  background-color: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.cart-empty-btn:hover {
  background-color: #333;
}


.cart-drawer__items {
  padding: 8px 0;
}

.cart-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 28px;
  border-bottom: 1px solid #f5f5f5;
  transition: background-color 0.2s;
}

.cart-item:hover {
  background-color: #fafafa;
}

.cart-item__image {
  width: 72px;
  height: 72px;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
  background-color: #f5f5f5;
}

.cart-item__image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cart-item__info {
  flex: 1;
  min-width: 0;
}

.cart-item__name {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 4px;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cart-item__meta {
  font-size: 0.72rem;
  color: #999;
  margin: 0 0 6px;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.cart-item__price {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.cart-item__remove {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: #bbb;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}

.cart-item__remove:hover {
  background-color: #fef2f2;
  color: #dc2626;
}

.cart-item__remove:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.cart-item__removing-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid #eee;
  border-top-color: #999;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  display: block;
}


.cart-item-enter-active {
  transition: all 0.3s ease;
}
.cart-item-leave-active {
  transition: all 0.25s ease;
}
.cart-item-enter-from {
  opacity: 0;
  transform: translateX(20px);
}
.cart-item-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}


.cart-drawer__footer {
  border-top: 1px solid #eee;
  padding: 24px 28px;
  flex-shrink: 0;
  background-color: #fafafa;
}

.cart-drawer__total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.cart-drawer__total span:first-child {
  font-size: 0.8rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #666;
}

.cart-drawer__total-price {
  font-family: 'Playfair Display', serif;
  font-size: 1.3rem;
  font-weight: 700;
  color: #1a1a1a;
}

.cart-drawer__preorder-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 0.72rem;
  color: #999;
  margin-bottom: 20px;
  line-height: 1.5;
}

.cart-drawer__preorder-note svg {
  flex-shrink: 0;
  margin-top: 1px;
}

.cart-drawer__checkout-btn {
  width: 100%;
  padding: 16px 32px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.82rem;
  font-weight: 500;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  background-color: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.1s;
}

.cart-drawer__checkout-btn:hover:not(:disabled) {
  background-color: #333;
}

.cart-drawer__checkout-btn:active:not(:disabled) {
  transform: scale(0.99);
}

.cart-drawer__checkout-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.checkout-msg {
  font-size: 0.85rem;
  margin-bottom: 12px;
  padding: 8px 12px;
  border-radius: 4px;
  background-color: #fef2f2;
  color: #dc2626;
  text-align: center;
}

.checkout-msg.success {
  background-color: #f0fdf4;
  color: #16a34a;
}


@media (max-width: 480px) {
  .cart-drawer {
    width: 100vw;
    max-width: 100vw;
  }

  .cart-item {
    padding: 14px 20px;
  }

  .cart-drawer__header,
  .cart-drawer__footer {
    padding-left: 20px;
    padding-right: 20px;
  }
}
</style>
