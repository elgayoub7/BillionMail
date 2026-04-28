<template>
	<div class="mailbox-page">
		<div class="page-header">
			<h1 class="page-title">Mailboxes</h1>
			<div class="header-actions">
				<bt-search
					v-model:value="tableParams.keyword"
					placeholder="Search mailboxes..."
					@search="resetTable" />
				<domain-select v-model:value="tableParams.domain" @update:value="resetTable" />
				<n-button type="primary" @click="handleAdd">
					<template #icon><i class="i-mdi:plus"></i></template>
					Add
				</n-button>
				<n-button @click="handleBatchAdd">Batch</n-button>
				<n-button @click="handleImport">Import</n-button>
			</div>
		</div>

		<!-- Mailbox Cards -->
		<div v-if="tableProps.data?.length" class="cards-grid">
			<div v-for="row in tableProps.data" :key="row.username" class="mailbox-card">
				<div class="card-top">
					<div class="card-info">
						<h3 class="card-email">{{ row.username }}</h3>
						<span class="card-domain">{{ row.domain }}</span>
					</div>
					<n-switch
						:value="row.active"
						:checked-value="1"
						:unchecked-value="0"
						size="small"
						@update:value="val => handleStatusChange(row, val)" />
				</div>

				<!-- Status -->
				<div class="health-section">
					<div class="health-header">
						<span class="health-label">Status</span>
						<span class="health-value" :class="row.active ? 'health-good' : 'health-bad'">{{ row.active ? 'Active' : 'Inactive' }}</span>
					</div>
				</div>

				<!-- Meta Info -->
				<div class="card-meta">
					<span v-if="row.full_name" class="meta-item">
						<i class="i-mdi:account-outline"></i>
						{{ row.full_name }}
					</span>
					<span class="meta-item">
						<i class="i-custom:smtp"></i>
						{{ row.is_admin ? 'Admin' : 'General' }}
					</span>
				</div>

				<!-- Actions -->
				<div class="card-actions">
					<n-button quaternary size="small" @click="handleCopyLogin(row)">
						<i class="i-mdi:content-copy"></i>
						Login Info
					</n-button>
					<n-button quaternary size="small" @click="handleEdit(row)">
						<i class="i-mdi:pencil-outline"></i>
						Edit
					</n-button>
					<n-popconfirm @positive-click="handleDelete(row)">
						<template #trigger>
							<n-button quaternary size="small" type="error">
								<i class="i-mdi:delete-outline"></i>
							</n-button>
						</template>
						Delete {{ row.username }}?
					</n-popconfirm>
				</div>
			</div>
		</div>

		<!-- Empty -->
		<div v-else class="empty-state">
			<i class="i-mdi:email-fast-outline empty-icon"></i>
			<h3>No mailboxes</h3>
			<p>Add your first mailbox to get started</p>
		</div>

		<!-- Pagination -->
		<div v-if="pageProps.total > 0" class="pagination-bar">
			<bt-table-page v-bind="pageProps" @refresh="fetchTable" />
		</div>

		<!-- Modals -->
		<form-modal />
		<batch-add-modal ref="batchAddRef" @refresh="fetchTable" />
		<import-modal ref="importRef" @refresh="fetchTable" />
	</div>
</template>

<script lang="ts" setup>
import { useBrowserLocation } from '@vueuse/core'
import { confirm } from '@/utils'
import { useModal } from '@/hooks/modal/useModal'
import { useCopy } from '@/hooks/useCopy'
import { useDataTable } from '@/hooks/useDataTable'
import { deleteMailbox, getMailboxList, updateMailbox } from '@/api/modules/mailbox'
import type { MailBox, MailBoxParams } from './interface'

import DomainSelect from './components/DomainSelect.vue'
import MailboxForm from './components/MailboxForm.vue'
import BatchAddModal from './components/MailboxBatchAdd.vue'
import ImportModal from './components/MailboxImport.vue'

const { t } = useI18n()
const location = useBrowserLocation()
const { copyText } = useCopy()
const batchAddRef = useTemplateRef('batchAddRef')
const importRef = useTemplateRef('importRef')

const getHealthClass = (score: number) => {
	if (score >= 80) return 'health-good'
	if (score >= 50) return 'health-warn'
	return 'health-bad'
}

const getHealthLabel = (score: number) => `${score}/100`

const { tableParams, tableProps, pageProps, fetchTable, resetTable } = useDataTable<MailBox, MailBoxParams>({
	loading: true,
	immediate: true,
	params: {
		page: 1,
		page_size: 12,
		domain: location.value?.state?.domain || '',
		keyword: '',
	},
	rowKey: row => row.username,
	fetchFn: getMailboxList,
})

const [FormModal, formModalApi] = useModal({
	component: MailboxForm,
	state: { isEdit: false, refresh: fetchTable },
})

const handleAdd = () => {
	formModalApi.setState({ isEdit: false, row: null })
	formModalApi.open()
}

const handleBatchAdd = () => {
	batchAddRef.value?.open()
}

const handleImport = () => {
	importRef.value?.open()
}

const handleCopyLogin = (row: MailBox) => {
	copyText(
		t('mailbox.loginInfo.template', {
			webmail: window.location.origin + '/roundcube',
			username: row.username,
			password: row.password,
			mx: row.mx,
		})
	)
}

const handleStatusChange = async (row: MailBox, val: number) => {
	await updateMailbox({
		full_name: row.full_name,
		domain: row.domain,
		password: row.password,
		quota: row.quota,
		isAdmin: row.is_admin,
		active: val,
	})
	row.active = val
}

const handleEdit = (row: MailBox) => {
	formModalApi.setState({ isEdit: true, row })
	formModalApi.open()
}

const handleDelete = (row: MailBox) => {
	confirm({
		title: t('mailbox.delete.title'),
		content: t('mailbox.delete.confirm', { name: row.username }),
		confirmText: t('common.actions.delete'),
		confirmType: 'error',
		onConfirm: async () => {
			await deleteMailbox({ emails: [row.username] })
			fetchTable()
		},
	})
}
</script>

<style lang="scss" scoped>
.mailbox-page {
	padding: 24px;
}

.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24px;
}

.page-title {
	margin: 0;
	font-size: 22px;
	font-weight: 700;
	color: #e2e8f0;
}

.header-actions {
	display: flex;
	gap: 10px;
	align-items: center;
}

.cards-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
	gap: 16px;
}

.mailbox-card {
	background: #1a1d27;
	border: 1px solid #2e3142;
	border-radius: 12px;
	padding: 20px;
	transition: all 0.2s ease;

	&:hover {
		border-color: #6c5ce7;
	}
}

.card-top {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 16px;
}

.card-info {
	min-width: 0;
	flex: 1;
}

.card-email {
	margin: 0;
	font-size: 16px;
	font-weight: 600;
	color: #e2e8f0;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.card-domain {
	font-size: 12px;
	color: #6b7280;
}

.health-section,
.daily-section {
	margin-bottom: 12px;
}

.health-header,
.daily-header {
	display: flex;
	justify-content: space-between;
	margin-bottom: 6px;
}

.health-label,
.daily-label {
	font-size: 12px;
	color: #8892a8;
}

.health-value {
	font-size: 13px;
	font-weight: 600;

	&.health-good { color: #22c55e; }
	&.health-warn { color: #f59e0b; }
	&.health-bad { color: #ef4444; }
}

.daily-value {
	font-size: 13px;
	font-weight: 600;
	color: #e2e8f0;
}

.health-bar,
.daily-bar {
	height: 6px;
	background: #2e3142;
	border-radius: 3px;
	overflow: hidden;
}

.health-fill {
	height: 100%;
	background: linear-gradient(90deg, #22c55e, #4ade80);
	border-radius: 3px;
	transition: width 0.3s ease;
}

.daily-fill {
	height: 100%;
	background: linear-gradient(90deg, #6c5ce7, #a78bfa);
	border-radius: 3px;
	transition: width 0.3s ease;
}

.card-meta {
	display: flex;
	gap: 16px;
	margin-bottom: 12px;
}

.meta-item {
	display: flex;
	align-items: center;
	gap: 6px;
	font-size: 12px;
	color: #6b7280;

	i {
		font-size: 14px;
	}
}

.card-actions {
	display: flex;
	gap: 4px;
	border-top: 1px solid #2e3142;
	padding-top: 12px;
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 80px 20px;
	color: #6b7280;

	h3 {
		font-size: 18px;
		color: #8892a8;
		margin: 16px 0 8px;
	}

	p {
		font-size: 14px;
		margin: 0;
	}
}

.empty-icon {
	font-size: 48px;
	color: #2e3142;
}

.pagination-bar {
	display: flex;
	justify-content: flex-end;
	margin-top: 24px;
}
</style>
