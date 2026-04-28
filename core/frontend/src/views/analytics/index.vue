<template>
	<div class="analytics-page">
		<div class="page-header">
			<h1 class="page-title">Analytics</h1>
			<n-select v-model:value="period" :options="periodOptions" style="width:160px" @update:value="refresh" />
		</div>

		<!-- KPI Row -->
		<div class="kpi-grid">
			<div class="kpi-card">
				<span class="kpi-icon" style="color:#6c5ce7"><i class="i-mdi:send-outline"></i></span>
				<div>
					<div class="kpi-val">{{ stats.delivery_rate?.toFixed(1) || '0' }}%</div>
					<div class="kpi-lbl">Delivery Rate</div>
				</div>
			</div>
			<div class="kpi-card">
				<span class="kpi-icon" style="color:#00d68f"><i class="i-mdi:email-open-outline"></i></span>
				<div>
					<div class="kpi-val">{{ stats.open_rate?.toFixed(1) || '0' }}%</div>
					<div class="kpi-lbl">Open Rate</div>
				</div>
			</div>
			<div class="kpi-card">
				<span class="kpi-icon" style="color:#0095ff"><i class="i-mdi:cursor-default-click-outline"></i></span>
				<div>
					<div class="kpi-val">{{ stats.click_rate?.toFixed(1) || '0' }}%</div>
					<div class="kpi-lbl">Click Rate</div>
				</div>
			</div>
			<div class="kpi-card">
				<span class="kpi-icon" style="color:#ff3d71"><i class="i-mdi:email-alert-outline"></i></span>
				<div>
					<div class="kpi-val">{{ stats.bounce_rate?.toFixed(1) || '0' }}%</div>
					<div class="kpi-lbl">Bounce Rate</div>
				</div>
			</div>
		</div>

		<!-- Charts Row -->
		<div class="charts-row">
			<div class="chart-card">
				<div class="card-title">Sending Trends</div>
				<div ref="sendChartRef" class="chart-container"></div>
			</div>
			<div class="chart-card">
				<div class="card-title">Bounce Rate</div>
				<div ref="bounceChartRef" class="chart-container"></div>
			</div>
		</div>

		<!-- Mail Providers -->
		<div class="providers-section">
			<div class="card-title">Mail Providers Performance</div>
			<n-data-table :columns="providerColumns" :data="providers" :bordered="false" size="small" />
		</div>
	</div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getOverviewInfo } from '@/api/modules/overview'
import { isObject, isArray, getDayTimeRange } from '@/utils'

interface ChartDataPoint { time: string; sends: number; delivered: number; bounced: number; open_rate: number; click_rate: number; bounce_rate: number }
interface MailProvider { provider: string; sends: number; delivered: number; bounced: number; open_rate: number; bounce_rate: number }

const period = ref('7d')
const periodOptions = [
	{ label: 'Today', value: '1d' },
	{ label: 'Last 7 days', value: '7d' },
	{ label: 'Last 30 days', value: '30d' },
]

const stats = ref({ delivery_rate: 0, open_rate: 0, click_rate: 0, bounce_rate: 0 })
const providers = ref<MailProvider[]>([])

const sendChartRef = ref<HTMLElement>()
const bounceChartRef = ref<HTMLElement>()
let sendChart: echarts.ECharts | null = null
let bounceChart: echarts.ECharts | null = null

const providerColumns = [
	{ key: 'provider', title: 'Provider', minWidth: 120 },
	{ key: 'sends', title: 'Sent', width: 100 },
	{ key: 'delivered', title: 'Delivered', width: 100 },
	{ key: 'bounced', title: 'Bounced', width: 100 },
	{ key: 'open_rate', title: 'Open %', width: 100, render: (row: MailProvider) => `${(row.open_rate || 0).toFixed(1)}%` },
	{ key: 'bounce_rate', title: 'Bounce %', width: 100, render: (row: MailProvider) => `${(row.bounce_rate || 0).toFixed(1)}%` },
]

function getDateRange() {
	const now = Date.now()
	const day = 86400000
	switch (period.value) {
		case '1d': return getDayTimeRange()
		case '7d': return [now - 7 * day, now]
		case '30d': return [now - 30 * day, now]
		default: return [now - 7 * day, now]
	}
}

function initCharts() {
	if (sendChartRef.value) sendChart = echarts.init(sendChartRef.value, 'dark')
	if (bounceChartRef.value) bounceChart = echarts.init(bounceChartRef.value, 'dark')
}

async function refresh() {
	const [start, end] = getDateRange()
	try {
		const res = await getOverviewInfo({
			domain: '',
			start_time: Math.floor(start / 1000),
			end_time: Math.floor(end / 1000),
		})
		if (isObject(res)) {
			stats.value = {
				delivery_rate: res.dashboard?.delivery_rate || 0,
				open_rate: res.dashboard?.open_rate || 0,
				click_rate: res.dashboard?.click_rate || 0,
				bounce_rate: res.dashboard?.bounce_rate || 0,
			}
			providers.value = isArray(res.mail_providers) ? res.mail_providers : []
			updateCharts(res)
		}
	} catch {}
}

function updateCharts(res: any) {
	if (!sendChart || !bounceChart) return

	const sendData = res.send_mail_chart?.data || []
	const bounceData = res.bounce_rate_chart?.data || []

	const sendXAxis = sendData.map((d: any) => {
		const dt = new Date(d.time * 1000)
		return `${dt.getMonth()+1}/${dt.getDate()}`
	})

	sendChart.setOption({
		backgroundColor: 'transparent',
		grid: { top: 20, right: 20, bottom: 30, left: 50 },
		xAxis: { type: 'category', data: sendXAxis, axisLine: { lineStyle: { color: '#2e3142' } }, axisLabel: { color: '#6b7280' } },
		yAxis: { type: 'value', splitLine: { lineStyle: { color: '#1e2030' } }, axisLabel: { color: '#6b7280' } },
		series: [
			{ name: 'Sent', type: 'line', smooth: true, data: sendData.map((d: any) => d.sends || 0), lineStyle: { color: '#6c5ce7' }, itemStyle: { color: '#6c5ce7' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(108,92,231,0.3)' }, { offset: 1, color: 'rgba(108,92,231,0)' }] } } },
			{ name: 'Delivered', type: 'line', smooth: true, data: sendData.map((d: any) => d.delivered || 0), lineStyle: { color: '#00d68f' }, itemStyle: { color: '#00d68f' } },
		],
		tooltip: { trigger: 'axis', backgroundColor: '#1a1d27', borderColor: '#2e3142', textStyle: { color: '#e2e8f0' } },
	})

	const bounceXAxis = bounceData.map((d: any) => {
		const dt = new Date(d.time * 1000)
		return `${dt.getMonth()+1}/${dt.getDate()}`
	})

	bounceChart.setOption({
		backgroundColor: 'transparent',
		grid: { top: 20, right: 20, bottom: 30, left: 50 },
		xAxis: { type: 'category', data: bounceXAxis, axisLine: { lineStyle: { color: '#2e3142' } }, axisLabel: { color: '#6b7280' } },
		yAxis: { type: 'value', splitLine: { lineStyle: { color: '#1e2030' } }, axisLabel: { color: '#6b7280', formatter: '{value}%' } },
		series: [
			{ name: 'Bounce Rate', type: 'bar', data: bounceData.map((d: any) => ((d.bounce_rate || 0) * 100).toFixed(1)), itemStyle: { color: '#ff3d71', borderRadius: [4, 4, 0, 0] } },
		],
		tooltip: { trigger: 'axis', backgroundColor: '#1a1d27', borderColor: '#2e3142', textStyle: { color: '#e2e8f0' }, formatter: (params: any) => `${params[0].name}<br/>${params[0].value}%` },
	})
}

function handleResize() { sendChart?.resize(); bounceChart?.resize() }

onMounted(async () => {
	await nextTick()
	initCharts()
	refresh()
	window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
	window.removeEventListener('resize', handleResize)
	sendChart?.dispose()
	bounceChart?.dispose()
})
</script>

<style lang="scss" scoped>
.analytics-page { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { margin: 0; font-size: 22px; font-weight: 700; color: #e2e8f0; }

.kpi-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
.kpi-card {
	display: flex; align-items: center; gap: 12px; padding: 16px;
	background: #1a1d27; border: 1px solid #2e3142; border-radius: 10px;
}
.kpi-icon { font-size: 22px; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; background: rgba(108,92,231,0.08); border-radius: 10px; }
.kpi-val { font-size: 20px; font-weight: 700; color: #f1f3f7; }
.kpi-lbl { font-size: 12px; color: #6b7280; margin-top: 2px; }

.charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 24px; }
.chart-card {
	background: #1a1d27; border: 1px solid #2e3142; border-radius: 12px; padding: 20px;
}
.chart-container { height: 280px; }
.card-title { font-size: 15px; font-weight: 600; color: #e2e8f0; margin-bottom: 12px; }

.providers-section {
	background: #1a1d27; border: 1px solid #2e3142; border-radius: 12px; padding: 20px;
}

:deep(.n-data-table) { --n-td-color: transparent; --n-th-color: transparent; --n-border-color: #2e3142; --n-text-color: #e2e8f0; }
</style>
