<template>
	<n-modal v-model:show="show" preset="card" title="A/B Test Results" style="width: 680px">
		<n-spin :show="loading">
			<template v-if="result">
				<div class="result-header">
					<span class="test-name">{{ result.test_name }}</span>
					<n-tag :type="result.status === 2 ? 'success' : 'info'" size="small">
						{{ result.status === 2 ? 'Completed' : 'Running' }}
					</n-tag>
					<n-tag v-if="result.is_significant" type="success" size="small">
						Statistically Significant
					</n-tag>
				</div>

				<div v-if="result.winner >= 0" class="winner-banner">
					<i class="i-mdi-trophy" /> Winner: Variant {{ result.winner === 0 ? 'A' : 'B' }}
				</div>

				<div class="variants-compare">
					<variant-card
						label="A"
						:stats="result.variant_a"
						:is-winner="result.winner === 0"
					/>
					<div class="vs-divider">VS</div>
					<variant-card
						label="B"
						:stats="result.variant_b"
						:is-winner="result.winner === 1"
					/>
				</div>
			</template>
		</n-spin>
	</n-modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { getAbTestResults } from '@/api/modules/abtest'
import VariantCard from './VariantCard.vue'

interface VariantStats {
	variant_id: number
	label: string
	sent_count: number
	open_rate: number
	click_rate: number
	reply_rate: number
	bounce_rate: number
}

interface TestResult {
	test_id: number
	test_name: string
	status: number
	winner: number
	variant_a: VariantStats | null
	variant_b: VariantStats | null
	is_significant: boolean
}

const props = defineProps<{
	show: boolean
	testId: number
}>()

const emit = defineEmits<{
	'update:show': [value: boolean]
}>()

const show = computed({
	get: () => props.show,
	set: (val) => emit('update:show', val),
})

const loading = ref(false)
const result = ref<TestResult | null>(null)

async function fetchResults() {
	if (!props.testId) return
	loading.value = true
	try {
		const res = await getAbTestResults({ id: props.testId })
		result.value = res || null
	} catch (e) {
		console.error('Failed to fetch results', e)
	} finally {
		loading.value = false
	}
}

watch(() => props.show, (val) => {
	if (val && props.testId) fetchResults()
})
</script>

<style lang="scss" scoped>
.result-header {
	display: flex;
	align-items: center;
	gap: 10px;
	margin-bottom: 16px;

	.test-name {
		font-size: 16px;
		font-weight: 600;
	}
}

.winner-banner {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 10px 16px;
	background: rgba(24, 160, 88, 0.08);
	border-radius: 8px;
	margin-bottom: 16px;
	font-weight: 600;
	color: #18a058;
	font-size: 15px;
}

.variants-compare {
	display: flex;
	gap: 16px;
	align-items: stretch;

	.vs-divider {
		display: flex;
		align-items: center;
		font-size: 18px;
		font-weight: 700;
		color: var(--n-text-color-3);
	}
}
</style>
