import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/analytics',
	component: Layout,
	meta: { sort: 0, key: 'analytics', title: 'Analytics' },
	children: [
		{
			path: '/analytics',
			name: 'Analytics',
			component: () => import('@/views/analytics/index.vue'),
		},
	],
}

export default route
