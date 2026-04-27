<template>
  <div class="p-4">
    <n-page-header @back="router.back()" :title="detail?.name || ''" />
    <n-tabs v-if="detail" type="line" class="mt-4">
      <n-tab-pane :name="'steps'" :tab="t('sequences.detail.steps')">
        <n-space vertical>
          <n-card v-for="step in detail.steps" :key="step.id" size="small" class="ml-6">
            <template #header>
              <n-space align="center">
                <n-tag :type="stepTypeColor(step.step_type)" size="small">{{ step.step_type }}</n-tag>
                <span>{{ t('sequences.form.step') }} {{ step.step_order }}</span>
                <span v-if="step.step_type === 'email'" class="text-gray-500">- {{ step.subject }}</span>
              </n-space>
            </template>
            <template v-if="step.step_type === 'email'">
              <n-grid :cols="5" :x-gap="16">
                <n-gi>
                  <n-statistic :label="t('sequences.stats.sent')" :value="step.sent_count" />
                </n-gi>
                <n-gi>
                  <n-statistic :label="t('sequences.stats.opened')" :value="step.opened_count" />
                </n-gi>
                <n-gi>
                  <n-statistic :label="t('sequences.stats.clicked')" :value="step.clicked_count" />
                </n-gi>
                <n-gi>
                  <n-statistic :label="t('sequences.stats.bounced')" :value="step.bounced_count" />
                </n-gi>
                <n-gi>
                  <n-statistic label="Replied" :value="step.replied_count" />
                </n-gi>
              </n-grid>
            </template>
            <template v-if="step.step_type === 'wait'">
              {{ step.wait_days }} {{ t('sequences.form.days') }} {{ step.wait_hours }} {{ t('sequences.form.hours') }}
            </template>
            <template v-if="step.step_type === 'condition'">
              {{ step.condition_type }}
              <template v-if="step.condition_step_id"> ({{ t('sequences.form.step') }} {{ step.condition_step_id }})</template>
            </template>
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane :name="'enrollments'" :tab="t('sequences.detail.enrollments')">
        <bt-table-layout>
          <template #toolsRight>
            <n-select v-model:value="enrollmentStatus" :options="enrollmentStatusOptions" class="w-140px" @update:value="fetchEnrollments" />
            <n-button type="primary" @click="handleEnroll">
              {{ t('sequences.enrollContacts') }}
            </n-button>
          </template>
          <template #table>
            <n-data-table :columns="enrollmentColumns" :data="enrollments" :loading="enrollmentsLoading" :pagination="enrollmentPagination" />
          </template>
        </bt-table-layout>
      </n-tab-pane>
      <n-tab-pane :name="'lead-scoring'" :tab="'Lead Scoring'">
        <n-space vertical size="large">
          <n-grid :cols="4" :x-gap="16">
            <n-gi>
              <n-card size="small">
                <n-statistic label="Hot Leads" :value="hotCount">
                  <template #suffix>
                    <n-tag type="error" size="small" style="margin-left: 8px">&gt; 40</n-tag>
                  </template>
                </n-statistic>
              </n-card>
            </n-gi>
            <n-gi>
              <n-card size="small">
                <n-statistic label="Warm Leads" :value="warmCount">
                  <template #suffix>
                    <n-tag type="warning" size="small" style="margin-left: 8px">10-40</n-tag>
                  </template>
                </n-statistic>
              </n-card>
            </n-gi>
            <n-gi>
              <n-card size="small">
                <n-statistic label="Cold Leads" :value="coldCount">
                  <template #suffix>
                    <n-tag type="info" size="small" style="margin-left: 8px">&lt; 10</n-tag>
                  </template>
                </n-statistic>
              </n-card>
            </n-gi>
            <n-gi>
              <n-card size="small">
                <n-statistic label="Avg Score" :value="avgScore" :precision="1" />
              </n-card>
            </n-gi>
          </n-grid>
          <n-data-table :columns="scoringColumns" :data="scoredEnrollments" :pagination="{ pageSize: 10 }" :default-sort="{ columnKey: 'score', order: 'descend' }" />
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NTag, NPopconfirm, NButton, NProgress } from 'naive-ui'
import { Message } from '@/utils'
import { getSequenceDetail, getEnrollments, enrollContacts, removeEnrollment } from '@/api/modules/sequences/sequence'
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
  enrollments.value
    .map((e) => {
      const score = calculateLeadScore(e)
      return { ...e, score, tier: getScoreTier(score) }
    })
    .sort((a, b) => b.score - a.score)
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
  {
    key: 'email',
    title: 'Email',
    minWidth: 220,
    ellipsis: { tooltip: true },
  },
  {
    key: 'score',
    title: 'Score',
    width: 90,
    sorter: 'default' as const,
    defaultSortOrder: 'descend' as const,
    render: (row: ScoredEnrollment) => {
      const maxScore = 100
      const pct = Math.min((row.score / maxScore) * 100, 100)
      const status = row.score >= 40 ? 'error' : row.score >= 10 ? 'warning' : 'success'
      return h(NProgress, {
        type: 'line',
        percentage: pct,
        status,
        indicatorPlacement: 'inside',
        showIndicator: true,
        style: { width: '80px' },
      })
    },
  },
  {
    key: 'tier',
    title: 'Tier',
    width: 100,
    render: (row: ScoredEnrollment) => {
      const cfg = tierConfig[row.tier]
      return h(NTag, { size: 'small', type: cfg.type }, { default: () => cfg.label })
    },
  },
  { key: 'total_opens', title: 'Opens', width: 80 },
  { key: 'total_clicks', title: 'Clicks', width: 80 },
  { key: 'total_replies', title: 'Replies', width: 80 },
  { key: 'total_emails_sent', title: 'Emails', width: 80 },
  {
    key: 'status',
    title: 'Status',
    width: 100,
    render: (row: ScoredEnrollment) => {
      const s = enrollmentStatusMap[row.status]
      return s ? h(NTag, { size: 'small', type: s.type }, { default: () => t(s.label) }) : row.status
    },
  },
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
    key: 'status',
    title: t('sequences.enrollments.status'),
    width: 120,
    render: (row: Enrollment) => {
      const s = enrollmentStatusMap[row.status]
      return s ? h(NTag, { size: 'small', type: s.type }, { default: () => t(s.label) }) : row.status
    },
  },
  { key: 'total_emails_sent', title: t('sequences.enrollments.emailsSent'), width: 120 },
  { key: 'total_opens', title: t('sequences.enrollments.opens'), width: 90 },
  { key: 'total_clicks', title: t('sequences.enrollments.clicks'), width: 90 },
  {
    key: 'enrolled_at',
    title: t('sequences.enrollments.enrolledAt'),
    width: 170,
    render: (row: Enrollment) => formatTime(row.enrolled_at),
  },
  {
    title: t('common.columns.actions'),
    key: 'actions',
    width: 100,
    render: (row: Enrollment) => {
      return h(NPopconfirm, {
        onPositiveClick: () => handleRemoveEnrollment(row.id),
      }, {
        trigger: () => h(NButton, { type: 'error', text: true, size: 'small' }, { default: () => t('common.actions.delete') }),
        default: () => t('sequences.enrollments.removeConfirm'),
      })
    },
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
