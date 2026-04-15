<template>
  <div>
    <bt-table-layout>
      <template #toolsLeft>
        <span>{{ t('sequences.status') }}</span>
        <div class="w-140px">
          <n-select v-model:value="tableParams.status" :options="statusOptions" @update:value="resetTable" />
        </div>
      </template>
      <template #toolsRight>
        <bt-search v-model:value="tableParams.keyword" :placeholder="t('sequences.search.placeholder')" @search="resetTable" />
        <n-button type="primary" @click="handleCreate">
          {{ t('common.actions.create') }}
        </n-button>
      </template>
      <template #table>
        <n-data-table v-bind="tableProps" :columns="columns" :row-props="rowProps">
          <template #empty>
            <bt-table-help />
          </template>
        </n-data-table>
      </template>
      <template #pageRight>
        <bt-table-page v-bind="pageProps" @refresh="fetchTable" />
      </template>
    </bt-table-layout>
  </div>
</template>

<script lang="tsx" setup>
import { DataTableColumns, NButton, NFlex, NTag, NPopconfirm } from 'naive-ui'
import { Message } from '@/utils'
import { useDataTable } from '@/hooks/useDataTable'
import { getSequenceList, deleteSequence, activateSequence, pauseSequence, resumeSequence } from '@/api/modules/sequences/sequence'
import { formatTime } from '@/utils'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import type { Sequence, SequenceParams } from './interface'

const { t } = useI18n()
const router = useRouter()

const statusOptions = [
  { label: t('common.all.text'), value: -1 },
  { label: t('sequences.statusDraft'), value: 0 },
  { label: t('sequences.statusActive'), value: 1 },
  { label: t('sequences.statusPaused'), value: 2 },
  { label: t('sequences.statusArchived'), value: 3 },
]

const statusMap: Record<number, { label: string; type: 'default' | 'success' | 'warning' | 'error' }> = {
  0: { label: 'sequences.statusDraft', type: 'default' },
  1: { label: 'sequences.statusActive', type: 'success' },
  2: { label: 'sequences.statusPaused', type: 'warning' },
  3: { label: 'sequences.statusArchived', type: 'error' },
}

const { tableParams, tableProps, pageProps, fetchTable, resetTable } = useDataTable<Sequence, SequenceParams>({
  params: {
    page: 1,
    page_size: 10,
    keyword: '',
    status: -1,
  },
  useParams: (params) => params,
  fetchFn: getSequenceList as any,
})

const columns = computed<DataTableColumns<Sequence>>(() => [
  {
    key: 'name',
    title: t('sequences.columns.name'),
    minWidth: 180,
    ellipsis: { tooltip: true },
  },
  {
    key: 'status',
    title: t('sequences.columns.status'),
    width: 120,
    render: (row) => {
      const s = statusMap[row.status]
      return s ? <NTag size="small" type={s.type}>{t(s.label)}</NTag> : row.status
    },
  },
  {
    key: 'step_count',
    title: t('sequences.columns.steps'),
    width: 100,
  },
  {
    key: 'total_enrolled',
    title: t('sequences.columns.enrolled'),
    width: 110,
  },
  {
    key: 'total_completed',
    title: t('sequences.columns.completed'),
    width: 110,
  },
  {
    key: 'create_time',
    title: t('sequences.columns.created'),
    width: 170,
    render: (row) => formatTime(row.create_time),
  },
  {
    title: t('common.columns.actions'),
    key: 'actions',
    width: 280,
    render: (row) => (
      <NFlex inline={true}>
        <NButton type="primary" text size="small" onClick={() => router.push(`/sequences/${row.id}`)}>
          {t('common.actions.view')}
        </NButton>
        {row.status === 0 && (
          <NButton type="success" text size="small" onClick={() => handleActivate(row.id)}>
            {t('sequences.activate')}
          </NButton>
        )}
        {row.status === 1 && (
          <NButton type="warning" text size="small" onClick={() => handlePause(row.id)}>
            {t('sequences.pause')}
          </NButton>
        )}
        {row.status === 2 && (
          <NButton type="success" text size="small" onClick={() => handleResume(row.id)}>
            {t('sequences.resume')}
          </NButton>
        )}
        <NButton type="primary" text size="small" onClick={() => router.push(`/sequences/${row.id}/edit`)}>
          {t('common.actions.edit')}
        </NButton>
        <NPopconfirm onPositiveClick={() => handleDelete(row)}>
          {{
            trigger: () => <NButton type="error" text size="small">{t('common.actions.delete')}</NButton>,
            default: () => t('sequences.delete.confirm'),
          }}
        </NPopconfirm>
      </NFlex>
    ),
  },
])

const rowProps = (row: Sequence) => ({
  style: 'cursor: pointer',
  onClick: () => router.push(`/sequences/${row.id}`),
})

function handleCreate() {
  router.push('/sequences/create')
}

async function handleActivate(id: number) {
  await activateSequence({ id })
  fetchTable()
}

async function handlePause(id: number) {
  await pauseSequence({ id })
  fetchTable()
}

async function handleResume(id: number) {
  await resumeSequence({ id })
  fetchTable()
}

async function handleDelete(row: Sequence) {
  await deleteSequence({ id: row.id })
  Message.success(t('sequences.delete.success'))
  fetchTable()
}
</script>
