import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import ModView from '@/views/ModView.vue'

const routes = [
  {
    path: '/',
    name: 'main',
    component: HomeView,
  },
  {
    path: '/mod',
    name: 'mod',
    component: ModView,
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router