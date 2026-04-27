<template>
	<modal :title="$t('contacts.subscribers.import.title')" width="900" :footer="false">
		<div class="pt-12px">
			<!-- Step indicator -->
			<n-steps :current="currentStep" class="mb-24px" size="small">
				<n-step :title="$t('contacts.subscribers.import.steps.upload')" />
				<n-step :title="$t('contacts.subscribers.import.steps.preview')" />
				<n-step :title="$t('contacts.subscribers.import.steps.import')" />
			</n-steps>

			<!-- STEP 1: Upload / Paste -->
			<div v-if="currentStep === 1">
				<bt-form ref="formRef" :model="form" :rules="rules">
					<n-form-item :label="$t('contacts.subscribers.import.data')" path="contacts">
						<div class="flex-1">
							<n-radio-group v-model:value="form.import_type">
								<n-radio-button :value="1">
									{{ $t('contacts.subscribers.import.uploadFile') }}
								</n-radio-button>
								<n-radio-button :value="2">
									{{ $t('contacts.subscribers.import.pasteData') }}
								</n-radio-button>
							</n-radio-group>
							<div class="mt-8px">
								<template v-if="form.import_type === 1">
									<bt-file-upload :is-upload="false" :accept="['csv', 'txt']" @change="handleChangeFile">
									</bt-file-upload>
									<div class="mt-8px text-12px text-desc">
										{{ $t('contacts.subscribers.import.fileTypeHint') }}
										<n-button text type="primary" @click="handleDownloadTemplate">
											{{ $t('common.actions.download') }}
										</n-button>
									</div>
								</template>
								<template v-if="form.import_type === 2">
									<n-input
										v-model:value="form.contacts"
										type="textarea"
										:rows="8"
										:placeholder="$t('contacts.subscribers.import.pastePlaceholder.0')"
									/>
								</template>
							</div>
						</div>
					</n-form-item>
					<n-form-item :label="$t('contacts.subscribers.import.exampleCsvLabel')">
						<pre class="csv-example">{{ example }}</pre>
					</n-form-item>
				</bt-form>
				<div class="flex justify-end gap-12px mt-16px">
					<n-button @click="handleCancel">{{ $t('common.actions.cancel') }}</n-button>
					<n-button type="primary" :loading="previewLoading" :disabled="!hasData" @click="handlePreview">
						{{ $t('contacts.subscribers.import.actions.next') }}
					</n-button>
				</div>
			</div>

			<!-- STEP 2: Preview & Mapping -->
			<div v-if="currentStep === 2">
				<!-- Stats row -->
				<n-grid :cols="3" :x-gap="16" class="mb-16px">
					<n-gi>
						<n-statistic :label="$t('contacts.subscribers.import.stats.totalRows')" :value="previewData.total_rows" />
					</n-gi>
					<n-gi>
						<n-statistic :label="$t('contacts.subscribers.import.stats.duplicates')" :value="previewData.duplicates" />
					</n-gi>
					<n-gi>
						<n-statistic :label="$t('contacts.subscribers.import.stats.invalidRows')" :value="previewData.invalid_rows" />
					</n-gi>
				</n-grid>

				<!-- Column mapping -->
				<n-card size="small" class="mb-16px" :title="$t('contacts.subscribers.import.mapping.title')">
					<n-grid :cols="2" :x-gap="16" :y-gap="8">
						<n-gi v-for="(header, idx) in previewData.headers" :key="idx">
							<div class="flex items-center gap-8px">
								<n-tag size="small" type="info" class="flex-shrink-0" style="min-width:100px">{{ header }}</n-tag>
								<n-select
									v-model:value="columnMapping[idx]"
									:options="fieldOptions"
									size="small"
									class="flex-1"
								/>
							</div>
						</n-gi>
					</n-grid>
					<div v-if="previewData.unmapped?.length" class="mt-8px text-12px text-warning">
						{{ $t('contacts.subscribers.import.mapping.unmapped') }}: {{ previewData.unmapped.join(', ') }}
					</div>
				</n-card>

				<!-- Preview table -->
				<n-card size="small" :title="$t('contacts.subscribers.import.preview.title')">
					<n-data-table
						:columns="previewColumns"
						:data="previewData.preview_rows"
						:bordered="true"
						size="small"
						:max-height="300"
						:scroll-x="previewData.headers.length * 160"
					/>
				</n-card>

				<div class="flex justify-between mt-16px">
					<n-button @click="currentStep = 1">{{ $t('contacts.subscribers.import.actions.back') }}</n-button>
					<n-button type="primary" :disabled="!hasEmailMapped" @click="currentStep = 3">
						{{ $t('contacts.subscribers.import.actions.next') }}
					</n-button>
				</div>
			</div>

			<!-- STEP 3: Import Settings -->
			<div v-if="currentStep === 3">
				<n-alert type="info" class="mb-16px">
					{{ $t('contacts.subscribers.import.summary.ready', {
						count: previewData.total_rows - previewData.duplicates - previewData.invalid_rows,
						skipped: previewData.duplicates + previewData.invalid_rows
					}) }}
				</n-alert>

				<bt-form ref="settingsFormRef" :model="form" :rules="settingsRules">
					<n-form-item :label="$t('contacts.subscribers.import.overwriteLabel')">
						<div class="flex-1">
							<div class="mb-8px text-desc">{{ $t('contacts.subscribers.import.overwriteDescription') }}</div>
							<div class="w-260px">
								<n-select v-model:value="form.overwrite" :options="overwriteOptions" />
							</div>
						</div>
					</n-form-item>
					<n-grid :cols="2" :x-gap="24">
						<n-form-item-gi :span="1" :label="$t('contacts.subscribers.import.subscriptionStatus')">
							<n-select v-model:value="form.default_active" :options="modeOptions" />
						</n-form-item-gi>
						<n-form-item-gi :span="1" :label="$t('contacts.subscribers.import.status')">
							<n-select v-model:value="form.status" :options="statusOptions" />
						</n-form-item-gi>
					</n-grid>
					<n-form-item :label="$t('contacts.subscribers.import.group')" path="group_id">
						<group-select v-model:value="form.group_ids" />
					</n-form-item>
				</bt-form>

				<!-- Import result -->
				<n-result v-if="importResult" :status="importResult.status" :title="importResult.title" :description="importResult.description" class="mt-16px" />

				<div class="flex justify-between mt-16px">
					<n-button :disabled="importing" @click="currentStep = 2">{{ $t('contacts.subscribers.import.actions.back') }}</n-button>
					<div class="flex gap-12px">
						<n-button @click="handleCancel">{{ $t('common.actions.cancel') }}</n-button>
						<n-button type="primary" :loading="importing" :disabled="!!importResult" @click="handleImport">
							{{ $t('contacts.subscribers.import.actions.import') }}
						</n-button>
					</div>
				</div>
			</div>
		</div>
	</modal>
</template>

<script lang="ts" setup>
import { FormRules, UploadFileInfo, DataTableColumns } from 'naive-ui'
import { isEmpty } from 'lodash-es'
import { Message } from '@/utils'
import { useModal } from '@/hooks/modal/useModal'
import { downloadFile } from '@/api/modules/public'
import { importSubscribers, previewImport } from '@/api/modules/contacts/subscribers'

import GroupSelect from './GroupMultipleSelect.vue'

const { t } = useI18n()

const formRef = useTemplateRef('formRef')
const settingsFormRef = useTemplateRef('settingsFormRef')

// -- Step control --
const currentStep = ref(1)
const previewLoading = ref(false)
const importing = ref(false)
const importResult = ref<{ status: 'success' | 'error'; title: string; description: string } | null>(null)

// -- Preview data from API --
const previewData = reactive({
	headers: [] as string[],
	mapping: {} as Record<number, string>,
	unmapped: [] as string[],
	preview_rows: [] as Record<string, string>[],
	total_rows: 0,
	duplicates: 0,
	invalid_rows: 0,
})

// -- Column mapping (user-adjustable) --
const columnMapping = reactive<Record<number, string>>({})

// -- Field options for mapping --
const fieldOptions = computed(() => [
	{ label: 'Email', value: 'email' },
	{ label: 'First Name', value: 'first_name' },
	{ label: 'Last Name', value: 'last_name' },
	{ label: 'Company', value: 'company' },
	{ label: 'Phone', value: 'phone' },
	{ label: 'Title', value: 'title' },
	{ label: 'Website', value: 'website' },
	{ label: 'City', value: 'city' },
	{ label: 'Country', value: 'country' },
	{ label: 'Industry', value: 'industry' },
	{ label: 'LinkedIn', value: 'linkedin' },
	{ label: t('contacts.subscribers.import.mapping.skip'), value: 'skip' },
])

// -- Check at least one column is mapped to email --
const hasEmailMapped = computed(() => {
	return Object.values(columnMapping).includes('email')
})

// -- Form data --
const form = reactive({
	group_ids: [] as number[],
	import_type: 1,
	contacts: '',
	file_data: '',
	file_type: '',
	overwrite: 0,
	default_active: 1,
	status: 1,
})

// -- Derived: does user have data ready? --
const hasData = computed(() => {
	if (form.import_type === 1) return !isEmpty(form.file_data)
	return !isEmpty(form.contacts)
})

// -- Rules for step 1 --
const rules: FormRules = {
	contacts: {
		validator: () => {
			if (form.import_type === 2 && isEmpty(form.contacts)) {
				return new Error(t('contacts.subscribers.import.validation.dataRequired'))
			}
			return true
		},
	},
}

// -- Rules for step 3 --
const settingsRules: FormRules = {
	group_id: {
		trigger: 'change',
		validator: () => {
			if (form.group_ids.length === 0) {
				return new Error(t('contacts.subscribers.import.validation.groupRequired'))
			}
			return true
		},
	},
}

const overwriteOptions = computed(() => [
	{ label: t('contacts.subscribers.import.overwriteOptions.no'), value: 0 },
	{ label: t('contacts.subscribers.import.overwriteOptions.yes'), value: 1 },
])

const modeOptions = computed(() => [
	{ label: t('contacts.subscribers.import.subscriptionOptions.subscribe'), value: 1 },
	{ label: t('contacts.subscribers.import.subscriptionOptions.unsubscribe'), value: 0 },
])

const statusOptions = computed(() => [
	{ label: t('contacts.subscribers.import.statusOptions.confirmed'), value: 1 },
	{ label: t('contacts.subscribers.import.statusOptions.unconfirmed'), value: 0 },
])

const example = `email,first_name,last_name,company,phone,title
john@example.com,John,Doe,Acme Inc,+1234567890,CEO
jane@example.com,Jane,Smith,Globex,+0987654321,CTO
`

// -- Preview table columns (dynamic) --
const previewColumns = computed<DataTableColumns>(() => {
	if (!previewData.headers.length) return []
	return previewData.headers.map((h) => ({
		title: h,
		key: h,
		ellipsis: true,
		width: 160,
		render: (row: Record<string, string>) => row[h] ?? '',
	}))
})

// -- File upload handler --
const handleChangeFile = (file: UploadFileInfo) => {
	const reader = new FileReader()
	reader.onload = e => {
		form.file_type = file.name.substring(file.name.lastIndexOf('.') + 1)
		form.file_data = e.target?.result as string
	}
	if (file.file) {
		reader.readAsText(file.file)
	}
}

const handleDownloadTemplate = async () => {
	await downloadFile({ file_path: './template/example_recipients.csv' })
}

// -- Step 1 -> Step 2: Call preview API --
const handlePreview = async () => {
	if (!hasData.value) {
		Message.error(t('contacts.subscribers.import.validation.dataRequired'))
		return
	}
	try {
		previewLoading.value = true
		const payload: Record<string, unknown> = {
			import_type: form.import_type,
		}
		if (form.import_type === 1) {
			payload.file_data = form.file_data
			payload.file_type = form.file_type
		} else {
			payload.file_data = form.contacts
			payload.file_type = 'csv'
		}
		const res = await previewImport(payload)
		if (res && typeof res === 'object') {
			const data = res as {
				headers: string[]
				mapping: Record<number, string>
				unmapped: string[]
				preview_rows: Record<string, string>[]
				total_rows: number
				duplicates: number
				invalid_rows: number
			}
			previewData.headers = data.headers || []
			previewData.mapping = data.mapping || {}
			previewData.unmapped = data.unmapped || []
			previewData.preview_rows = data.preview_rows || []
			previewData.total_rows = data.total_rows || 0
			previewData.duplicates = data.duplicates || 0
			previewData.invalid_rows = data.invalid_rows || 0

			// Initialize column mapping from API auto-mapping
			Object.keys(columnMapping).forEach(k => delete columnMapping[Number(k)])
			previewData.headers.forEach((_, idx) => {
				columnMapping[idx] = previewData.mapping[idx] || 'skip'
			})

			currentStep.value = 2
		}
	} catch (e: unknown) {
		Message.error(t('contacts.subscribers.import.preview.error'))
	} finally {
		previewLoading.value = false
	}
}

// -- Step 3: Execute import --
const handleImport = async () => {
	await settingsFormRef.value?.validate()

	// Build column_mapping from user selections: { csv_col_index: field_name }
	const mappingPayload: Record<number, string> = {}
	previewData.headers.forEach((_, idx) => {
		if (columnMapping[idx] && columnMapping[idx] !== 'skip') {
			mappingPayload[idx] = columnMapping[idx]
		}
	})

	try {
		importing.value = true
		const payload: Record<string, unknown> = {
			...toRaw(form),
			column_mapping: mappingPayload,
		}
		await importSubscribers(payload as Parameters<typeof importSubscribers>[0])
		importResult.value = {
			status: 'success',
			title: t('contacts.subscribers.import.result.success.title'),
			description: t('contacts.subscribers.import.result.success.description', {
				count: previewData.total_rows - previewData.duplicates - previewData.invalid_rows,
			}),
		}
		const state = modalApi.getState<{ refresh: () => void }>()
		state?.refresh()
	} catch {
		importResult.value = {
			status: 'error',
			title: t('contacts.subscribers.import.result.error.title'),
			description: t('contacts.subscribers.import.result.error.description'),
		}
	} finally {
		importing.value = false
	}
}

const handleCancel = () => {
	modalApi.close()
}

// -- Modal lifecycle --
const [Modal, modalApi] = useModal({
	onChangeState: isOpen => {
		if (isOpen) {
			const state = modalApi.getState<{ group_id: number }>()
			if (state?.group_id) {
				form.group_ids = [state.group_id]
			} else {
				form.group_ids = []
			}
		} else {
			// Full reset
			currentStep.value = 1
			importResult.value = null
			previewData.headers = []
			previewData.mapping = {}
			previewData.unmapped = []
			previewData.preview_rows = []
			previewData.total_rows = 0
			previewData.duplicates = 0
			previewData.invalid_rows = 0
			Object.keys(columnMapping).forEach(k => delete columnMapping[Number(k)])
			form.import_type = 1
			form.contacts = ''
			form.file_data = ''
			form.file_type = ''
			form.group_ids = []
			form.overwrite = 0
			form.default_active = 1
			form.status = 1
		}
	},
})
</script>

<style lang="scss" scoped>
.csv-example {
	width: 100%;
	margin: 0;
	padding: 16px 24px;
	background: none;
	border: 1px solid var(--color-border-1);
	border-radius: 4px;
	overflow-x: auto;
	white-space: pre;
	word-wrap: normal;
}
</style>
