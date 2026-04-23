import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/domain-health',
	component: Layout,
	meta: {
		sort: 9,
		key: 'domain-health',
		title: 'Domain Health',
		icon: 'i-mdi-shield-check-outline',
	},
	children: [
		{
			path: '/domain-health',
			name: 'DomainHealth',
			component: () => import('@/views/domain-health/index.vue'),
		},
	],
}

export default route
