<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-container glass-strong" ref="modalRef">
        <!-- Close Button -->
        <button class="modal-close" @click="$emit('close')" aria-label="Close modal">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M18 6L6 18M6 6l12 12"/>
          </svg>
        </button>

        <!-- Modal Header -->
        <div class="modal-header">
          <div class="modal-badge badge badge-amber">Pre-Order</div>
          <h2 class="modal-title">{{ product.name }}</h2>
          <p class="modal-vendor">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
            </svg>
            {{ product.vendor?.name || 'Vendor NTT' }}
          </p>
        </div>

        <!-- Product Info -->
        <div class="modal-product-info">
          <div class="info-row">
            <span class="info-label">Harga per item</span>
            <span class="info-value text-gradient">{{ formatPrice(product.price) }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Min. order</span>
            <span class="info-value">{{ product.min_order }} pcs</span>
          </div>
          <div class="info-row">
            <span class="info-label">Durasi PO</span>
            <span class="info-value">{{ product.po_duration }} hari</span>
          </div>
          <div class="info-row">
            <span class="info-label">Daerah</span>
            <span class="info-value">{{ product.region }}</span>
          </div>
        </div>

        <!-- Order Form -->
        <form @submit.prevent="submitOrder" class="modal-form">
          <div class="form-group">
            <label class="form-label" for="order-qty">Jumlah Pesanan</label>
            <div class="qty-input-group">
              <button type="button" class="qty-btn" @click="decrementQty">−</button>
              <input
                id="order-qty"
                v-model.number="quantity"
                type="number"
                class="input qty-input"
                :min="product.min_order"
                :max="product.stock"
                required
              />
              <button type="button" class="qty-btn" @click="incrementQty">+</button>
            </div>
            <span v-if="quantity < product.min_order" class="form-error">
              Minimum order: {{ product.min_order }} pcs
            </span>
            <span v-else-if="quantity > product.stock" class="form-error">
              Stok tidak mencukupi (Tersedia: {{ product.stock }} pcs)
            </span>
            <span v-else class="stock-warning">
              Stok tersedia: <strong>{{ product.stock }} pcs</strong>
            </span>
          </div>

          <div class="form-group">
            <label class="form-label" for="order-note">Catatan (opsional)</label>
            <textarea
              id="order-note"
              v-model="note"
              class="input note-input"
              placeholder="Contoh: warna merah, ukuran XL..."
              rows="3"
            ></textarea>
          </div>

          <!-- Total -->
          <div class="order-total">
            <span class="total-label">Total Pembayaran</span>
            <span class="total-value text-gradient">{{ formatPrice(totalPrice) }}</span>
          </div>

          <!-- Actions -->
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="$emit('close')">
              Batal
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              :disabled="submitting || quantity < product.min_order"
            >
              <svg v-if="!submitting" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
              </svg>
              <span v-if="submitting" class="spinner"></span>
              {{ submitting ? 'Memproses...' : 'Konfirmasi Pre-Order' }}
            </button>
          </div>

          <!-- Status Messages -->
          <div v-if="statusMessage" class="status-message" :class="statusType">
            {{ statusMessage }}
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { animate } from 'animejs'

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

const API_BASE = 'http://localhost:8080/api'

const quantity = ref(props.product.min_order || 1)
const note = ref('')
const submitting = ref(false)
const statusMessage = ref('')
const statusType = ref('')
const modalRef = ref(null)

const totalPrice = computed(() => {
  return (props.product.price || 0) * (quantity.value || 0)
})

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(price)
}

function incrementQty() {
  if (quantity.value < props.product.stock) {
    quantity.value++
  }
}

function decrementQty() {
  if (quantity.value > props.product.min_order) {
    quantity.value--
  }
}

async function submitOrder() {
  if (quantity.value < props.product.min_order) {
    statusType.value = 'error'
    statusMessage.value = `Minimal order adalah ${props.product.min_order} pcs`
    return
  }

  if (quantity.value > props.product.stock) {
    statusType.value = 'error'
    statusMessage.value = `Kuantitas pesanan melebihi stok tersedia (${props.product.stock} pcs)`
    return
  }

  submitting.value = true
  statusMessage.value = ''

  try {
    const response = await fetch(`${API_BASE}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        customer_id: 4, // Demo customer ID
        product_id: props.product.id,
        quantity: quantity.value,
        note: note.value
      })
    })

    const data = await response.json()

    if (data.success) {
      statusType.value = 'success'
      statusMessage.value = '✅ Pre-order berhasil! Pesanan Anda telah dicatat.'
      setTimeout(() => emit('close'), 2500)
    } else {
      statusType.value = 'error'
      statusMessage.value = data.message || 'Gagal membuat pesanan'
    }
  } catch (error) {
    statusType.value = 'error'
    statusMessage.value = 'Tidak dapat terhubung ke server. Pastikan backend berjalan.'
  } finally {
    submitting.value = false
  }
}

// Entrance animation
onMounted(() => {
  if (modalRef.value) {
    animate(modalRef.value, {
      opacity: [0, 1],
      scale: [0.9, 1],
      translateY: [30, 0],
      duration: 400,
      ease: 'outExpo'
    })
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-xl);
}

.modal-container {
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  border-radius: var(--radius-2xl);
  padding: var(--space-2xl);
  position: relative;
}

.modal-close {
  position: absolute;
  top: var(--space-lg);
  right: var(--space-lg);
  background: rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-full);
  width: 2.2rem;
  height: 2.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  transition: all var(--transition-fast);
}

.modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text-primary);
}

/* Header */
.modal-header {
  margin-bottom: var(--space-lg);
}

.modal-badge {
  margin-bottom: var(--space-sm);
}

.modal-title {
  font-size: 1.3rem;
  font-weight: 700;
  line-height: 1.3;
  margin-bottom: var(--space-xs);
}

.modal-vendor {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.modal-vendor svg {
  color: var(--color-amber);
}

/* Product Info */
.modal-product-info {
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-lg);
  padding: var(--space-md) var(--space-lg);
  margin-bottom: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.info-label {
  font-size: 0.82rem;
  color: var(--color-text-muted);
}

.info-value {
  font-size: 0.9rem;
  font-weight: 600;
}

/* Form */
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.form-label {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.qty-input-group {
  display: flex;
  gap: var(--space-xs);
}

.qty-btn {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-md);
  background: var(--color-surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--color-text-primary);
  font-size: 1.2rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.qty-btn:hover {
  background: var(--color-surface-hover);
  border-color: var(--color-indigo-glow);
}

.qty-input {
  width: 80px;
  text-align: center;
  font-weight: 700;
  font-size: 1.1rem;
}

.note-input {
  resize: vertical;
  min-height: 80px;
}

.form-error {
  font-size: 0.75rem;
  color: var(--color-error);
}

.stock-warning {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.2rem;
}

/* Total */
.order-total {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-md) var(--space-lg);
  background: rgba(245, 166, 35, 0.06);
  border: 1px solid rgba(245, 166, 35, 0.15);
  border-radius: var(--radius-lg);
}

.total-label {
  font-size: 0.9rem;
  font-weight: 600;
}

.total-value {
  font-size: 1.4rem;
  font-weight: 800;
}

/* Actions */
.modal-actions {
  display: flex;
  gap: var(--space-sm);
  justify-content: flex-end;
}

.modal-actions .btn {
  flex: 1;
}

/* Status */
.status-message {
  padding: var(--space-md);
  border-radius: var(--radius-md);
  font-size: 0.85rem;
  text-align: center;
  font-weight: 500;
}

.status-message.success {
  background: rgba(46, 204, 113, 0.1);
  color: var(--color-success);
  border: 1px solid rgba(46, 204, 113, 0.2);
}

.status-message.error {
  background: rgba(231, 76, 60, 0.1);
  color: var(--color-error);
  border: 1px solid rgba(231, 76, 60, 0.2);
}

/* Spinner */
.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(11, 10, 18, 0.3);
  border-top-color: #0b0a12;
  border-radius: 50%;
  animation: spin-slow 0.6s linear infinite;
}
</style>
