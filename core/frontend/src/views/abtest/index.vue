<template>
	<div class="abtest-page">
		<div class="bt-title">A/B Tests</div>

		<!-- Top Actions -->
		<div class="actions-row">
			<n-space>
				<n-select
					v-model:value="selectedSequence"
					:options="sequenceOptions"
					placeholder="Select a sequence"
					style="width: 240px"
					size="small"
					@update:value="fetchTestsForSequence"
				/>
				<n-button type="primary" size="small" @click="showCreateModal = true" :disabled="!selectedSequence">
					New A/B Test
				</n-button>
			</n-space>
		</div>

		<!-- Tests List -->
		<n-spin :show="loadingTests">
			<div v-if="tests.length === 0 && !loadingTests" class="empty-state">
				<n-empty description="Select a sequence to view A/B tests" />
			</div>
			<div v-else class="test-cards">
				<n-card v-for="test in tests" :key="test.id" class="test-card">
					<template #header>
						<div class="test-header">
							<span class="test-name">{{ test.name || 'Test #' + test.id }}</span>
							<n-tag :type="statusTagType(test.status)" size="small" bordered>
								{{ statusLabel(test.status) }}
							</n-tag>
						</div>
					</template>
					<template #header-extra>
						<n-space size="small">
							<n-button size="tiny" @click="viewResults(test)">Results</n-button>
							<n-button
								v-if="test.status === 1"
								size="tiny"
								type="primary"
								@click="handlePickWinner(test)"
							>
								Pick Winner
							</n-button>
							<n-button size="tiny" type="error" quaternary @click="handleDelete(test)">
								Delete
							</n-button>
						</n-space>
					</template>
					<div class="test-meta">
						<span>Step {{ test.step_id }}</span>
						<span>Split: {{ (test.split_ratio * 100).toFixed(0) }}/{{ ((1 - test.split_ratio) * 100).toFixed(0) }}</span>
						<span>Criteria: {{ test.winner_criteria }}</span>
						<span v-if="test.winner_variant >= 0">
							Winner: <n-tag size="tiny" type="success">{{ test.winner_variant === 0 ? 'A' : 'B' }}</n-tag>
						</span>
					</div>
				</n-card>
			</div>
		</n-spin>

		<!-- Results Modal -->
		<results-modal
			v-model:show="showResultsModal"
			:test-id="selectedTestId"
		/>

		<!-- Create Modal -->
		<create-modal
			v-model:show="showCreateModal"
			:sequence-id="selectedSequence"
			@created="onTestCreated"
		/>
	</div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Message } from '@/utils'
import { getSequenceList } from '@/api/modules/sequences/sequence'
import { listAbTests, deleteAbTest, pickWinner } from '@/api/modules/abtest'
import ResultsModal from './components/ResultsModal.vue'
import CreateModal from './components/CreateModal.vue'

interface AbTest {
	id: number
	sequence_id: number
	step_id: number
	name: string
	status: number
	winner_variant: number
	winner_criteria: string
	split_ratio: number
	created_at: number
	completed_at: number
}

// using Message from utils

const selectedSequence = ref<number | null>(null)
const sequenceOptions = ref<Array<{ label: string; value: number }>>([])
const tests = ref<AbTest[]>([])
const loadingTests = ref(false)
const showResultsModal = ref(false)
const showCreateModal = ref(false)
const selectedTestId = ref(0)

function statusLabel(s: number) {
	return ['Draft', 'Running', 'Completed'][s] || 'Unknown'
}

function statusTagType(s: number) {
	return (['default', 'success', 'info'] as const)[s] || 'default'
}

async function fetchSequences() {
	try {
		const res = await getSequenceList({ page: 1, page_size: 100 })
		if (res?.list) {
			sequenceOptions.value = res.list.map((s: any) => ({
				label: s.name || 'Sequence #' + s.id,
				value: s.id,
			}))
		}
	} catch (e) {
		console.error('Failed to fetch sequences', e)
	}
}

async function fetchTestsForSequence() {
	if (!selectedSequence.value) return
	loadingTests.value = true
	try {
		const res = await listAbTests({ sequence_id: selectedSequence.value })
		tests.value = Array.isArray(res) ? res : []
	} catch (e) {
		console.error('Failed to fetch tests', e)
	} finally {
		loadingTests.value = false
	}
}

function viewResults(test: AbTest) {
	selectedTestId.value = test.id
	showResultsModal.value = true
}

async function handlePickWinner(test: AbTest) {
	try {
		await pickWinner({ id: test.id, winner_variant: 0 }) // Default pick A
		Message.success('Winner picked')
		fetchTestsForSequence()
	} catch (e) {
		Message.error('Failed to pick winner')
	}
}

async function handleDelete(test: AbTest) {
	try {
		await deleteAbTest({ id: test.id })
		Message.success('Test deleted')
		fetchTestsForSequence()
	} catch (e) {
		Message.error('Failed to delete test')
	}
}

function onTestCreated() {
	showCreateModal.value = false
	fetchTestsForSequence()
}

onMounted(() => {
	fetchSequences()
})
</script>

<style lang="scss" scoped>
.abtest-page {
	padding: 20px;
}

.actions-row {
	margin-bottom: 16px;
}

.empty-state {
	padding: 60px 0;
}

.test-cards {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.test-card {
	.test-header {
		display: flex;
		align-items: center;
		gap: 10px;

		.test-name {
			font-weight: 600;
		}
	}

	.test-meta {
		display: flex;
		gap: 16px;
		font-size: 13px;
		color: var(--n-text-color-3);
	}
}
</style>
