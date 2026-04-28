<template>
	<div class="rate-charts">
		<div class="card">
			<div class="card-header">
				<span class="card-title">Open Rate</span>
			</div>
			<line-chart
				chart-color="#6c5ce7"
				:date-type="open.column_type"
				chart-name="Open Rate"
				:chart-data="openRateData" />
		</div>
		<div class="card">
			<div class="card-header">
				<span class="card-title">Click Rate</span>
			</div>
			<line-chart
				chart-color="#22c55e"
				:date-type="click.column_type"
				chart-name="Click Rate"
				:chart-data="clickRateData" />
		</div>
		<div class="card">
			<div class="card-header">
				<span class="card-title">Bounce Rate</span>
			</div>
			<line-chart
				chart-color="#ef4444"
				:date-type="bounce.column_type"
				chart-name="Bounce Rate"
				:chart-data="bounceRateData" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { PropType } from 'vue'
import { formatTime } from '@/utils'
import { MailOverview } from '../types'

import LineChart from './LineChart.vue'

const { open, click, bounce } = defineProps({
	bounce: {
		type: Object as PropType<MailOverview['bounce_rate_chart']>,
		required: true,
	},
	click: {
		type: Object as PropType<MailOverview['click_rate_chart']>,
		required: true,
	},
	open: {
		type: Object as PropType<MailOverview['open_rate_chart']>,
		required: true,
	},
})

const getChartTime = (type: string, x: number) => {
	let date = new Date()
	let format = ''
	if (type === 'hourly') {
		date.setMinutes(0)
		date.setSeconds(0)
		date.setHours(x)
		format = 'yyyy-MM-dd HH:mm'
	} else if (type === 'daily') {
		date = new Date(x * 1000)
		format = 'yyyy-MM-dd'
	}
	return formatTime(date, format)
}

const openRateData = computed(() =>
	open.data.map(item => [getChartTime(open.column_type, item.x), item.open_rate] as [string, number])
)

const clickRateData = computed(() =>
	click.data.map(item => [getChartTime(click.column_type, item.x), item.click_rate] as [string, number])
)

const bounceRateData = computed(() =>
	bounce.data.map(item => [getChartTime(bounce.column_type, item.x), item.bounce_rate] as [string, number])
)
</script>

<style lang="scss" scoped>
.rate-charts {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 16px;
}

.card {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
}

.card-header {
	margin-bottom: 12px;
}

.card-title {
	font-size: 15px;
	font-weight: 600;
	color: #e2e8f0;
}

@media (max-width: 1200px) {
	.rate-charts {
		grid-template-columns: 1fr;
	}
}
</style>
