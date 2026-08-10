import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/koleksi',
      name: 'koleksi',
      component: () => import('../views/KoleksiView.vue')
    },
    {
      path: '/tentang',
      name: 'tentang',
      component: () => import('../views/TentangView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/register-seller',
      name: 'register-seller',
      component: () => import('../views/RegisterSellerView.vue')
    },
    {
      path: '/seller/dashboard',
      name: 'seller-dashboard',
      component: () => import('../views/seller/SellerDashboardView.vue'),
      meta: { requiresAuth: true, requiresSeller: true }
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: () => import('../views/CheckoutView.vue'),
      meta: { requiresAuth: true }
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return savedPosition || { top: 0 }
  }
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresSeller && auth.user?.role !== 'seller') {
    return { name: 'home' }
  }
})

export default router
