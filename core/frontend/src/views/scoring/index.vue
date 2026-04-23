<template>
	<div class="scoring-page">
		<div class="bt-title">Lead Scoring</div>

		<!-- Stats Row -->
		<div class="stats-row">
			<score-stat-card
				v-for="stat in statCards"
				:key="stat.label"
				:label="stat.label"
				:count="stat.count"
				:color="stat.color"
				:icon="stat.icon"
				:active="filter === stat.filter"
				@click="filter = stat.filter"
			/>
			<n-card class="avg-card">
				<div class="avg-value">{{ stats.average_score?.toFixed(1) || '0' }}</div>
				<div class="avg-label">Avg Score</div>
			</n-card>
		</div>

		<!-- Actions -->
		<div class="actions-row">
			<n-space>
				<n-select
					v-model:value="filter"
					:options="filterOptions"
					placeholder="Filter by level"
					style="width: 160px"
					size="small"
				/>
				<n-button type="primary" size="small" @click="handleRecalculate" :loading="recalculating">
					Recalculate All
				</n-button>
			</n-space>
		</div>

		<!-- Leads Table -->
		<n-card>
			<n-data-table
				:columns="columns"
				:data="leads"
				:loading="loading"
				:pagination="pagination"
				:bordered="false"
				:row-key="(row: any) => row.id"
				remote
				@update:page="handlePageChange"
			/>
		</n-card>
	</div>
</template>

<script lang="ts" setup>
import { h, onMounted, ref, watch } from 'vue'
import { NTag, NProgress } from 'naive-ui'
import { Message } from '@/utils'
import { getScoringLeads, getScoringStats, recalculateScores } from '@/api/modules/scoring'
import ScoreStatCard from './components/ScoreStatCard.vue'

interface ScoringStats {
	cold_count: number
	warm_count: number
	hot_count: number
	converted_count: number
	average_score: number
	total_leads: number
}

interface LeadScore {
	id: number
	contact_email: string
	score: number
	engagement_level: string
	total_opens: number
	total_clicks: number
	total_replies: number
	total_bounces: number
	last_engagement_at: number
	last_scored_at: number
}

const stats = ref<ScoringStats>({
	cold_count: 0,
	warm_count: 0,
	hot_count: 0,
	converted_count: 0,
	average_score: 0,
	total_leads: 0,
})
const leads = ref<LeadScore[]>([])
const loading = ref(false)
const recalculating = ref(false)
const filter = ref('all')
const pagination = ref({ page: 1, pageSize: 20, itemCount: 0 })

const filterOptions = [
	{ label: 'All Leads', value: 'all' },
	{ label: 'Cold', value: 'cold' },
	{ label: 'Warm', value: 'warm' },
	{ label: 'Hot', value: 'hot' },
	{ label: 'Converted', value: 'converted' },
]

const levelConfig: Record<string, { color: string; type: 'success' | 'warning' | 'error' | 'info' | 'default' }> = {
	cold: { color: '#909399', type: 'default' },
	warm: { color: '#e6a23c', type: 'warning' },
	hot: { color: '#d03050', type: 'error' },
	converted: { color: '#18a058', type: 'success' },
}

const statCards = ref([
	{ label: 'Total Leads', count: 0, color: '#606266', icon: 'i-mdi-account-group-outline', filter: 'all' },
	{ label: 'Cold', count: 0, color: '#909399', icon: 'i-mdi-snowflake', filter: 'cold' },
	{ label: 'Warm', count: 0, color: '#e6a23c', icon: 'i-mdi-fire', filter: 'warm' },
	{ label: 'Hot', count: 0, color: '#d03050', icon: 'i-mdi-fire-alert', filter: 'hot' },
	{ label: 'Converted', count: 0, color: '#18a058', icon: 'i-mdi-check-circle-outline', filter: 'converted' },
])

const columns = [
	{
		title: 'Email',
		key: 'contact_email',
		ellipsis: true,
	},
	{
		title: 'Score',
		key: 'score',
		width: 120,
		render: (row: LeadScore) => {
			const pct = Math.min(row.score, 100)
			const cfg = levelConfig[row.engagement_level] || levelConfig.cold
			return h(NProgress, {
				type: 'line',
				percentage: pct,
				indicatorPlacement: 'inside',
				color: cfg.color,
				height: 18,
			})
		},
		sorter: true,
	},
	{
		title: 'Level',
		key: 'engagement_level',
		width: 100,
		render: (row: LeadScore) => {
			const cfg = levelConfig[row.engagement_level] || levelConfig.cold
			return h(NTag, { size: 'small', type: cfg.type, bordered: false }, () => row.engagement_level)
		},
	},
	{
		title: 'Opens',
		key: 'total_opens',
		width: 70,
	},
	{
		title: 'Clicks',
		key: 'total_clicks',
		width: 70,
	},
	{
		title: 'Replies',
		key: 'total_replies',
		width: 70,
	},
	{
		title: 'Bounces',
		key: 'total_bounces',
		width: 70,
	},
]

async function fetchStats() {
	try {
		const res = await getScoringStats()
		if (res) Object.assign(stats.value, res)
		statCards.value[0].count = res.total_leads
		statCards.value[1].count = res.cold_count
		statCards.value[2].count = res.warm_count
		statCards.value[3].count = res.hot_count
		statCards.value[4].count = res.converted_count
	} catch (e) {
		console.error('Failed to fetch scoring stats', e)
	}
}

async function fetchLeads() {
	loading.value = true
	try {
		const res = await getScoringLeads({
			level: filter.value,
			page: pagination.value.page,
			page_size: pagination.value.pageSize,
		})
		if (res) {
			leads.value = res.list || []
			pagination.value.itemCount = res.total || 0
		}
	} catch (e) {
		console.error('Failed to fetch leads', e)
	} finally {
		loading.value = false
	}
}

async function handleRecalculate() {
	recalculating.value = true
	try {
		await recalculateScores()
		Message.success('Scores recalculated')
		await fetchStats()
		await fetchLeads()
	} catch (e) {
		Message.error('Recalculation failed')
	} finally {
		recalculating.value = false
	}
}

function handlePageChange(page: number) {
	pagination.value.page = page
	fetchLeads()
}

watch(filter, () => {
	pagination.value.page = 1
	fetchLeads()
})

onMounted(() => {
	fetchStats()
	fetchLeads()
})
</script>

<style lang="scss" scoped>
.scoring-page {
	padding: 20px;
}

.stats-row {
	display: flex;
	gap: 16px;
	margin-bottom: 16px;
}

.avg-card {
	text-align: center;
	min-width: 100px;

	.avg-value {
		font-size: 24px;
		font-weight: 700;
	}

	.avg-label {
		font-size: 12px;
		color: var(--n-text-color-3);
		margin-top: 2px;
	}
}

.actions-row {
	margin-bottom: 16px;
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>
