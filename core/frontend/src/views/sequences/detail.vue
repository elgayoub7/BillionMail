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
              <n-grid :cols="4" :x-gap="16">
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
    </n-tabs>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NTag, NPopconfirm, NButton } from 'naive-ui'
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
    enrollments.value = res.data.list
    enrollmentPagination.value.itemCount = res.data.total
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
  detail.value = res.data
}

onMounted(async () => {
  await loadDetail()
  fetchEnrollments()
})
</script>
