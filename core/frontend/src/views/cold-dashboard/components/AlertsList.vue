<template>
	<div class="alerts-list">
		<n-spin :show="loading">
			<div v-if="alerts.length === 0 && !loading" class="empty-state">
				<n-empty description="No alerts — everything looks good!" />
			</div>
			<div v-else class="alert-items">
				<div
					v-for="(alert, idx) in alerts"
					:key="idx"
					class="alert-item"
					:class="alertTypeClass(alert.type)"
				>
					<div class="alert-icon">
						<i :class="alertIcon(alert.type)" />
					</div>
					<div class="alert-body">
						<div class="alert-msg">{{ alert.message }}</div>
						<div v-if="alert.time" class="alert-time">{{ formatTime(alert.time) }}</div>
					</div>
				</div>
			</div>
		</n-spin>
	</div>
</template>

<script lang="ts" setup>
import { format } from 'date-fns'

defineProps<{
	alerts: Array<{
		type: string
		message: string
		time: number
	}>
	loading: boolean
}>()

function alertTypeClass(type: string) {
	const map: Record<string, string> = {
		hard_bounce: 'alert-error',
		domain_health: 'alert-warning',
		abtest_winner: 'alert-success',
	}
	return map[type] || 'alert-info'
}

function alertIcon(type: string) {
	const map: Record<string, string> = {
		hard_bounce: 'i-mdi-alert-circle-outline',
		domain_health: 'i-mdi-alert-outline',
		abtest_winner: 'i-mdi-trophy-outline',
	}
	return map[type] || 'i-mdi-information-outline'
}

function formatTime(ts: number) {
	if (!ts) return ''
	try {
		return format(new Date(ts * 1000), 'MMM d, HH:mm')
	} catch {
		return ''
	}
}
</script>

<style lang="scss" scoped>
.alerts-list {
	max-height: 400px;
	overflow-y: auto;
}

.empty-state {
	padding: 40px 0;
}

.alert-items {
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.alert-item {
	display: flex;
	align-items: flex-start;
	gap: 10px;
	padding: 10px 12px;
	border-radius: 8px;
	border-left: 3px solid transparent;

	&.alert-error {
		background: rgba(208, 48, 80, 0.06);
		border-left-color: #d03050;
	}
	&.alert-warning {
		background: rgba(240, 160, 48, 0.06);
		border-left-color: #f0a030;
	}
	&.alert-success {
		background: rgba(24, 160, 88, 0.06);
		border-left-color: #18a058;
	}
	&.alert-info {
		background: rgba(48, 120, 208, 0.06);
		border-left-color: #3078d0;
	}
}

.alert-icon {
	font-size: 18px;
	margin-top: 1px;
	flex-shrink: 0;
}

.alert-body {
	flex: 1;
}

.alert-msg {
	font-size: 13px;
	line-height: 1.4;
}

.alert-time {
	font-size: 11px;
	color: var(--n-text-color-3);
	margin-top: 2px;
}
</style>
