import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/scoring',
	component: Layout,
	meta: {
		sort: 7,
		key: 'scoring',
		title: 'Lead Scoring',
		icon: 'i-mdi-star-circle-outline',
	},
	children: [
		{
			path: '/scoring',
			name: 'Scoring',
			component: () => import('@/views/scoring/index.vue'),
		},
	],
}

export default route
