<template>
	<div class="chart-wrapper">
		<bt-charts :options="chartOptions" />
	</div>
</template>

<script setup lang="ts">
import { formatTime } from '@/utils'

const { chartName, chartColor, dateType, chartData } = defineProps({
	chartName: {
		type: String,
		required: true,
	},
	chartColor: {
		type: String,
		required: true,
	},
	dateType: {
		type: String,
		required: true,
	},
	chartData: {
		type: Array as PropType<[string, number][]>,
		required: true,
	},
})

const chartOptions = computed(() => ({
	tooltip: {
		trigger: 'axis',
		background: '#1a1d27',
		borderColor: '#2e3142',
		textStyle: { color: '#e2e8f0', fontSize: 12 },
		axisPointer: { type: 'line', lineStyle: { color: '#2e3142' } },
	},
	grid: {
		top: '8%',
		left: '2%',
		right: '4%',
		bottom: '2%',
		containLabel: true,
	},
	xAxis: {
		type: 'category',
		axisLine: { lineStyle: { color: '#2e3142' } },
		axisLabel: {
			color: '#6b7280',
			fontSize: 11,
			formatter: (val: string) => {
				if (dateType === 'hourly') return formatTime(val, 'HH:mm')
				return formatTime(val, 'MM-dd')
			},
		},
		axisTick: { show: false },
	},
	yAxis: {
		name: '%',
		type: 'value',
		max: 100,
		splitLine: { lineStyle: { type: 'dashed', width: 1, color: '#1e2030' } },
		axisLabel: { color: '#6b7280', fontSize: 11 },
		axisLine: { show: false },
	},
	series: [
		{
			name: chartName,
			type: 'line',
			data: chartData,
			itemStyle: { color: chartColor },
			lineStyle: { width: 2 },
			areaStyle: {
				color: {
					type: 'linear',
					x: 0, y: 0, x2: 0, y2: 1,
					colorStops: [
						{ offset: 0, color: chartColor + '30' },
						{ offset: 1, color: chartColor + '05' },
					],
				},
			},
			smooth: true,
			showSymbol: false,
			sampling: 'average',
		},
	],
}))
</script>

<style lang="scss" scoped>
.chart-wrapper {
	height: 200px;
}
</style>
