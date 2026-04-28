<template>
	<div class="sequences-page">
		<div class="page-header">
			<h1 class="page-title">Campaigns</h1>
			<div class="header-actions">
				<div class="filter-group">
					<n-select
						v-model:value="tableParams.status"
						:options="statusOptions"
						class="status-select"
						@update:value="resetTable" />
					<bt-search
						v-model:value="tableParams.keyword"
						placeholder="Search campaigns..."
						@search="resetTable" />
				</div>
				<n-button type="primary" @click="handleCreate">
					<template #icon><i class="i-mdi:plus"></i></template>
					Create Campaign
				</n-button>
			</div>
		</div>

		<!-- Cards Grid -->
		<div v-if="tableProps.data?.length" class="cards-grid">
			<div
				v-for="row in tableProps.data"
				:key="row.id"
				class="campaign-card"
				@click="router.push(`/sequences/${row.id}`)">
				<div class="card-top">
					<div class="card-info">
						<h3 class="card-name">{{ row.name }}</h3>
						<span v-if="row.description" class="card-desc">{{ row.description }}</span>
					</div>
					<n-tag :type="getStatusType(row.status)" size="small" round>
						{{ getStatusLabel(row.status) }}
					</n-tag>
				</div>

				<div class="card-stats">
					<div class="stat-row">
						<div class="stat">
							<span class="stat-value">{{ row.total_enrolled }}</span>
							<span class="stat-label">Enrolled</span>
						</div>
						<div class="stat">
							<span class="stat-value">{{ row.total_completed }}</span>
							<span class="stat-label">Completed</span>
						</div>
						<div class="stat">
							<span class="stat-value">{{ row.step_count }}</span>
							<span class="stat-label">Steps</span>
						</div>
						<div v-if="row.total_bounced > 0" class="stat stat-danger">
							<span class="stat-value">{{ row.total_bounced }}</span>
							<span class="stat-label">Bounced</span>
						</div>
					</div>
				</div>

				<div class="card-progress">
					<div class="progress-bar">
						<div
							class="progress-fill"
							:style="{ width: getProgress(row) + '%' }"></div>
					</div>
					<span class="progress-text">{{ getProgress(row) }}% complete</span>
				</div>

				<div class="card-actions" @click.stop>
					<n-button quaternary size="small" @click="router.push(`/sequences/${row.id}`)">
						<i class="i-mdi:eye-outline"></i>
					</n-button>
					<n-button quaternary size="small" @click="router.push(`/sequences/${row.id}/edit`)">
						<i class="i-mdi:pencil-outline"></i>
					</n-button>
					<n-button
						v-if="row.status === 0"
						quaternary
						size="small"
						type="success"
						@click="handleActivate(row.id)">
						<i class="i-mdi:play-outline"></i>
					</n-button>
					<n-button
						v-if="row.status === 1"
						quaternary
						size="small"
						type="warning"
						@click="handlePause(row.id)">
						<i class="i-mdi:pause-outline"></i>
					</n-button>
					<n-button
						v-if="row.status === 2"
						quaternary
						size="small"
						type="success"
						@click="handleResume(row.id)">
						<i class="i-mdi:play-outline"></i>
					</n-button>
					<n-popconfirm @positive-click="handleDelete(row)">
						<template #trigger>
							<n-button quaternary size="small" type="error">
								<i class="i-mdi:delete-outline"></i>
							</n-button>
						</template>
						Delete this campaign?
					</n-popconfirm>
				</div>
			</div>
		</div>

		<!-- Empty State -->
		<div v-else class="empty-state">
			<i class="i-mdi:rocket-launch-outline empty-icon"></i>
			<h3>No campaigns yet</h3>
			<p>Create your first outreach campaign to get started</p>
			<n-button type="primary" @click="handleCreate">
				<template #icon><i class="i-mdi:plus"></i></template>
				Create Campaign
			</n-button>
		</div>

		<!-- Pagination -->
		<div v-if="pageProps.total > 0" class="pagination-bar">
			<bt-table-page v-bind="pageProps" @refresh="fetchTable" />
		</div>
	</div>
</template>

<script lang="ts" setup>
import { useDataTable } from '@/hooks/useDataTable'
import { getSequenceList, deleteSequence, activateSequence, pauseSequence, resumeSequence } from '@/api/modules/sequences/sequence'
import { Message } from '@/utils'
import type { Sequence, SequenceParams } from './interface'

const { t } = useI18n()
const router = useRouter()

const statusOptions = [
	{ label: 'All Status', value: -1 },
	{ label: 'Draft', value: 0 },
	{ label: 'Active', value: 1 },
	{ label: 'Paused', value: 2 },
	{ label: 'Archived', value: 3 },
]

const statusMap: Record<number, { label: string; type: 'default' | 'success' | 'warning' | 'error' }> = {
	0: { label: 'Draft', type: 'default' },
	1: { label: 'Active', type: 'success' },
	2: { label: 'Paused', type: 'warning' },
	3: { label: 'Archived', type: 'error' },
}

const getStatusType = (status: number) => statusMap[status]?.type || 'default'
const getStatusLabel = (status: number) => statusMap[status]?.label || 'Unknown'

const getProgress = (row: Sequence) => {
	if (row.total_enrolled === 0) return 0
	return Math.round((row.total_completed / row.total_enrolled) * 100)
}

const { tableParams, tableProps, pageProps, fetchTable, resetTable } = useDataTable<Sequence, SequenceParams>({
	params: {
		page: 1,
		page_size: 12,
		keyword: '',
		status: -1,
	},
	useParams: (params) => params,
	fetchFn: getSequenceList as any,
})

function handleCreate() {
	router.push('/sequences/create')
}

async function handleActivate(id: number) {
	await activateSequence({ id })
	fetchTable()
}

async function handlePause(id: number) {
	await pauseSequence({ id })
	fetchTable()
}

async function handleResume(id: number) {
	await resumeSequence({ id })
	fetchTable()
}

async function handleDelete(row: Sequence) {
	await deleteSequence({ id: row.id })
	Message.success('Campaign deleted')
	fetchTable()
}
</script>

<style lang="scss" scoped>
.sequences-page {
	padding: 24px;
}

.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24px;
}

.page-title {
	margin: 0;
	font-size: 22px;
	font-weight: 700;
	color: #e2e8f0;
}

.header-actions {
	display: flex;
	gap: 12px;
	align-items: center;
}

.filter-group {
	display: flex;
	gap: 10px;
}

.status-select {
	width: 140px;
}

.cards-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
	gap: 16px;
}

.campaign-card {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
	cursor: pointer;
	transition: all 0.2s ease;

	&:hover {
		border-color: #6c5ce7;
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(108, 92, 231, 0.1);
	}
}

.card-top {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 16px;
}

.card-info {
	flex: 1;
	min-width: 0;
}

.card-name {
	margin: 0 0 4px;
	font-size: 16px;
	font-weight: 600;
	color: #e2e8f0;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.card-desc {
	font-size: 13px;
	color: #6b7280;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	display: block;
}

.card-stats {
	margin-bottom: 16px;
}

.stat-row {
	display: flex;
	gap: 24px;
}

.stat {
	display: flex;
	flex-direction: column;
	gap: 2px;

	.stat-value {
		font-size: 18px;
		font-weight: 700;
		color: #e2e8f0;
	}

	.stat-label {
		font-size: 11px;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	&.stat-danger .stat-value {
		color: #ef4444;
	}
}

.card-progress {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 16px;
}

.progress-bar {
	flex: 1;
	height: 6px;
	background: #2e3142;
	border-radius: 3px;
	overflow: hidden;
}

.progress-fill {
	height: 100%;
	background: linear-gradient(90deg, #6c5ce7, #a78bfa);
	border-radius: 3px;
	transition: width 0.3s ease;
}

.progress-text {
	font-size: 12px;
	color: #6b7280;
	white-space: nowrap;
}

.card-actions {
	display: flex;
	gap: 4px;
	justify-content: flex-end;
	border-top: 1px solid #2e3142;
	padding-top: 12px;
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 80px 20px;
	color: #6b7280;

	h3 {
		font-size: 18px;
		color: #8892a8;
		margin: 16px 0 8px;
	}

	p {
		font-size: 14px;
		margin: 0 0 24px;
	}
}

.empty-icon {
	font-size: 48px;
	color: #2e3142;
}

.pagination-bar {
	display: flex;
	justify-content: flex-end;
	margin-top: 24px;
}
</style>
