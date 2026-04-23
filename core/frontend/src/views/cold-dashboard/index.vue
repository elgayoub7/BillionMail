<template>
	<div class="cold-dashboard">
		<div class="bt-title">Cold Dashboard</div>

		<!-- KPI Cards -->
		<div class="kpi-row">
			<n-card v-for="kpi in kpiCards" :key="kpi.label" class="kpi-card">
				<div class="kpi-content">
					<div class="kpi-value" :style="{ color: kpi.color || 'var(--n-text-color)' }">
						{{ kpi.value }}
					</div>
					<div class="kpi-label">{{ kpi.label }}</div>
				</div>
			</n-card>
		</div>

		<!-- Two columns: Active Sequences + Alerts -->
		<div class="content-row">
			<n-card title="Active Sequences" class="flex-1">
				<template #header-extra>
					<n-button text type="primary" @click="router.push('/sequences')">
						View All
					</n-button>
				</template>
				<active-sequences-table :data="sequences" :loading="loadingSequences" />
			</n-card>

			<n-card title="Recent Alerts" class="flex-1">
				<template #header-extra>
					<n-tag :type="alertCount > 0 ? 'warning' : 'success'" size="small">
						{{ alertCount }} alerts
					</n-tag>
				</template>
				<alerts-list :alerts="alerts" :loading="loadingAlerts" />
			</n-card>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardStats, getActiveSequences, getAlerts } from '@/api/modules/cold-dashboard'
import ActiveSequencesTable from './components/ActiveSequencesTable.vue'
import AlertsList from './components/AlertsList.vue'

const router = useRouter()

interface DashboardStats {
	emails_sent_today: number
	emails_sent_week: number
	emails_sent_month: number
	open_rate: number
	click_rate: number
	reply_rate: number
	bounce_rate: number
	active_sequences: number
	active_warmups: number
}

interface ActiveSequence {
	id: number
	name: string
	status: number
	enrolled: number
	completed: number
	reply_rate: number
	open_rate: number
	bounce_rate: number
	has_ab_test: boolean
}

interface Alert {
	type: string
	message: string
	time: number
}

const stats = ref<DashboardStats>({
	emails_sent_today: 0,
	emails_sent_week: 0,
	emails_sent_month: 0,
	open_rate: 0,
	click_rate: 0,
	reply_rate: 0,
	bounce_rate: 0,
	active_sequences: 0,
	active_warmups: 0,
})

const sequences = ref<ActiveSequence[]>([])
const alerts = ref<Alert[]>([])
const loadingSequences = ref(false)
const loadingAlerts = ref(false)

const alertCount = computed(() => alerts.value.length)

const kpiCards = computed(() => [
	{ label: 'Sent Today', value: stats.value.emails_sent_today.toLocaleString(), color: undefined },
	{ label: 'Sent This Week', value: stats.value.emails_sent_week.toLocaleString(), color: undefined },
	{ label: 'Open Rate', value: stats.value.open_rate.toFixed(1) + '%', color: stats.value.open_rate > 20 ? '#18a058' : '#d03050' },
	{ label: 'Click Rate', value: stats.value.click_rate.toFixed(1) + '%', color: stats.value.click_rate > 5 ? '#18a058' : '#d03050' },
	{ label: 'Reply Rate', value: stats.value.reply_rate.toFixed(1) + '%', color: stats.value.reply_rate > 3 ? '#18a058' : '#d03050' },
	{ label: 'Bounce Rate', value: stats.value.bounce_rate.toFixed(1) + '%', color: stats.value.bounce_rate < 5 ? '#18a058' : '#d03050' },
])

async function fetchStats() {
	try {
		const res = await getDashboardStats()
		if (res) Object.assign(stats.value, res)
	} catch (e) {
		console.error('Failed to fetch dashboard stats', e)
	}
}

async function fetchSequences() {
	loadingSequences.value = true
	try {
		const res = await getActiveSequences()
		if (res) sequences.value = Array.isArray(res) ? res : []
	} catch (e) {
		console.error('Failed to fetch sequences', e)
	} finally {
		loadingSequences.value = false
	}
}

async function fetchAlerts() {
	loadingAlerts.value = true
	try {
		const res = await getAlerts({ limit: 20 })
		if (res) alerts.value = Array.isArray(res) ? res : []
	} catch (e) {
		console.error('Failed to fetch alerts', e)
	} finally {
		loadingAlerts.value = false
	}
}

onMounted(() => {
	fetchStats()
	fetchSequences()
	fetchAlerts()
})
</script>

<style lang="scss" scoped>
.cold-dashboard {
	padding: 20px;
}

.kpi-row {
	display: grid;
	grid-template-columns: repeat(6, 1fr);
	gap: 16px;
	margin-bottom: 20px;
}

.kpi-card {
	.kpi-content {
		text-align: center;
	}

	.kpi-value {
		font-size: 24px;
		font-weight: 700;
		line-height: 1.2;
	}

	.kpi-label {
		font-size: 13px;
		color: var(--n-text-color-3);
		margin-top: 4px;
	}
}

.content-row {
	display: flex;
	gap: 20px;

	.flex-1 {
		flex: 1;
	}
}
</style>
