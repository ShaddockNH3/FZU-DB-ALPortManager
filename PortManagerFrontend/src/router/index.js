import { createRouter, createWebHistory } from 'vue-router'
import { PhChartBar, PhListDashes } from '@phosphor-icons/vue'

// 导航菜单直接从这份路由配置生成（看 Layout.vue），加页面只加一条路由。
const routes = [
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/views/Layout.vue'),
    redirect: '/ships',
    children: [
      {
        path: '/ships',
        name: 'ShipList',
        component: () => import('@/views/ShipList.vue'),
        meta: { title: '舰船管理', icon: PhListDashes }
      },
      {
        path: '/statistics',
        name: 'Statistics',
        component: () => import('@/views/Statistics.vue'),
        meta: { title: '数据统计', icon: PhChartBar }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 港区管理系统` : '港区管理系统'
  next()
})

export default router
