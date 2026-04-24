<template>
	<div class="dashboard">
		<div class="dash-header">
			<h1 class="dash-title">Dashboard</h1>
			<span class="dash-date">{{ today }}</span>
		</div>
		<div class="kpi-grid">
			<div v-for="kpi in kpiCards" :key="kpi.label" class="kpi-card">
				<div class="kpi-icon" :style="{ background: kpi.bg }">
					<span :style="{ color: kpi.accent }">{{ kpi.emoji }}</span>
				</div>
				<div class="kpi-info">
					<div class="kpi-value" :style="{ color: kpi.color || '#f1f3f7' }">{{ kpi.value }}</div>
					<div class="kpi-label">{{ kpi.label }}</div>
				</div>
			</div>
		</div>
		<div class="content-grid">
			<n-card title="Active Sequences" class="dark-card flex-1">
				<template #header-extra>
					<n-button text type="primary" @click="router.push('/sequences')">View All</n-button>
				</template>
				<active-sequences-table :data="sequences" :loading="loadingSequences" />
			</n-card>
			<n-card title="Recent Alerts" class="dark-card flex-1">
				<template #header-extra>
					<n-tag :type="alertCount > 0 ? 'warning' : 'success'" size="small">{{ alertCount }} alerts</n-tag>
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
const today = new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

interface DashboardStats { emails_sent_today: number; emails_sent_week: number; emails_sent_month: number; open_rate: number; click_rate: number; reply_rate: number; bounce_rate: number; active_sequences: number; active_warmups: number }
interface ActiveSequence { id: number; name: string; status: number; enrolled: number; completed: number; reply_rate: number; open_rate: number; bounce_rate: number; has_ab_test: boolean }
interface Alert { type: string; message: string; time: number }

const stats = ref<DashboardStats>({ emails_sent_today: 0, emails_sent_week: 0, emails_sent_month: 0, open_rate: 0, click_rate: 0, reply_rate: 0, bounce_rate: 0, active_sequences: 0, active_warmups: 0 })
const sequences = ref<ActiveSequence[]>([])
const alerts = ref<Alert[]>([])
const loadingSequences = ref(false)
const loadingAlerts = ref(false)
const alertCount = computed(() => alerts.value.length)

const kpiCards = computed(() => [
	{ label: 'Sent Today', value: stats.value.emails_sent_today.toLocaleString(), emoji: '✈', accent: '#6c5ce7', bg: 'rgba(108,92,231,0.1)', color: undefined },
	{ label: 'This Week', value: stats.value.emails_sent_week.toLocaleString(), emoji: '📊', accent: '#6c5ce7', bg: 'rgba(108,92,231,0.1)', color: undefined },
	{ label: 'Open Rate', value: stats.value.open_rate.toFixed(1) + '%', emoji: '📧', accent: stats.value.open_rate > 20 ? '#00d68f' : '#ff3d71', bg: stats.value.open_rate > 20 ? 'rgba(0,214,143,0.1)' : 'rgba(255,61,113,0.1)', color: stats.value.open_rate > 20 ? '#00d68f' : '#ff3d71' },
	{ label: 'Click Rate', value: stats.value.click_rate.toFixed(1) + '%', emoji: '🔥', accent: stats.value.click_rate > 5 ? '#00d68f' : '#ff3d71', bg: stats.value.click_rate > 5 ? 'rgba(0,214,143,0.1)' : 'rgba(255,61,113,0.1)', color: stats.value.click_rate > 5 ? '#00d68f' : '#ff3d71' },
	{ label: 'Reply Rate', value: stats.value.reply_rate.toFixed(1) + '%', emoji: '💬', accent: stats.value.reply_rate > 3 ? '#00d68f' : '#ffaa00', bg: stats.value.reply_rate > 3 ? 'rgba(0,214,143,0.1)' : 'rgba(255,170,0,0.1)', color: stats.value.reply_rate > 3 ? '#00d68f' : '#ffaa00' },
	{ label: 'Bounce Rate', value: stats.value.bounce_rate.toFixed(1) + '%', emoji: '⚠', accent: stats.value.bounce_rate < 5 ? '#00d68f' : '#ff3d71', bg: stats.value.bounce_rate < 5 ? 'rgba(0,214,143,0.1)' : 'rgba(255,61,113,0.1)', color: stats.value.bounce_rate < 5 ? '#00d68f' : '#ff3d71' },
])

async function fetchStats() { try { const r = await getDashboardStats(); if (r) Object.assign(stats.value, r) } catch {} }
async function fetchSequences() { loadingSequences.value = true; try { const r = await getActiveSequences(); if (r) sequences.value = Array.isArray(r) ? r : [] } catch {} finally { loadingSequences.value = false } }
async function fetchAlerts() { loadingAlerts.value = true; try { const r = await getAlerts({ limit: 20 }); if (r) alerts.value = Array.isArray(r) ? r : [] } catch {} finally { loadingAlerts.value = false } }

onMounted(() => { fetchStats(); fetchSequences(); fetchAlerts() })
</script>
<style scoped>
.dashboard { padding: 24px; }
.dash-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.dash-title { font-size: 24px; font-weight: 700; color: #f1f3f7; margin: 0; }
.dash-date { font-size: 14px; color: #6b7084; }
.kpi-grid { display: grid; grid-template-columns: repeat(6, 1fr); gap: 16px; margin-bottom: 24px; }
.kpi-card {
	display: flex; align-items: center; gap: 12px; padding: 16px;
	background: #1a1d27; border-radius: 12px; border: 1px solid #2e3142;
	transition: all 0.2s ease;
}
.kpi-card:hover { border-color: #363a4a; transform: translateY(-1px); }
.kpi-icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 18px; flex-shrink: 0; }
.kpi-value { font-size: 20px; font-weight: 700; line-height: 1.2; }
.kpi-label { font-size: 12px; color: #6b7084; margin-top: 2px; }
.content-grid { display: flex; gap: 20px; }
.dark-card { background: #1a1d27 !important; border: 1px solid #2e3142 !important; }
:deep(.n-card .n-card-header) { color: #f1f3f7; }
:deep(.n-card) { background: #1a1d27; border-color: #2e3142; }
</style>
