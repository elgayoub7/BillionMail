<template>
	<n-modal v-model:show="show" preset="card" title="Create A/B Test" style="width: 720px">
		<n-form ref="formRef" :model="form" label-placement="left" label-width="120">
			<n-form-item label="Test Name">
				<n-input v-model:value="form.name" placeholder="e.g. Subject line test" />
			</n-form-item>
			<n-form-item label="Step ID">
				<n-input-number v-model:value="form.step_id" :min="0" placeholder="Step number" style="width: 100%" />
			</n-form-item>
			<n-form-item label="Split Ratio">
				<n-slider v-model:value="form.split_ratio" :min="0.1" :max="0.9" :step="0.05" :marks="{ 0.5: '50/50' }" />
				<span class="ml-2 text-sm">{{ (form.split_ratio * 100).toFixed(0) }}% / {{ ((1 - form.split_ratio) * 100).toFixed(0) }}%</span>
			</n-form-item>
			<n-form-item label="Winner Criteria">
				<n-select v-model:value="form.winner_criteria" :options="criteriaOptions" />
			</n-form-item>

			<n-divider>Variant A</n-divider>
			<n-form-item label="Subject A">
				<n-input v-model:value="form.variant_a.subject" placeholder="Subject line for variant A" />
			</n-form-item>
			<n-form-item label="Body HTML A">
				<n-input v-model:value="form.variant_a.body_html" type="textarea" :rows="3" placeholder="HTML body for variant A" />
			</n-form-item>

			<n-divider>Variant B</n-divider>
			<n-form-item label="Subject B">
				<n-input v-model:value="form.variant_b.subject" placeholder="Subject line for variant B" />
			</n-form-item>
			<n-form-item label="Body HTML B">
				<n-input v-model:value="form.variant_b.body_html" type="textarea" :rows="3" placeholder="HTML body for variant B" />
			</n-form-item>
		</n-form>

		<template #action>
			<n-space justify="end">
				<n-button @click="show = false">Cancel</n-button>
				<n-button type="primary" @click="handleCreate" :loading="creating">Create Test</n-button>
			</n-space>
		</template>
	</n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { Message } from '@/utils'
import { createAbTest } from '@/api/modules/abtest'

const props = defineProps<{
	show: boolean
	sequenceId: number | null
}>()

const emit = defineEmits<{
	'update:show': [value: boolean]
	created: []
}>()

// using Message from utils
const creating = ref(false)

const show = computed({
	get: () => props.show,
	set: (val) => emit('update:show', val),
})

const form = ref({
	name: '',
	step_id: 0,
	split_ratio: 0.5,
	winner_criteria: 'open_rate',
	variant_a: { subject: '', body_html: '', body_text: '' },
	variant_b: { subject: '', body_html: '', body_text: '' },
})

const criteriaOptions = [
	{ label: 'Open Rate', value: 'open_rate' },
	{ label: 'Click Rate', value: 'click_rate' },
	{ label: 'Reply Rate', value: 'reply_rate' },
]

async function handleCreate() {
	if (!props.sequenceId) {
		Message.error('Select a sequence first')
		return
	}
	if (!form.value.variant_a.subject || !form.value.variant_b.subject) {
		Message.error('Both variants need a subject')
		return
	}
	creating.value = true
	try {
		await createAbTest({
			sequence_id: props.sequenceId,
			step_id: form.value.step_id,
			name: form.value.name,
			split_ratio: form.value.split_ratio,
			winner_criteria: form.value.winner_criteria,
			variant_a: form.value.variant_a,
			variant_b: form.value.variant_b,
		})
		Message.success('A/B test created')
		emit('created')
		show.value = false
	} catch (e) {
		Message.error('Failed to create test')
	} finally {
		creating.value = false
	}
}
</script>
