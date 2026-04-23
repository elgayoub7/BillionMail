<template>
	<n-data-table
		:columns="columns"
		:data="data"
		:loading="loading"
		:pagination="{ pageSize: 5 }"
		:bordered="false"
		size="small"
	/>
</template>

<script lang="ts" setup>
import { h } from 'vue'
import { NTag } from 'naive-ui'

defineProps<{
	data: Array<{
		id: number
		name: string
		status: number
		enrolled: number
		completed: number
		reply_rate: number
		open_rate: number
		bounce_rate: number
		has_ab_test: boolean
	}>
	loading: boolean
}>()

const columns = [
	{
		title: 'Sequence',
		key: 'name',
		ellipsis: true,
	},
	{
		title: 'Status',
		key: 'status',
		width: 100,
		render: (row: any) => {
			const map: Record<number, { label: string; type: 'success' | 'warning' | 'info' | 'default' }> = {
				0: { label: 'Draft', type: 'default' },
				1: { label: 'Active', type: 'success' },
				2: { label: 'Paused', type: 'warning' },
			}
			const s = map[row.status] || { label: 'Unknown', type: 'default' as const }
			return h(NTag, { size: 'small', type: s.type, bordered: false }, () => s.label)
		},
	},
	{
		title: 'Enrolled',
		key: 'enrolled',
		width: 80,
	},
	{
		title: 'Open %',
		key: 'open_rate',
		width: 80,
		render: (row: any) => row.open_rate?.toFixed(1) + '%',
	},
	{
		title: 'Reply %',
		key: 'reply_rate',
		width: 80,
		render: (row: any) => row.reply_rate?.toFixed(1) + '%',
	},
	{
		title: 'AB Test',
		key: 'has_ab_test',
		width: 75,
		render: (row: any) =>
			row.has_ab_test
				? h(NTag, { size: 'small', type: 'info', bordered: false }, () => 'Active')
				: h('span', { style: 'color: var(--n-text-color-3)' }, '—'),
	},
]
</script>
