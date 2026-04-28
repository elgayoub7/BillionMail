import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/inbox',
	component: Layout,
	meta: {
		sort: 2,
		key: 'inbox',
		title: 'Inbox',
	},
	children: [
		{
			path: '/inbox',
			name: 'Inbox',
			component: () => import('@/views/inbox/index.vue'),
		},
	],
}

export default route
