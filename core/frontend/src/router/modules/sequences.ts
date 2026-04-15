import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/sequences',
  component: Layout,
  meta: {
    sort: 6,
    key: 'sequences',
    title: 'Sequences',
    icon: 'i-mdi-email-sync-outline',
  },
  children: [
    {
      path: '/sequences',
      name: 'Sequences',
      component: () => import('@/views/sequences/index.vue'),
    },
    {
      path: '/sequences/create',
      name: 'SequenceCreate',
      component: () => import('@/views/sequences/edit.vue'),
      meta: { hidden: true },
    },
    {
      path: '/sequences/:id/edit',
      name: 'SequenceEdit',
      component: () => import('@/views/sequences/edit.vue'),
      meta: { hidden: true },
    },
    {
      path: '/sequences/:id',
      name: 'SequenceDetail',
      component: () => import('@/views/sequences/detail.vue'),
      meta: { hidden: true },
    },
  ],
}

export default route
