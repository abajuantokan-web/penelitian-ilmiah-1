import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from './auth'

export const useDashboardStore = defineStore('dashboard', () => {
  const pendingCount = ref(0)
  const processingCount = ref(0)
  const completedCount = ref(0)
  const totalRevenue = ref(0)
  const recentOrders = ref([])
  
  const chartData = ref([])
  const chartLabels = ref([])
  
  const isLoading = ref(false)
  const error = ref(null)

  const authStore = useAuthStore()

  const fetchInitialStats = async () => {
    isLoading.value = true
    error.value = null
    try {
      const response = await axios.get('http://localhost:8081/api/seller/dashboard/stats', {
        headers: {
          Authorization: `Bearer ${authStore.token || localStorage.getItem('token')}`
        }
      })
      if (response.data.success) {
        const data = response.data.data
        pendingCount.value = data.total_orders_pending || 0
        processingCount.value = data.total_orders_processing || 0
        completedCount.value = data.total_completed_orders || 0
        totalRevenue.value = data.total_revenue || 0
        recentOrders.value = data.recent_orders || []
      }
    } catch (err) {
      console.error('Failed to fetch dashboard stats:', err)
      error.value = err.response?.data?.message || 'Gagal memuat statistik dasbor'
    } finally {
      isLoading.value = false
    }
  }

  const fetchChartData = async (range = '7_days') => {
    try {
      const response = await axios.get(`http://localhost:8081/api/seller/dashboard/chart?range=${range}`, {
        headers: {
          Authorization: `Bearer ${authStore.token || localStorage.getItem('token')}`
        }
      })
      if (response.data.success) {
        const dbData = response.data.data || []
        
        // Generate continuous dates based on range
        let days = 7;
        if (range === '30_days') days = 30;
        else if (range === '12_months') days = 365;

        const dateMap = new Map();
        const today = new Date();
        
        for (let i = days - 1; i >= 0; i--) {
          const d = new Date(today);
          d.setDate(d.getDate() - i);
          const dateStr = d.toISOString().split('T')[0];
          
          let displayLabel = d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
          if (range === '12_months') {
            displayLabel = d.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' });
          }
          
          dateMap.set(dateStr, { label: displayLabel, total: 0 });
        }
        
        // Merge with DB data
        dbData.forEach(item => {
          // item.date is format 'YYYY-MM-DD' from SQL DATE()
          const dtStr = item.date.split('T')[0];
          if (dateMap.has(dtStr)) {
            dateMap.get(dtStr).total = item.total;
          }
        });
        
        // If 12 months, we should probably group by month in JS, but for now we'll just plot days or simplify it.
        // The prompt says "fill in missing dates with 0 so the chart line is continuous".
        // We'll extract labels and data directly from the continuous map.
        
        // If 12 months, it's 365 days. We might want to group them by month for display, 
        // but the DB is returning dates. We'll aggregate by label if multiple dates share the same label (e.g. 12 months)
        const aggregatedMap = new Map();
        Array.from(dateMap.values()).forEach(val => {
           if (!aggregatedMap.has(val.label)) {
              aggregatedMap.set(val.label, 0);
           }
           aggregatedMap.set(val.label, aggregatedMap.get(val.label) + val.total);
        });

        chartLabels.value = Array.from(aggregatedMap.keys());
        chartData.value = Array.from(aggregatedMap.values());
      }
    } catch (err) {
      console.error('Failed to fetch chart data:', err)
    }
  }

  const handleNewOrder = (order) => {
    pendingCount.value++
    recentOrders.value.unshift(order)
    if (recentOrders.value.length > 5) {
      recentOrders.value.pop()
    }
  }

  const handleOrderStatusUpdated = (order) => {
    if (order.status === 'Pesanan Sedang Diproses') {
      if (pendingCount.value > 0) pendingCount.value--
      processingCount.value++
    }
    const index = recentOrders.value.findIndex(o => o.id === order.id)
    if (index !== -1) {
      recentOrders.value[index] = order
    }
  }

  return {
    pendingCount,
    processingCount,
    completedCount,
    totalRevenue,
    recentOrders,
    chartData,
    chartLabels,
    isLoading,
    error,
    fetchInitialStats,
    fetchChartData,
    handleNewOrder,
    handleOrderStatusUpdated
  }
})
