import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/chat',
      name: 'chat',
      component: () => import('@/views/ChatView.vue')
    },
    {
      path: '/dashboard/:modelId',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue')
    },
    {
      path: '/microscope/:modelId',
      name: 'microscope',
      component: () => import('@/views/MicroscopeView.vue')
    },
    {
      path: '/attention/:modelId',
      name: 'attention',
      component: () => import('@/views/AttentionView.vue')
    },
    {
      path: '/activation/:modelId',
      name: 'activation',
      component: () => import('@/views/ActivationView.vue')
    }
  ]
})

export default router
