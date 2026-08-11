import { createRouter, createWebHistory } from 'vue-router'

import HomePage from '../views/HomePage.vue'
import LoginPage from '../views/LoginPage.vue'
import AdminDashboard from '../views/AdminDashboard.vue'
import ShopPage from '../views/ShopPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/shop',
    name: 'Shop',
    component: ShopPage,
    meta: { requiresAuth: true, role: 'customer' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    return savedPosition || { top: 0 }
  }
})


router.beforeEach((to, _from, next) => {
  const user = JSON.parse(localStorage.getItem('openpeo_user') || 'null')

  if (to.meta.requiresAuth && !user) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.role && user && user.role !== to.meta.role) {
    
    if (user.role === 'admin') {
      next({ name: 'AdminDashboard' })
    } else {
      next({ name: 'Shop' })
    }
    return
  }

  
  if (to.name === 'Login' && user) {
    if (user.role === 'admin') {
      next({ name: 'AdminDashboard' })
    } else {
      next({ name: 'Shop' })
    }
    return
  }

  next()
})

export default router
