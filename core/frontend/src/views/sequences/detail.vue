<template>
	<div class="seq-detail">
		<!-- Header -->
		<div class="detail-header">
			<div class="header-left">
				<n-button quaternary @click="router.back()">
					<i class="i-mdi:arrow-left" style="font-size:18px"></i>
				</n-button>
				<div>
					<h1 class="detail-name">{{ detail?.name || 'Loading...' }}</h1>
					<span v-if="detail?.description" class="detail-desc">{{ detail.description }}</span>
				</div>
				<n-tag v-if="detail" :type="statusMap[detail.status]?.type" size="small" round>
					{{ statusMap[detail.status]?.label }}
				</n-tag>
			</div>
			<div class="header-actions">
				<n-button @click="router.push(`/sequences/${route.params.id}/edit`)">
					<i class="i-mdi:pencil-outline"></i> Edit
				</n-button>
				<n-button v-if="detail?.status === 0" type="success" @click="handleActivate">Activate</n-button>
				<n-button v-if="detail?.status === 1" type="warning" @click="handlePause">Pause</n-button>
				<n-button v-if="detail?.status === 2" type="success" @click="handleResume">Resume</n-button>
				<n-popconfirm @positive-click="handleDelete">
					<template #trigger>
						<n-button type="error">Delete</n-button>
					</template>
					Delete this campaign?
				</n-popconfirm>
			</div>
		</div>

		<template v-if="detail">
			<!-- KPI Row -->
			<div class="kpi-row">
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#6c5ce7"><i class="i-mdi:send-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ totalSent }}</div>
						<div class="kpi-lbl">Emails Sent</div>
					</div>
				</div>
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#00d68f"><i class="i-mdi:email-open-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ totalOpened }}</div>
						<div class="kpi-lbl">Opened ({{ openRate }}%)</div>
					</div>
				</div>
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#0095ff"><i class="i-mdi:cursor-default-click-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ totalClicked }}</div>
						<div class="kpi-lbl">Clicked ({{ clickRate }}%)</div>
					</div>
				</div>
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#ffaa00"><i class="i-mdi:reply-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ totalReplied }}</div>
						<div class="kpi-lbl">Replied ({{ replyRate }}%)</div>
					</div>
				</div>
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#ff3d71"><i class="i-mdi:email-alert-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ detail.total_bounced || 0 }}</div>
						<div class="kpi-lbl">Bounced ({{ bounceRate }}%)</div>
					</div>
				</div>
				<div class="kpi-card">
					<span class="kpi-icon" style="color:#6c5ce7"><i class="i-mdi:account-group-outline"></i></span>
					<div class="kpi-data">
						<div class="kpi-val">{{ detail.total_enrolled }}</div>
						<div class="kpi-lbl">Enrolled</div>
					</div>
				</div>
			</div>

			<!-- Steps Timeline -->
			<div class="section-title">
				<span>Steps</span>
				<n-tag size="small" :bordered="false">{{ detail.steps?.length || 0 }} steps</n-tag>
			</div>
			<div class="steps-timeline">
				<div v-for="(step, idx) in detail.steps" :key="step.id" class="step-node">
					<div class="step-connector" v-if="idx > 0">
						<div class="connector-line"></div>
						<i class="i-mdi:chevron-down connector-arrow"></i>
					</div>
					<div class="step-card" :class="'step-' + step.step_type">
						<div class="step-header">
							<div class="step-badge">
								<i v-if="step.step_type === 'email'" class="i-mdi:email-outline"></i>
								<i v-else-if="step.step_type === 'wait'" class="i-mdi:clock-outline"></i>
								<i v-else class="i-mdi:source-branch"></i>
							</div>
							<div class="step-meta">
								<span class="step-order">Step {{ step.step_order }}</span>
								<n-tag size="small" :type="stepTypeColor(step.step_type)" round>{{ step.step_type }}</n-tag>
							</div>
						</div>
						<template v-if="step.step_type === 'email'">
							<div class="step-subject">{{ step.subject || 'No subject' }}</div>
							<div class="step-stats">
								<div class="st"><span class="st-v">{{ step.sent_count }}</span><span class="st-l">Sent</span></div>
								<div class="st"><span class="st-v">{{ step.opened_count }}</span><span class="st-l">Opened</span></div>
								<div class="st"><span class="st-v">{{ step.clicked_count }}</span><span class="st-l">Clicked</span></div>
								<div class="st"><span class="st-v">{{ step.bounced_count }}</span><span class="st-l">Bounced</span></div>
							</div>
						</template>
						<template v-if="step.step_type === 'wait'">
							<div class="step-wait-info">
								<i class="i-mdi:timer-sand"></i>
								Wait {{ step.wait_days || 0 }}d {{ step.wait_hours || 0 }}h
							</div>
						</template>
						<template v-if="step.step_type === 'condition'">
							<div class="step-condition-info">
								If <strong>{{ step.condition_type }}</strong>
								<span v-if="step.condition_step_id"> → Step {{ step.condition_step_id }}</span>
							</div>
						</template>
					</div>
				</div>
			</div>

			<!-- Tabs: Enrollments + Lead Scoring -->
			<n-tabs type="line" class="detail-tabs" v-model:value="activeTab">
				<n-tab-pane name="enrollments" tab="Enrollments">
					<div class="tab-toolbar">
						<n-select v-model:value="enrollmentStatus" :options="enrollmentStatusOptions" style="width:160px" @update:value="fetchEnrollments" />
						<n-button type="primary" @click="handleEnroll">
							<i class="i-mdi:account-plus-outline"></i> Enroll Contacts
						</n-button>
					</div>
					<n-data-table :columns="enrollmentColumns" :data="enrollments" :loading="enrollmentsLoading" :pagination="enrollmentPagination" />
				</n-tab-pane>
				<n-tab-pane name="scoring" tab="Lead Scoring">
					<div class="scoring-cards">
						<div class="scoring-card hot">
							<div class="sc-val">{{ hotCount }}</div>
							<div class="sc-label">Hot Leads</div>
							<div class="sc-threshold">Score &gt; 40</div>
						</div>
						<div class="scoring-card warm">
							<div class="sc-val">{{ warmCount }}</div>
							<div class="sc-label">Warm Leads</div>
							<div class="sc-threshold">Score 10-40</div>
						</div>
						<div class="scoring-card cold">
							<div class="sc-val">{{ coldCount }}</div>
							<div class="sc-label">Cold Leads</div>
							<div class="sc-threshold">Score &lt; 10</div>
						</div>
						<div class="scoring-card avg">
							<div class="sc-val">{{ avgScore.toFixed(1) }}</div>
							<div class="sc-label">Avg Score</div>
						</div>
					</div>
					<n-data-table :columns="scoringColumns" :data="scoredEnrollments" :pagination="{ pageSize: 10 }" :default-sort="{ columnKey: 'score', order: 'descend' }" />
				</n-tab-pane>
			</n-tabs>
		</template>
	</div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NTag, NPopconfirm, NButton, NProgress } from 'naive-ui'
import { Message } from '@/utils'
import { getSequenceDetail, getEnrollments, enrollContacts, removeEnrollment, activateSequence, pauseSequence, resumeSequence, deleteSequence } from '@/api/modules/sequences/sequence'
import { formatTime } from '@/utils'
import type { SequenceDetail, Enrollment } from './interface'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const detail = ref<SequenceDetail | null>(null)
const enrollments = ref<Enrollment[]>([])
const enrollmentsLoading = ref(false)
const enrollmentStatus = ref(-1)
const enrollmentPage = ref(1)
const enrollmentPageSize = ref(10)
const activeTab = ref('enrollments')

const statusMap: Record<number, { label: string; type: 'default' | 'success' | 'warning' | 'error' }> = {
	0: { label: 'Draft', type: 'default' },
	1: { label: 'Active', type: 'success' },
	2: { label: 'Paused', type: 'warning' },
	3: { label: 'Archived', type: 'error' },
}

const enrollmentStatusOptions = [
	{ label: t('common.all.text'), value: -1 },
	{ label: t('sequences.enrollStatus.active'), value: 0 },
	{ label: t('sequences.enrollStatus.completed'), value: 1 },
	{ label: t('sequences.enrollStatus.paused'), value: 2 },
	{ label: t('sequences.enrollStatus.exited'), value: 3 },
]

const enrollmentPagination = computed(() => ({
	page: enrollmentPage.value,
	pageSize: enrollmentPageSize.value,
	itemCount: 0,
	onChange: (page: number) => {
		enrollmentPage.value = page
		fetchEnrollments()
	},
}))

function stepTypeColor(type: string): 'info' | 'success' | 'warning' {
	if (type === 'email') return 'info'
	if (type === 'wait') return 'warning'
	return 'success'
}

// Computed KPIs
const totalSent = computed(() => detail.value?.steps?.reduce((s, st) => s + (st.sent_count || 0), 0) || 0)
const totalOpened = computed(() => detail.value?.steps?.reduce((s, st) => s + (st.opened_count || 0), 0) || 0)
const totalClicked = computed(() => detail.value?.steps?.reduce((s, st) => s + (st.clicked_count || 0), 0) || 0)
const totalReplied = computed(() => enrollments.value.reduce((s, e) => s + (e.total_replies || 0), 0))
const openRate = computed(() => totalSent.value ? ((totalOpened.value / totalSent.value) * 100).toFixed(1) : '0')
const clickRate = computed(() => totalSent.value ? ((totalClicked.value / totalSent.value) * 100).toFixed(1) : '0')
const replyRate = computed(() => totalSent.value ? ((totalReplied.value / totalSent.value) * 100).toFixed(1) : '0')
const bounceRate = computed(() => {
	const b = detail.value?.total_bounced || 0
	return totalSent.value ? ((b / totalSent.value) * 100).toFixed(1) : '0'
})

// Actions
async function handleActivate() {
	await activateSequence({ id: Number(route.params.id) })
	Message.success('Campaign activated')
	loadDetail()
}
async function handlePause() {
	await pauseSequence({ id: Number(route.params.id) })
	Message.success('Campaign paused')
	loadDetail()
}
async function handleResume() {
	await resumeSequence({ id: Number(route.params.id) })
	Message.success('Campaign resumed')
	loadDetail()
}
async function handleDelete() {
	await deleteSequence({ id: Number(route.params.id) })
	Message.success('Campaign deleted')
	router.push('/sequences')
}

// --- Lead Scoring ---
interface ScoredEnrollment extends Enrollment {
	score: number
	tier: 'hot' | 'warm' | 'cold'
}

function calculateLeadScore(e: Enrollment): number {
	return e.total_opens * 5 + e.total_clicks * 10 + e.total_replies * 30
}

function getScoreTier(score: number): 'hot' | 'warm' | 'cold' {
	if (score >= 40) return 'hot'
	if (score >= 10) return 'warm'
	return 'cold'
}

const scoredEnrollments = computed<ScoredEnrollment[]>(() =>
	enrollments.value.map((e) => {
		const score = calculateLeadScore(e)
		return { ...e, score, tier: getScoreTier(score) }
	}).sort((a, b) => b.score - a.score)
)

const hotCount = computed(() => scoredEnrollments.value.filter((e) => e.tier === 'hot').length)
const warmCount = computed(() => scoredEnrollments.value.filter((e) => e.tier === 'warm').length)
const coldCount = computed(() => scoredEnrollments.value.filter((e) => e.tier === 'cold').length)
const avgScore = computed(() => {
	if (scoredEnrollments.value.length === 0) return 0
	return scoredEnrollments.value.reduce((sum, e) => sum + e.score, 0) / scoredEnrollments.value.length
})

const tierConfig: Record<string, { label: string; type: 'error' | 'warning' | 'info' }> = {
	hot: { label: 'Hot', type: 'error' },
	warm: { label: 'Warm', type: 'warning' },
	cold: { label: 'Cold', type: 'info' },
}

const scoringColumns = computed(() => [
	{ key: 'email', title: 'Email', minWidth: 220, ellipsis: { tooltip: true } },
	{
		key: 'score', title: 'Score', width: 90, sorter: 'default' as const, defaultSortOrder: 'descend' as const,
		render: (row: ScoredEnrollment) => {
			const pct = Math.min((row.score / 100) * 100, 100)
			const status = row.score >= 40 ? 'error' : row.score >= 10 ? 'warning' : 'success'
			return h(NProgress, { type: 'line', percentage: pct, status, indicatorPlacement: 'inside', showIndicator: true, style: { width: '80px' } })
		},
	},
	{
		key: 'tier', title: 'Tier', width: 100,
		render: (row: ScoredEnrollment) => h(NTag, { size: 'small', type: tierConfig[row.tier].type }, { default: () => tierConfig[row.tier].label }),
	},
	{ key: 'total_opens', title: 'Opens', width: 80 },
	{ key: 'total_clicks', title: 'Clicks', width: 80 },
	{ key: 'total_replies', title: 'Replies', width: 80 },
	{ key: 'total_emails_sent', title: 'Emails', width: 80 },
])

// --- Enrollments ---
const enrollmentStatusMap: Record<number, { label: string; type: 'default' | 'success' | 'warning' | 'error' }> = {
	0: { label: 'sequences.enrollStatus.active', type: 'success' },
	1: { label: 'sequences.enrollStatus.completed', type: 'default' },
	2: { label: 'sequences.enrollStatus.paused', type: 'warning' },
	3: { label: 'sequences.enrollStatus.exited', type: 'error' },
}

const enrollmentColumns = computed(() => [
	{ key: 'email', title: t('sequences.enrollments.email'), minWidth: 200, ellipsis: { tooltip: true } },
	{ key: 'current_step', title: t('sequences.enrollments.step'), width: 100 },
	{
		key: 'status', title: t('sequences.enrollments.status'), width: 120,
		render: (row: Enrollment) => {
			const s = enrollmentStatusMap[row.status]
			return s ? h(NTag, { size: 'small', type: s.type }, { default: () => t(s.label) }) : row.status
		},
	},
	{ key: 'total_emails_sent', title: t('sequences.enrollments.emailsSent'), width: 120 },
	{ key: 'total_opens', title: t('sequences.enrollments.opens'), width: 90 },
	{ key: 'total_clicks', title: t('sequences.enrollments.clicks'), width: 90 },
	{
		key: 'enrolled_at', title: t('sequences.enrollments.enrolledAt'), width: 170,
		render: (row: Enrollment) => formatTime(row.enrolled_at),
	},
	{
		title: t('common.columns.actions'), key: 'actions', width: 100,
		render: (row: Enrollment) => h(NPopconfirm, { onPositiveClick: () => handleRemoveEnrollment(row.id) }, {
			trigger: () => h(NButton, { type: 'error', text: true, size: 'small' }, { default: () => t('common.actions.delete') }),
			default: () => t('sequences.enrollments.removeConfirm'),
		}),
	},
])

async function fetchEnrollments() {
	if (!detail.value) return
	enrollmentsLoading.value = true
	try {
		const res = await getEnrollments({
			sequence_id: detail.value.id,
			page: enrollmentPage.value,
			page_size: enrollmentPageSize.value,
			status: enrollmentStatus.value,
		})
		enrollments.value = res.list
		enrollmentPagination.value.itemCount = res.total
	} finally {
		enrollmentsLoading.value = false
	}
}

async function handleEnroll() {
	await enrollContacts({ sequence_id: Number(route.params.id) })
	Message.success(t('sequences.enrollSuccess'))
	fetchEnrollments()
	loadDetail()
}

async function handleRemoveEnrollment(id: number) {
	await removeEnrollment({ enrollment_id: id })
	Message.success(t('sequences.removeSuccess'))
	fetchEnrollments()
}

async function loadDetail() {
	const res = await getSequenceDetail({ id: Number(route.params.id) })
	detail.value = res
}

onMounted(async () => {
	await loadDetail()
	fetchEnrollments()
})
</script>

<style lang="scss" scoped>
.seq-detail { padding: 24px; }

.detail-header {
	display: flex; justify-content: space-between; align-items: center;
	margin-bottom: 24px; gap: 16px;
}
.header-left { display: flex; align-items: center; gap: 12px; }
.detail-name { font-size: 22px; font-weight: 700; color: #e2e8f0; margin: 0; }
.detail-desc { font-size: 13px; color: #6b7280; display: block; }
.header-actions { display: flex; gap: 8px; }

.kpi-row { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; margin-bottom: 28px; }
.kpi-card {
	display: flex; align-items: center; gap: 12px; padding: 16px;
	background: #1a1d27; border: 1px solid #2e3142; border-radius: 10px;
}
.kpi-icon { font-size: 22px; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; background: rgba(108,92,231,0.08); border-radius: 10px; }
.kpi-val { font-size: 20px; font-weight: 700; color: #f1f3f7; }
.kpi-lbl { font-size: 12px; color: #6b7280; margin-top: 2px; }

.section-title {
	display: flex; align-items: center; gap: 10px;
	font-size: 16px; font-weight: 600; color: #e2e8f0; margin-bottom: 16px;
}

.steps-timeline { display: flex; flex-direction: column; align-items: center; margin-bottom: 32px; }
.step-connector { display: flex; flex-direction: column; align-items: center; padding: 4px 0; }
.connector-line { width: 2px; height: 16px; background: #2e3142; }
.connector-arrow { color: #4a5568; font-size: 16px; }
.step-node { display: flex; flex-direction: column; align-items: center; width: 100%; max-width: 600px; }

.step-card {
	width: 100%; background: #1a1d27; border: 1px solid #2e3142;
	border-radius: 12px; padding: 16px;
	&.step-email { border-left: 3px solid #6c5ce7; }
	&.step-wait { border-left: 3px solid #ffaa00; }
	&.step-condition { border-left: 3px solid #00d68f; }
}
.step-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.step-badge {
	width: 32px; height: 32px; border-radius: 8px;
	display: flex; align-items: center; justify-content: center;
	background: rgba(108,92,231,0.12); color: #a78bfa; font-size: 16px;
}
.step-meta { display: flex; align-items: center; gap: 8px; }
.step-order { font-size: 14px; font-weight: 600; color: #e2e8f0; }
.step-subject { font-size: 13px; color: #a1a5b3; margin-bottom: 12px; font-style: italic; }
.step-stats { display: flex; gap: 20px; }
.st { display: flex; flex-direction: column; gap: 2px; }
.st-v { font-size: 16px; font-weight: 700; color: #e2e8f0; }
.st-l { font-size: 11px; color: #6b7280; text-transform: uppercase; }
.step-wait-info {
	display: flex; align-items: center; gap: 8px;
	font-size: 14px; color: #ffaa00; font-weight: 500;
}
.step-condition-info { font-size: 13px; color: #a1a5b3; }

.detail-tabs { margin-top: 8px; }
.tab-toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }

.scoring-cards { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 20px; }
.scoring-card {
	padding: 16px; border-radius: 10px; text-align: center;
	background: #1a1d27; border: 1px solid #2e3142;
	&.hot .sc-val { color: #ef4444; }
	&.warm .sc-val { color: #f59e0b; }
	&.cold .sc-val { color: #3b82f6; }
	&.avg .sc-val { color: #a78bfa; }
}
.sc-val { font-size: 28px; font-weight: 700; }
.sc-label { font-size: 13px; color: #8892a8; margin-top: 4px; }
.sc-threshold { font-size: 11px; color: #4a5568; margin-top: 2px; }

:deep(.n-tabs .n-tabs-tab) { color: #8892a8; }
:deep(.n-tabs .n-tabs-tab--active) { color: #a78bfa; }
:deep(.n-data-table) { --n-td-color: transparent; --n-th-color: transparent; --n-border-color: #2e3142; }
</style>
