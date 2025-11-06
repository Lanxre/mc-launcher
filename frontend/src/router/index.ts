import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import ModView from '@/views/ModView.vue'
import DownloadsView from '@/views/DownloadsView.vue'

const routes = [
  {
    path: '/',
    name: 'main',
    component: HomeView,
  },
  {
    path: '/:name',
    name: 'mod',
    component: ModView,
  },
  {
    path: '/mod-downloads',
    name: 'mod-downloads',
    component: DownloadsView,
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router