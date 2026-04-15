<template>
  <div class="p-20px">
    <div class="flex justify-between items-center mb-20px">
      <n-breadcrumb>
        <n-breadcrumb-item>
          <router-link to="/market/task">{{ $t('market.task.title') }}</router-link>
        </n-breadcrumb-item>
        <n-breadcrumb-item>
          <router-link :to="`/market/task/analytics/${taskId}`">{{ $t('market.task.actions.analytics') }}</router-link>
        </n-breadcrumb-item>
        <n-breadcrumb-item>{{ $t('market.task.recipients.title') }}</n-breadcrumb-item>
      </n-breadcrumb>
    </div>

    <bt-table-layout>
      <template #toolsLeft>
        <span>{{ $t('market.task.recipients.status') }}</span>
        <div class="w-140px">
          <n-select v-model:value="tableParams.status" :options="statusOptions" @update:value="resetTable" />
        </div>
      </template>
      <template #toolsRight>
        <bt-search v-model:value="tableParams.search" :placeholder="$t('market.task.recipients.search')" @search="resetTable" />
      </template>
      <template #table>
        <n-data-table v-bind="tableProps" :columns="columns">
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
import { DataTableColumns, NTag } from 'naive-ui'
import { useDataTable } from '@/hooks/useDataTable'
import { getTaskRecipients } from '@/api/modules/market/task'
import { formatTime, getNumber } from '@/utils'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()

const taskId = computed(() => getNumber(route.params.id || '0'))

interface RecipientParams {
  task_id: number
  status: string
  search: string
  page: number
  page_size: number
}

const statusOptions = computed(() => [
  { label: t('common.all.text'), value: '' },
  { label: t('market.task.recipients.filterSent'), value: 'sent' },
  { label: t('market.task.recipients.filterBounced'), value: 'bounced' },
  { label: t('market.task.recipients.filterOpened'), value: 'opened' },
  { label: t('market.task.recipients.filterClicked'), value: 'clicked' },
  { label: t('market.task.recipients.filterUnopened'), value: 'unopened' },
])

const statusTypeMap: Record<string, 'success' | 'error' | 'warning' | 'default' | 'info'> = {
  sent: 'success',
  opened: 'info',
  clicked: 'info',
  bounced: 'error',
  unopened: 'warning',
  pending: 'warning',
  deferred: 'warning',
}

const { tableParams, tableProps, pageProps, fetchTable, resetTable } = useDataTable<any, RecipientParams>({
  params: {
    task_id: taskId.value,
    status: '',
    search: '',
    page: 1,
    page_size: 20,
  },
  useParams: (params) => params,
  fetchFn: getTaskRecipients as any,
})

const columns = computed<DataTableColumns>(() => [
  {
    key: 'recipient',
    title: t('market.task.recipients.email'),
    minWidth: 200,
    ellipsis: { tooltip: true },
  },
  {
    key: 'status',
    title: t('market.task.recipients.status'),
    width: 110,
    render: (row) => {
      const type = statusTypeMap[row.status] || 'default'
      return h(NTag, { type, size: 'small' }, { default: () => row.status })
    },
  },
  {
    key: 'sent_time',
    title: t('market.task.recipients.sentTime'),
    width: 170,
    render: (row) => (row.sent_time ? formatTime(row.sent_time * 1000) : '-'),
  },
  {
    key: 'opened',
    title: t('market.task.recipients.opens'),
    width: 80,
    render: (row) => row.opened || 0,
  },
  {
    key: 'first_open',
    title: t('market.task.recipients.firstOpen'),
    width: 170,
    render: (row) => (row.first_open ? formatTime(row.first_open * 1000) : '-'),
  },
  {
    key: 'clicked',
    title: t('market.task.recipients.clicks'),
    width: 80,
    render: (row) => row.clicked || 0,
  },
  {
    key: 'first_click',
    title: t('market.task.recipients.firstClick'),
    width: 170,
    render: (row) => (row.first_click ? formatTime(row.first_click * 1000) : '-'),
  },
])
</script>
