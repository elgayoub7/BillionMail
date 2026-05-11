import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/deliverability',
	component: Layout,
	meta: {
		sort: 3,
		key: 'deliverability',
		title: 'Inbox Placement',
	},
	children: [
		{
			path: '/deliverability',
			name: 'InboxPlacement',
			component: () => import('@/views/deliverability/InboxPlacement.vue'),
		},
	],
}

export default route
