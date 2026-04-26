<template>
	<modal :title="title" :width="1200">
		<bt-form class="pt-8px" :model="form">
			<n-grid :cols="24" :x-gap="24">
				<n-form-item-gi :span="18" label="Name" path="temp_name">
					<n-input
						v-model:value="form.temp_name"
						:placeholder="$t('market.template.add.namePlaceholder')">
					</n-input>
				</n-form-item-gi>
				<n-form-item-gi :span="6" label="Type" path="add_type">
					<n-select v-model:value="form.add_type" :disabled="isEdit" :options="typeOptions">
					</n-select>
				</n-form-item-gi>
			</n-grid>

			<!-- A/B Test Toggle -->
			<n-form-item label="A/B Test">
				<div class="flex items-center gap-12px w-full">
					<n-switch v-model:value="form.ab_enabled" size="small" />
					<span class="text-sm text-gray-400">
						{{ form.ab_enabled ? 'A/B test enabled — recipients will receive one of two variants' : 'Disabled' }}
					</span>
				</div>
			</n-form-item>

			<!-- A/B Settings -->
			<template v-if="form.ab_enabled">
				<n-grid :cols="24" :x-gap="24">
					<n-form-item-gi :span="12" label="Split Ratio">
						<div class="flex items-center gap-12px w-full">
							<n-slider v-model:value="form.ab_split_ratio" :min="10" :max="90" :step="5" class="flex-1" />
							<span class="text-sm font-mono w-80px text-right">{{ form.ab_split_ratio }}% / {{ 100 - form.ab_split_ratio }}%</span>
						</div>
					</n-form-item-gi>
					<n-form-item-gi :span="12" label="Winner Criteria">
						<n-select v-model:value="form.ab_winner_criteria" :options="criteriaOptions" />
					</n-form-item-gi>
				</n-grid>
			</template>

			<n-divider v-if="form.ab_enabled" style="margin: 4px 0">Variant A</n-divider>

			<n-form-item :label="typeLabel" :show-feedback="false">
				<div class="flex flex-col flex-1 h-580px">
					<div class="flex justify-end gap-12px mb-8px">
						<div>
							<bt-file-upload
								mode="button"
								button-type="primary"
								:is-upload="false"
								:accept="['html']"
								@change="handleFileUpload">
								{{ $t('market.template.uploadHtml') }}
							</bt-file-upload>
						</div>
						<n-button type="primary" ghost @click="handlePreview">
							{{ $t('common.actions.preview') }}
						</n-button>
					</div>
					<div class="flex-1">
						<template v-if="form.add_type === 1">
							<email-editor v-model:config="form.drag_data" v-model:html="form.html_content">
							</email-editor>
						</template>
						<template v-if="form.add_type === 0">
							<bt-editor v-model:value="form.html_content" language="html"> </bt-editor>
						</template>
					</div>
				</div>
			</n-form-item>

			<!-- Variant B (A/B Test) -->
			<template v-if="form.ab_enabled">
				<n-divider style="margin: 4px 0">Variant B</n-divider>
				<n-form-item label="Subject B">
					<n-input v-model:value="form.variant_b_subject" placeholder="Alternative subject line for variant B" />
				</n-form-item>
				<n-form-item label="Body B" :show-feedback="false">
					<div class="flex flex-col flex-1 h-400px">
						<bt-editor v-model:value="form.variant_b_html" language="html" />
					</div>
				</n-form-item>
			</template>
		</bt-form>
	</modal>
</template>

<script lang="ts" setup>
import { UploadFileInfo } from 'naive-ui'
import { useModal } from '@/hooks/modal/useModal'
import { addTemplate, updateTemplate } from '@/api/modules/market/template'
import type { Template } from '../interface'

const EmailEditor = defineAsyncComponent(() => import('@/features/EmailEditor/index.vue'))

const isEdit = ref(false)

const title = computed(() => {
	return isEdit.value ? 'Edit template' : 'New template'
})

const id = ref(0)

const form = reactive({
	temp_name: '',
	add_type: 0,
	html_content: '',
	drag_data: '',
	ab_enabled: false,
	ab_split_ratio: 50,
	ab_winner_criteria: 'open_rate',
	variant_b_subject: '',
	variant_b_html: '',
})

const typeOptions = [
	{ label: 'HTML', value: 0 },
	{ label: 'Drag', value: 1 },
]

const criteriaOptions = [
	{ label: 'Open Rate', value: 'open_rate' },
	{ label: 'Click Rate', value: 'click_rate' },
	{ label: 'Reply Rate', value: 'reply_rate' },
]

const typeLabel = computed(() => {
	return typeOptions.find(item => item.value === form.add_type)?.label || ''
})

const handleFileUpload = (file: UploadFileInfo) => {
	const reader = new FileReader()
	reader.onload = e => {
		form.html_content = e.target?.result as string
	}
	if (file.file) {
		reader.readAsText(file.file)
	}
}

const handlePreview = () => {
	const state = modalApi.getState<{ preview: (content: string) => void }>()
	if (state.preview) {
		state.preview(form.html_content)
	}
}

const resetForm = () => {
	id.value = 0
	form.temp_name = ''
	form.add_type = 0
	form.html_content = ''
	form.drag_data = ''
	form.ab_enabled = false
	form.ab_split_ratio = 50
	form.ab_winner_criteria = 'open_rate'
	form.variant_b_subject = ''
	form.variant_b_html = ''
}

const [Modal, modalApi] = useModal({
	onChangeState: isOpen => {
		if (isOpen) {
			const state = modalApi.getState<{ isEdit: boolean; type: number; row: Template }>()
			const { row } = state
			isEdit.value = state.isEdit
			resetForm()
			if (state.type) {
				form.add_type = state.type
			}
			if (row) {
				id.value = row.id
				form.temp_name = row.temp_name
				form.add_type = row.add_type
				form.html_content = row.html_content
				form.drag_data = row.drag_data
				form.ab_enabled = !!row.ab_enabled
				form.ab_split_ratio = row.ab_split_ratio || 50
				form.ab_winner_criteria = row.ab_winner_criteria || 'open_rate'
				form.variant_b_subject = row.variant_b_subject || ''
				form.variant_b_html = row.variant_b_html || ''
			}
		} else {
			resetForm()
		}
	},
	onConfirm: async () => {
		const payload = toRaw(form)
		if (isEdit.value) {
			await updateTemplate({ id: id.value, ...payload })
		} else {
			await addTemplate(payload)
		}
		const state = modalApi.getState<{ refresh: () => void }>()
		state.refresh()
	},
})
</script>

<style lang="scss" scoped></style>
