import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/cold-dashboard',
	component: Layout,
	meta: {
		sort: 1,
		key: 'cold-dashboard',
		title: 'Cold Dashboard',
		icon: 'i-mdi-monitor-dashboard',
	},
	children: [
		{
			path: '/cold-dashboard',
			name: 'ColdDashboard',
			component: () => import('@/views/cold-dashboard/index.vue'),
		},
	],
}

export default route
