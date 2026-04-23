<template>
	<n-card
		class="score-stat-card"
		:class="{ active }"
		:style="{ borderLeftColor: color }"
		@click="$emit('click')"
		hoverable
	>
		<div class="stat-icon"><i :class="icon" :style="{ color }" /></div>
		<div class="stat-body">
			<div class="stat-count">{{ count }}</div>
			<div class="stat-label">{{ label }}</div>
		</div>
	</n-card>
</template>

<script lang="ts" setup>
defineProps<{
	label: string
	count: number
	color: string
	icon: string
	active: boolean
}>()

defineEmits<{
	click: []
}>()
</script>

<style lang="scss" scoped>
.score-stat-card {
	cursor: pointer;
	border-left: 3px solid transparent;
	transition: all 0.2s;

	&.active {
		box-shadow: 0 0 0 2px var(--n-primary-color);
	}

	display: flex;
	align-items: center;
	gap: 12px;

	.stat-icon {
		font-size: 28px;
	}

	.stat-body {
		flex: 1;
	}

	.stat-count {
		font-size: 22px;
		font-weight: 700;
		line-height: 1.2;
	}

	.stat-label {
		font-size: 12px;
		color: var(--n-text-color-3);
	}
}
</style>
