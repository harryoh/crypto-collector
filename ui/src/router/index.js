import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/redirect',
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirectView')
      }
    ]
  },
  {
    path: '/',
    redirect: '/premium'
  },
  {
    path: '/premium',
    name: 'PremiumMonitor',
    component: () => import('@/views/PremiumMonitor'),
    meta: { title: 'Preminum Monitor', icon: 'dashboard', affix: true }
  }
]

const router = new VueRouter({
  routes
})

export default router
