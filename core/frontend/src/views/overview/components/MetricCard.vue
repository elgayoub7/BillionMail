<template>
	<div class="metric-card" :class="{ clickable: !!onClick }" @click="onClick?.()">
		<div class="metric-header">
			<span class="metric-title">{{ title }}</span>
			<i v-if="icon" :class="icon" class="metric-icon"></i>
		</div>
		<div class="metric-value" :style="{ color: valueColor }">
			{{ formattedValue }}<span v-if="unit" class="metric-unit">{{ unit }}</span>
		</div>
		<div v-if="trend !== undefined" class="metric-trend" :class="trendClass">
			<i :class="trendIcon"></i>
			<span>{{ Math.abs(trend) }}%</span>
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps({
	title: { type: String, default: '' },
	value: { type: Number, default: 0 },
	unit: { type: String, default: '' },
	icon: { type: String, default: '' },
	valueColor: { type: String, default: '' },
	trend: { type: Number, default: undefined },
	onClick: { type: Function, default: undefined },
})

const formattedValue = computed(() => {
	if (props.value >= 1000000) return (props.value / 1000000).toFixed(1) + 'M'
	if (props.value >= 1000) return (props.value / 1000).toFixed(1) + 'K'
	return props.value.toLocaleString()
})

const trendClass = computed(() => {
	if (props.trend === undefined) return ''
	return props.trend >= 0 ? 'trend-up' : 'trend-down'
})

const trendIcon = computed(() => {
	if (props.trend === undefined) return ''
	return props.trend >= 0 ? 'i-mdi:trending-up' : 'i-mdi:trending-down'
})
</script>

<style lang="scss" scoped>
.metric-card {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
	transition: all 0.2s ease;

	&.clickable {
		cursor: pointer;

		&:hover {
			border-color: #6c5ce7;
			transform: translateY(-1px);
		}
	}
}

.metric-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 12px;
}

.metric-title {
	font-size: 13px;
	font-weight: 500;
	color: #8892a8;
}

.metric-icon {
	font-size: 18px;
	color: #6c5ce7;
	opacity: 0.6;
}

.metric-value {
	font-size: 26px;
	font-weight: 700;
	color: #e2e8f0;
	line-height: 1.2;
}

.metric-unit {
	font-size: 14px;
	font-weight: 500;
	color: #8892a8;
	margin-left: 2px;
}

.metric-trend {
	display: flex;
	align-items: center;
	gap: 4px;
	margin-top: 8px;
	font-size: 12px;
	font-weight: 600;

	&.trend-up {
		color: #22c55e;
	}

	&.trend-down {
		color: #ef4444;
	}
}
</style>
