<template>
	<div class="dashboard">
		<div class="dashboard-header">
			<h1 class="dashboard-title">Dashboard</h1>
			<filter-bar
				v-model:domain="domain"
				v-model:date-range="dateRange"
				@update:domain="handleDataUpdate"
				@update:date-range="handleDataUpdate" />
		</div>

		<!-- KPI Cards -->
		<div class="kpi-grid">
			<metric-card
				v-for="(item, key) in kpiCards"
				:key="key"
				:title="item.label"
				:value="item.value"
				:unit="item.unit"
				:icon="item.icon"
				:onClick="item.key === 'delayed_queue' ? onClickDelayedQueue : undefined" />
		</div>

		<!-- Main Chart -->
		<div class="chart-section">
			<div class="card">
				<div class="card-header">
					<span class="card-title">Sending Trends</span>
				</div>
				<send-today-stats :data="sendMail" @show-fail="handleShowFail" />
			</div>
		</div>

		<!-- Rate Charts -->
		<div class="rate-section">
			<rate-chart-panel :bounce="bounceRate" :click="clickRate" :open="openRate" />
		</div>

		<!-- Bottom Row -->
		<div class="bottom-row">
			<div class="card bottom-card">
				<div class="card-header">
					<span class="card-title">Mail Providers</span>
				</div>
				<provider-table v-model:value="providers" />
			</div>
		</div>

		<fail-modal />
	</div>
</template>

<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core'
import { useThemeVars } from 'naive-ui'
import { getDayTimeRange, isArray, isObject } from '@/utils'
import { useModal } from '@/hooks/modal/useModal'
import { getOverviewInfo } from '@/api/modules/overview'
import type { MailOverview, MailProvider, RateData } from './types'

import FilterBar from './components/FilterBar.vue'
import MetricCard from './components/MetricCard.vue'
import ProviderTable from './components/ProviderTable.vue'
import SendTodayStats from './components/SendTodayStats.vue'
import RateChartPanel from './components/RateChartPanel.vue'
import SendFailDetails from './components/SendFailDetails.vue'

const { t } = useI18n()

const theme = useThemeVars()

const domain = ref('')

const dateRange = ref(getDayTimeRange())

const providers = ref<MailProvider[]>([])

const rateData = reactive<RateData>({
	delivery_rate: { label: t('overview.delivered'), value: 0, unit: '%' },
	open_rate: { label: t('overview.opened'), value: 0, unit: '%' },
	click_rate: { label: t('overview.clicked'), value: 0, unit: '%' },
	bounce_rate: { label: t('overview.bounced'), value: 0, unit: '%' },
})

const kpiCards = computed(() => [
	{ key: 'delivery_rate', ...rateData.delivery_rate, icon: 'i-mdi:send-check-outline' },
	{ key: 'open_rate', ...rateData.open_rate, icon: 'i-mdi:email-open-outline' },
	{ key: 'click_rate', ...rateData.click_rate, icon: 'i-mdi:cursor-default-click-outline' },
	{ key: 'bounce_rate', ...rateData.bounce_rate, icon: 'i-mdi:email-alert-outline' },
])

const delayedQueue = ref(0)

const sendMail = ref<MailOverview['send_mail_chart']>({
	column_type: 'hourly',
	dashboard: {
		delivered: 0,
		delivery_rate: 0,
		failed: 0,
		failure_rate: 0,
		sends: 0,
	},
	data: [],
})

const bounceRate = ref<MailOverview['bounce_rate_chart']>({
	column_type: 'hourly',
	data: [],
})

const clickRate = ref<MailOverview['click_rate_chart']>({
	column_type: 'hourly',
	data: [],
})

const openRate = ref<MailOverview['open_rate_chart']>({
	column_type: 'hourly',
	data: [],
})

const router = useRouter()

const onClickDelayedQueue = () => {
	router.push('/settings/send-queue')
}

const [FailModal, failModalApi] = useModal({
	component: SendFailDetails,
})

const handleShowFail = () => {
	failModalApi.setState({
		domain: domain.value,
		startTime: dateRange.value[0],
		endTime: dateRange.value[1],
	})
	failModalApi.open()
}

const handleDataUpdate = useDebounceFn(fetchOverviewData, 300)

const updateRateData = (dashboard: MailOverview['dashboard']) => {
	Object.entries(dashboard).forEach(([key, value]) => {
		if (key in rateData) {
			rateData[key].value = value
		}
	})
	delayedQueue.value = dashboard.delayed_queue
}

async function fetchOverviewData() {
	const res = await getOverviewInfo({
		domain: domain.value,
		start_time: Math.floor(dateRange.value[0] / 1000),
		end_time: Math.floor(dateRange.value[1] / 1000),
	})

	if (isObject<MailOverview>(res)) {
		updateRateData(res.dashboard)
		providers.value = isArray(res.mail_providers) ? res.mail_providers : []
		sendMail.value = res.send_mail_chart
		bounceRate.value = res.bounce_rate_chart
		clickRate.value = res.click_rate_chart
		openRate.value = res.open_rate_chart
	}
}

onMounted(() => {
	fetchOverviewData()
})
</script>

<style lang="scss" scoped>
.dashboard {
	padding: 24px;
}

.dashboard-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24px;
}

.dashboard-title {
	margin: 0;
	font-size: 22px;
	font-weight: 700;
	color: #e2e8f0;
}

.kpi-grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 16px;
	margin-bottom: 24px;
}

.card {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
}

.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
}

.card-title {
	font-size: 15px;
	font-weight: 600;
	color: #e2e8f0;
}

.chart-section {
	margin-bottom: 24px;
}

.rate-section {
	margin-bottom: 24px;

	:deep(.n-card) {
		background: #1a1d27;
		border: 1px solid #2e3142;
		border-radius: 12px;
	}
}

.bottom-row {
	display: grid;
	grid-template-columns: 1fr;
	gap: 16px;
}

.bottom-card {
	:deep(.n-data-table) {
		--n-td-color: transparent;
		--n-th-color: transparent;
		--n-border-color: #2e3142;
	}
}

@media (max-width: 1200px) {
	.kpi-grid {
		grid-template-columns: repeat(2, 1fr);
	}
}
</style>
