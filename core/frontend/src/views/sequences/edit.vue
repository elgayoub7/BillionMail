<template>
	<div class="seq-edit-page">
		<!-- Header -->
		<div class="edit-header">
			<div class="header-left">
				<n-button quaternary @click="router.back()">
					<i class="i-mdi:arrow-left"></i>
				</n-button>
				<div>
					<h1 class="header-title">{{ isEdit ? 'Edit Campaign' : 'New Campaign' }}</h1>
					<span class="header-sub">Configure your outreach sequence</span>
				</div>
			</div>
			<div class="header-actions">
				<n-button @click="router.back()">Cancel</n-button>
				<n-button @click="handleSave">Save Draft</n-button>
				<n-button type="primary" :loading="saving" @click="handleSaveAndActivate">
					<template #icon><i class="i-mdi:play"></i></template>
					Save & Activate
				</n-button>
			</div>
		</div>

		<div class="edit-body">
			<!-- Left: Config -->
			<div class="config-col">
				<!-- Basic Info -->
				<div class="config-section">
					<div class="section-title">General</div>
					<div class="form-grid">
						<div class="form-field">
							<label class="field-label">Campaign Name</label>
							<n-input v-model:value="form.name" placeholder="e.g. Q2 Outreach Round 1" />
						</div>
						<div class="form-field">
							<label class="field-label">Contact Group</label>
							<group-select v-model:value="form.group_id" />
						</div>
						<div class="form-field full">
							<label class="field-label">Description</label>
							<n-input v-model:value="form.description" type="textarea" :rows="2" placeholder="Internal notes..." />
						</div>
					</div>
				</div>

				<!-- Senders -->
				<div class="config-section">
					<div class="section-title">Senders</div>
					<div class="form-grid">
						<div class="form-field full">
							<label class="field-label">Sender Pool</label>
							<n-select
								v-model:value="form.selected_senders"
								:options="mailboxOptions"
								multiple
								filterable
								placeholder="Select sender mailboxes..."
								@update:value="onSenderPoolChange"
							/>
						</div>
						<div class="form-field">
							<label class="field-label">Daily Limit / Sender</label>
							<n-input-number v-model:value="form.daily_limit_per_sender" :min="0" placeholder="0 = unlimited" class="w-full" />
						</div>
						<div class="form-field">
							<label class="field-label">Display Name</label>
							<n-input v-model:value="form.full_name" placeholder="Auto-filled from mailbox" />
						</div>
					</div>
				</div>

				<!-- Schedule -->
				<div class="config-section">
					<div class="section-title">Schedule</div>
					<div class="form-grid">
						<div class="form-field">
							<label class="field-label">Delay Between Emails</label>
							<n-select v-model:value="form.send_delay" :options="delayOptions" placeholder="Select delay..." />
						</div>
						<div class="form-field">
							<label class="field-label">Sending Window</label>
							<div class="time-range">
								<n-input-number v-model:value="form.schedule_start_hour" :min="0" :max="23" size="small" />
								<span class="range-sep">to</span>
								<n-input-number v-model:value="form.schedule_end_hour" :min="1" :max="24" size="small" />
							</div>
						</div>
						<div class="form-field full">
							<label class="field-label">Active Days</label>
							<div class="days-toggle">
								<button
									v-for="d in dayOptions"
									:key="d.value"
									class="day-btn"
									:class="{ active: form.schedule_days.includes(d.value) }"
									@click="toggleDay(d.value)">
									{{ d.label }}
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Right: Steps Timeline -->
			<div class="steps-col">
				<div class="steps-header">
					<div class="section-title">Steps</div>
					<n-button size="small" @click="addStep">
						<template #icon><i class="i-mdi:plus"></i></template>
						Add Step
					</n-button>
				</div>

				<div class="steps-timeline">
					<div v-for="(step, index) in form.steps" :key="index" class="step-item">
						<!-- Connector -->
						<div v-if="index > 0" class="step-connector">
							<div class="connector-line"></div>
							<div class="connector-dot">
								<i class="i-mdi:arrow-down" style="font-size:12px;color:#6c5ce7"></i>
							</div>
							<div class="connector-line"></div>
						</div>

						<div class="step-card" :class="`step-${step.step_type}`">
							<div class="step-top">
								<div class="step-badge">
									<span class="step-number">{{ step.step_order }}</span>
									<i :class="getStepIcon(step.step_type)" class="step-type-icon"></i>
									<span class="step-type-label">{{ getStepLabel(step.step_type) }}</span>
								</div>
								<n-button quaternary size="tiny" type="error" @click="removeStep(index)" :disabled="form.steps.length <= 1">
									<i class="i-mdi:close"></i>
								</n-button>
							</div>

							<div class="step-body">
								<!-- Type selector -->
								<n-select
									v-model:value="step.step_type"
									:options="stepTypeOptions"
									size="small"
									@update:value="(val: string) => onStepTypeChange(step, val)"
								/>

								<!-- Email step -->
								<template v-if="step.step_type === 'email'">
									<div class="subject-wrapper">
										<n-input
											v-model:value="step.subject"
											size="small"
											placeholder="Email subject line..."
											@input="(val: string) => debounceScoreSubject(index, val)"
										/>
										<div v-if="stepScores[index] !== undefined" class="score-bar">
											<div class="score-fill" :style="{ width: stepScores[index] + '%', background: getScoreColor(stepScores[index]) }"></div>
											<span class="score-value" :style="{ color: getScoreColor(stepScores[index]) }">{{ stepScores[index] }}/100</span>
											<template v-if="stepWarnings[index]?.length">
												<i class="i-mdi:alert-circle-outline" style="color:#f59e0b;margin-left:4px"></i>
												<span class="score-warnings">{{ stepWarnings[index].join(', ') }}</span>
											</template>
										</div>
									</div>
									<n-select
										v-model:value="step.template_id"
										:options="templateOptions"
										size="small"
										filterable
										placeholder="Select template..."
										clearable
									/>
								</template>

								<!-- Wait step -->
								<template v-if="step.step_type === 'wait'">
									<div class="wait-inputs">
										<n-input-number v-model:value="step.wait_days" :min="0" size="small" placeholder="Days">
											<template #suffix>days</template>
										</n-input-number>
										<n-input-number v-model:value="step.wait_hours" :min="0" :max="23" size="small" placeholder="Hours">
											<template #suffix>hrs</template>
										</n-input-number>
									</div>
								</template>

								<!-- Condition step -->
								<template v-if="step.step_type === 'condition'">
									<n-select v-model:value="step.condition_type" :options="conditionOptions" size="small" />
									<n-input-number v-model:value="step.condition_step_id" :min="1" size="small" placeholder="Go to step #" class="w-full" />
								</template>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Message } from '@/utils'
import { getSequenceDetail, createSequence, updateSequence, activateSequence } from '@/api/modules/sequences/sequence'
import { scoreSubject } from '@/api/modules/batch_mail'
import { getTemplateAll } from '@/api/modules/market/template'
import GroupSelect from '@/views/contacts/subscribers/components/GroupSelect.vue'
import type { StepInput } from './interface'
import { instance } from '@/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const saving = ref(false)

const stepScores = reactive<Record<number, number>>({})
const stepWarnings = reactive<Record<number, string[]>>({})
const scoreTimers = reactive<Record<number, ReturnType<typeof setTimeout>>>({})

function debounceScoreSubject(index: number, val: string) {
	if (scoreTimers[index]) clearTimeout(scoreTimers[index])
	if (!val.trim()) {
		delete stepScores[index]
		delete stepWarnings[index]
		return
	}
	scoreTimers[index] = setTimeout(async () => {
		try {
			const res = await scoreSubject(val)
			stepScores[index] = res.score
			stepWarnings[index] = res.warnings || []
		} catch {
			delete stepScores[index]
			delete stepWarnings[index]
		}
	}, 500)
}

function getScoreColor(score: number): string {
	if (score >= 60) return '#22c55e'
	if (score >= 40) return '#f59e0b'
	return '#ef4444'
}

function getStepIcon(type: string) {
	if (type === 'email') return 'i-mdi:email-outline'
	if (type === 'wait') return 'i-mdi:clock-outline'
	return 'i-mdi:source-branch'
}

function getStepLabel(type: string) {
	if (type === 'email') return 'Email'
	if (type === 'wait') return 'Wait'
	return 'Condition'
}

const isEdit = computed(() => !!route.params.id)

const templateOptions = ref<any[]>([])
const mailboxOptions = ref<any[]>([])

const delayOptions = [
	{ label: 'No delay', value: 0 },
	{ label: '30 seconds', value: 30 },
	{ label: '1 minute', value: 60 },
	{ label: '2 minutes', value: 120 },
	{ label: '5 minutes', value: 300 },
	{ label: '10 minutes', value: 600 },
	{ label: '15 minutes', value: 900 },
	{ label: '20 minutes', value: 1200 },
	{ label: '30 minutes', value: 1800 },
	{ label: '1 hour', value: 3600 },
]

const dayOptions = [
	{ label: 'Mon', value: 1 },
	{ label: 'Tue', value: 2 },
	{ label: 'Wed', value: 3 },
	{ label: 'Thu', value: 4 },
	{ label: 'Fri', value: 5 },
	{ label: 'Sat', value: 6 },
	{ label: 'Sun', value: 7 },
]

function toggleDay(val: number) {
	const idx = form.schedule_days.indexOf(val)
	if (idx >= 0) form.schedule_days.splice(idx, 1)
	else form.schedule_days.push(val)
}

const form = reactive({
	name: '',
	description: '',
	addresser: '',
	full_name: '',
	group_id: 0,
	tag_ids: [] as number[],
	tag_logic: 'AND',
	track_open: 1,
	track_click: 1,
	unsubscribe: 1,
	selected_senders: [] as string[],
	sender_pool: '[]',
	daily_limit_per_sender: 0,
	send_delay: 0,
	schedule_start_hour: 0,
	schedule_end_hour: 24,
	schedule_days: [1, 2, 3, 4, 5, 6, 7] as number[],
	steps: [
		{ step_order: 1, step_type: 'email', subject: '', template_id: null as number | null },
	] as StepInput[],
})

async function loadTemplates() {
	try {
		const res = await getTemplateAll()
		const list = Array.isArray(res) ? res : (res?.data || res?.list || [])
		templateOptions.value = list.map((t: any) => ({
			label: t.temp_name || t.name || `Template #${t.id}`,
			value: t.id,
		}))
	} catch { /* ignore */ }
}

async function loadMailboxes() {
	try {
		const res = await instance.get('/mailbox/all')
		const list = Array.isArray(res) ? res : (res?.data || [])
		mailboxOptions.value = list.map((m: any) => ({
			label: `${m.full_name || m.username} <${m.username}>`,
			value: m.username,
			email: m.username,
			full_name: m.full_name || '',
			domain: m.domain || '',
		}))
	} catch { /* ignore */ }
}

function onSenderPoolChange(ids: string[]) {
	if (ids.length === 0) {
		form.sender_pool = '[]'
		form.addresser = ''
		form.full_name = ''
		return
	}
	const selected = mailboxOptions.value.filter(m => ids.includes(m.value as unknown as string))
	const pool = selected.map(m => ({ email: m.email, name: m.full_name || m.email.split('@')[0] }))
	form.sender_pool = JSON.stringify(pool)
	form.addresser = pool[0].email
	form.full_name = pool[0].name
}

const stepTypeOptions = [
	{ label: t('sequences.stepType.email'), value: 'email' },
	{ label: t('sequences.stepType.wait'), value: 'wait' },
	{ label: t('sequences.stepType.condition'), value: 'condition' },
]

const conditionOptions = [
	{ label: t('sequences.condition.opened'), value: 'opened' },
	{ label: t('sequences.condition.notOpened'), value: 'not_opened' },
	{ label: t('sequences.condition.clicked'), value: 'clicked' },
	{ label: t('sequences.condition.notClicked'), value: 'not_clicked' },
	{ label: t('sequences.condition.bounced'), value: 'bounced' },
]

function addStep() {
	const nextOrder = form.steps.length + 1
	form.steps.push({ step_order: nextOrder, step_type: 'email', subject: '', template_id: null as number | null })
}

function removeStep(index: number) {
	form.steps.splice(index, 1)
	form.steps.forEach((step, i) => {
		step.step_order = i + 1
	})
}

function onStepTypeChange(step: StepInput, type: string) {
	if (type === 'wait') {
		step.wait_days = step.wait_days || 3
		step.wait_hours = 0
	} else if (type === 'email') {
		step.template_id = null as any
		step.subject = ''
	} else if (type === 'condition') {
		step.condition_type = 'opened'
		step.condition_step_id = 1
	}
}

function buildPayload() {
	let addr = form.addresser
	let fn = form.full_name
	if (form.sender_pool && form.sender_pool !== '[]') {
		try {
			const pool = JSON.parse(form.sender_pool)
			if (pool.length > 0) {
				addr = pool[0].email
				fn = pool[0].name
			}
		} catch {}
	}
	return {
		name: form.name,
		description: form.description,
		addresser: addr,
		full_name: fn,
		group_id: form.group_id,
		tag_ids: form.tag_ids,
		tag_logic: form.tag_logic,
		track_open: form.track_open,
		track_click: form.track_click,
		unsubscribe: form.unsubscribe,
		sender_pool: form.sender_pool,
		daily_limit_per_sender: form.daily_limit_per_sender,
		send_delay: form.send_delay,
		schedule_start_hour: form.schedule_start_hour,
		schedule_end_hour: form.schedule_end_hour,
		schedule_days: JSON.stringify(form.schedule_days),
		steps: form.steps.map(s => ({ ...s, template_id: s.template_id || 0 })),
	}
}

async function handleSave() {
	if (!form.name) {
		Message.warning(t('sequences.form.nameRequired'))
		return
	}
	saving.value = true
	try {
		const payload = buildPayload()
		if (isEdit.value) {
			await updateSequence({ id: Number(route.params.id), ...payload })
		} else {
			await createSequence(payload as any)
		}
		Message.success(t('sequences.form.saved'))
		router.push('/sequences')
	} catch {} finally {
		saving.value = false
	}
}

async function handleSaveAndActivate() {
	if (!form.name) {
		Message.warning(t('sequences.form.nameRequired'))
		return
	}
	saving.value = true
	try {
		const payload = buildPayload()
		let id: number
		if (isEdit.value) {
			await updateSequence({ id: Number(route.params.id), ...payload })
			id = Number(route.params.id)
		} else {
			const res = await createSequence(payload as any)
			id = res.id
		}
		await activateSequence({ id })
		Message.success(t('sequences.form.activated'))
		router.push('/sequences')
	} catch {} finally {
		saving.value = false
	}
}

onMounted(async () => {
	loadTemplates()
	loadMailboxes()

	if (isEdit.value) {
		const res = await getSequenceDetail({ id: Number(route.params.id) })
		const detail = res
		let sp: any[] = []
		try { sp = detail.sender_pool ? JSON.parse(detail.sender_pool) : [] } catch {}
		let sd: number[] = [1,2,3,4,5,6,7]
		try { sd = detail.schedule_days ? JSON.parse(detail.schedule_days) : [1,2,3,4,5,6,7] } catch {}

		Object.assign(form, {
			name: detail.name,
			description: detail.description,
			addresser: detail.addresser,
			full_name: detail.full_name,
			group_id: detail.group_id,
			tag_ids: detail.tag_ids || [],
			tag_logic: detail.tag_logic,
			track_open: detail.track_open,
			track_click: detail.track_click,
			unsubscribe: detail.unsubscribe,
			sender_pool: detail.sender_pool || '[]',
			daily_limit_per_sender: detail.daily_limit_per_sender || 0,
			send_delay: detail.send_delay || 0,
			schedule_start_hour: detail.schedule_start_hour ?? 0,
			schedule_end_hour: detail.schedule_end_hour ?? 24,
			schedule_days: sd,
			steps: detail.steps.map((s: any) => ({
				step_order: s.step_order,
				step_type: s.step_type,
				subject: s.subject || '',
				template_id: s.template_id || null,
				wait_days: s.wait_days || 0,
				wait_hours: s.wait_hours || 0,
				condition_type: s.condition_type || '',
				condition_step_id: s.condition_step_id || 0,
				on_true_go_to: s.on_true_go_to || 0,
				on_false_go_to: s.on_false_go_to || 0,
			})),
		})

		if (sp.length > 0) {
			const interval = setInterval(() => {
				if (mailboxOptions.value.length === 0) return
				clearInterval(interval)
				const ids = mailboxOptions.value
					.filter(m => sp.some((s: any) => s.email === m.email))
					.map(m => m.value as unknown as string)
				form.selected_senders = ids
			}, 300)
		}
	}
})
</script>

<style lang="scss" scoped>
.seq-edit-page {
	padding: 24px;
	min-height: 100vh;
}

.edit-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24px;
}

.header-left {
	display: flex;
	align-items: center;
	gap: 12px;
}

.header-title {
	margin: 0;
	font-size: 20px;
	font-weight: 700;
	color: #e2e8f0;
}

.header-sub {
	font-size: 13px;
	color: #6b7280;
}

.header-actions {
	display: flex;
	gap: 8px;
}

.edit-body {
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 24px;
	align-items: start;
}

/* Config column */
.config-col {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

.config-section {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
}

.section-title {
	font-size: 14px;
	font-weight: 600;
	color: #e2e8f0;
	margin-bottom: 16px;
	text-transform: uppercase;
	letter-spacing: 0.5px;
}

.form-grid {
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 12px;
}

.form-field {
	display: flex;
	flex-direction: column;
	gap: 6px;

	&.full {
		grid-column: 1 / -1;
	}
}

.field-label {
	font-size: 12px;
	font-weight: 500;
	color: #8892a8;
}

.time-range {
	display: flex;
	align-items: center;
	gap: 8px;

	.range-sep {
		color: #6b7280;
		font-size: 13px;
	}
}

.days-toggle {
	display: flex;
	gap: 6px;
}

.day-btn {
	padding: 6px 12px;
	border: 1px solid #2e3142;
	border-radius: 6px;
	background: transparent;
	color: #6b7280;
	font-size: 12px;
	font-weight: 500;
	cursor: pointer;
	transition: all 0.15s;

	&:hover {
		border-color: #6c5ce7;
		color: #e2e8f0;
	}

	&.active {
		background: rgba(108, 92, 231, 0.15);
		border-color: #6c5ce7;
		color: #a78bfa;
	}
}

/* Steps column */
.steps-col {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.steps-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.steps-timeline {
	display: flex;
	flex-direction: column;
}

.step-item {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.step-connector {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 4px 0;
}

.connector-line {
	width: 2px;
	height: 8px;
	background: #2e3142;
}

.connector-dot {
	width: 20px;
	height: 20px;
	border-radius: 50%;
	background: #1a1d27;
	border: 2px solid #2e3142;
	display: flex;
	align-items: center;
	justify-content: center;
}

.step-card {
	width: 100%;
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 16px;
	transition: border-color 0.2s;

	&:hover {
		border-color: #3d4058;
	}

	&.step-email {
		border-left: 3px solid #6c5ce7;
	}

	&.step-wait {
		border-left: 3px solid #f59e0b;
	}

	&.step-condition {
		border-left: 3px solid #22c55e;
	}
}

.step-top {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 12px;
}

.step-badge {
	display: flex;
	align-items: center;
	gap: 8px;
}

.step-number {
	width: 22px;
	height: 22px;
	border-radius: 50%;
	background: rgba(108, 92, 231, 0.15);
	color: #a78bfa;
	font-size: 11px;
	font-weight: 700;
	display: flex;
	align-items: center;
	justify-content: center;
}

.step-type-icon {
	font-size: 16px;
	color: #8892a8;
}

.step-type-label {
	font-size: 13px;
	font-weight: 600;
	color: #e2e8f0;
}

.step-body {
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.subject-wrapper {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.score-bar {
	display: flex;
	align-items: center;
	gap: 8px;
	height: 16px;
}

.score-fill {
	height: 4px;
	border-radius: 2px;
	transition: width 0.3s;
	flex: 0 0 auto;
	min-width: 20px;
	max-width: 100px;
}

.score-value {
	font-size: 11px;
	font-weight: 600;
}

.score-warnings {
	font-size: 11px;
	color: #f59e0b;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.wait-inputs {
	display: flex;
	gap: 8px;
}

.w-full { width: 100%; }

:deep(.n-input),
:deep(.n-select),
:deep(.n-input-number) {
	--n-border: 1px solid #2e3142;
	--n-border-hover: 1px solid #6c5ce7;
	--n-border-focus: 1px solid #6c5ce7;
	--n-color: #0f1117;
	--n-color-active: #0f1117;
	--n-text-color: #e2e8f0;
	--n-placeholder-color: #4a5568;
	--n-caret-color: #a78bfa;
}
</style>
