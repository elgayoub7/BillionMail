import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/abtest',
	component: Layout,
	meta: {
		sort: 8,
		key: 'abtest',
		title: 'AB Tests',
		icon: 'i-mdi-ab-testing',
	},
	children: [
		{
			path: '/abtest',
			name: 'AbTest',
			component: () => import('@/views/abtest/index.vue'),
		},
	],
}

export default route
