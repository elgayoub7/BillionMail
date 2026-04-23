<template>
	<n-card class="variant-card" :class="{ winner: isWinner }" size="small">
		<template #header>
			<div class="variant-label">
				<span v-if="isWinner" class="crown"><i class="i-mdi-crown" /></span>
				Variant {{ label }}
			</div>
		</template>
		<div class="stat-grid">
			<div class="stat-item">
				<div class="stat-val">{{ stats?.sent_count || 0 }}</div>
				<div class="stat-lbl">Sent</div>
			</div>
			<div class="stat-item">
				<div class="stat-val" :style="{ color: rateColor(stats?.open_rate) }">
					{{ (stats?.open_rate || 0).toFixed(1) }}%
				</div>
				<div class="stat-lbl">Open Rate</div>
			</div>
			<div class="stat-item">
				<div class="stat-val" :style="{ color: rateColor(stats?.click_rate) }">
					{{ (stats?.click_rate || 0).toFixed(1) }}%
				</div>
				<div class="stat-lbl">Click Rate</div>
			</div>
			<div class="stat-item">
				<div class="stat-val" :style="{ color: rateColor(stats?.reply_rate) }">
					{{ (stats?.reply_rate || 0).toFixed(1) }}%
				</div>
				<div class="stat-lbl">Reply Rate</div>
			</div>
			<div class="stat-item">
				<div class="stat-val" :style="{ color: bounceColor(stats?.bounce_rate) }">
					{{ (stats?.bounce_rate || 0).toFixed(1) }}%
				</div>
				<div class="stat-lbl">Bounce Rate</div>
			</div>
		</div>
	</n-card>
</template>

<script lang="ts" setup>
defineProps<{
	label: string
	stats: {
		sent_count: number
		open_rate: number
		click_rate: number
		reply_rate: number
		bounce_rate: number
	} | null
	isWinner: boolean
}>()

function rateColor(rate: number | undefined) {
	if (!rate) return 'inherit'
	if (rate > 20) return '#18a058'
	if (rate > 5) return '#e6a23c'
	return '#d03050'
}

function bounceColor(rate: number | undefined) {
	if (!rate) return 'inherit'
	if (rate < 5) return '#18a058'
	if (rate < 10) return '#e6a23c'
	return '#d03050'
}
</script>

<style lang="scss" scoped>
.variant-card {
	flex: 1;
	transition: all 0.2s;

	&.winner {
		border: 2px solid #18a058;
		box-shadow: 0 0 12px rgba(24, 160, 88, 0.15);
	}
}

.variant-label {
	display: flex;
	align-items: center;
	gap: 6px;
	font-weight: 600;

	.crown {
		color: #f0a030;
		font-size: 18px;
	}
}

.stat-grid {
	display: grid;
	grid-template-columns: repeat(5, 1fr);
	gap: 8px;
	text-align: center;
}

.stat-item {
	.stat-val {
		font-size: 18px;
		font-weight: 700;
	}

	.stat-lbl {
		font-size: 11px;
		color: var(--n-text-color-3);
		margin-top: 2px;
	}
}
</style>
