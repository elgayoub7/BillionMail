<template>
	<div class="domains-page">
		<div class="page-header">
			<h1 class="page-title">{{ t('layout.menu.domain') }}</h1>
			<div class="header-actions">
				<bt-help href="https://www.billionmail.com/start/domain.html" :text="t('domain.help')" />
				<n-button type="primary" @click="handleAddDomain">
					<i class="i-mdi:plus"></i> {{ t('domain.addDomain') }}
				</n-button>
			</div>
		</div>

		<!-- Domain Cards -->
		<div v-if="tableList.length" class="cards-grid">
			<div v-for="row in tableList" :key="row.domain" class="domain-card">
				<div class="card-top">
					<div class="card-info">
						<h3 class="card-domain">
							<i v-if="row.hasbrandinfo == 1" class="i-domain:brand-info" style="font-size:14px;margin-right:4px"></i>
							{{ row.domain }}
						</h3>
						<n-tag v-if="row.default === 1" size="small" :bordered="false" type="info">Default</n-tag>
					</div>
				</div>

				<div class="card-stats">
					<div class="stat">
						<i class="i-mdi:email-outline stat-icon"></i>
						<span class="stat-val">{{ row.mailboxes || 0 }}</span>
						<span class="stat-lbl">Mailboxes</span>
					</div>
					<div class="stat">
						<i class="i-mdi:harddisk stat-icon"></i>
						<span class="stat-val">{{ getByteUnit(row.quota) }}</span>
						<span class="stat-lbl">Quota</span>
					</div>
					<div class="stat">
						<i class="i-mdi:chart-bar stat-icon"></i>
						<span class="stat-val">{{ getByteUnit(row.current_usage) }}</span>
						<span class="stat-lbl">Usage</span>
					</div>
				</div>

				<!-- SSL & Spam -->
				<div class="card-badges">
					<div class="badge-item" @click="handleShowSsl(row)">
						<template v-if="row.cert_info && row.cert_info.endtime">
							<span :class="getSslClass(row.cert_info.endtime)">
								<i class="i-mdi:shield-check-outline"></i>
								{{ getSslLabel(row.cert_info.endtime) }}
							</span>
						</template>
						<template v-else>
							<span class="badge-warn">
								<i class="i-mdi:shield-alert-outline"></i> No SSL
							</span>
						</template>
					</div>
					<div class="badge-item" @click="handleBlacklistCheck(row)">
						<span :class="row.black_check_result === 1 ? 'badge-ok' : 'badge-warn'">
							<i class="i-mdi:shield-outline"></i>
							{{ row.black_check_result === 1 ? 'Clean' : 'Check Spam' }}
						</span>
					</div>
					<div v-if="row.multi_ip_domains" class="badge-item" @click="handleMultiIpDomains(row)">
						<span class="badge-neutral">
							<i class="i-mdi:ip-network-outline"></i>
							{{ row.multi_ip_domains.outbound_ip }}
						</span>
					</div>
				</div>

				<!-- Actions -->
				<div class="card-actions">
					<n-button quaternary size="small" @click="handleDNSRecord(row)">
						<i class="i-mdi:dns-outline"></i> DNS
					</n-button>
					<n-button quaternary size="small" @click="handleEdit(row)">
						<i class="i-mdi:pencil-outline"></i> Edit
					</n-button>
					<n-button quaternary size="small" :disabled="row.default === 1" @click="handleSetDefault(row)">
						<i class="i-mdi:star-outline"></i> Default
					</n-button>
					<n-popconfirm @positive-click="handleDeleteConfirm(row)">
						<template #trigger>
							<n-button quaternary size="small" type="error"><i class="i-mdi:delete-outline"></i></n-button>
						</template>
						{{ t('domain.delete.confirm', { domain: row.domain }) }}
					</n-popconfirm>
				</div>
			</div>
		</div>

		<div v-else-if="!loading" class="empty-state">
			<i class="i-mdi:dns-outline empty-icon"></i>
			<h3>No domains</h3>
			<p>Add your first domain to start sending emails</p>
		</div>

		<div v-if="tableTotal > 0" class="pagination-bar">
			<bt-table-page
				v-model:page="tableParams.page"
				v-model:page-size="tableParams.page_size"
				:item-count="tableTotal"
				@refresh="getTableData"
			/>
		</div>

		<form-modal />
		<ssl-modal />
		<dns-modal />
		<domain-ip-set-modal />
		<blacklist-modal />
		<check-logs-modal />
	</div>
</template>

<script lang="ts" setup>
import { confirm, getByteUnit } from '@/utils'
import { useModal } from '@/hooks/modal/useModal'
import { useTableData } from '@/hooks/useTableData'
import { deleteDomain, getDomainList, setDefaultDomain } from '@/api/modules/domain'
import type { MailDomain, MailDomainParams } from './interface'

import DomainForm from './components/DomainForm.vue'
import DomainSsl from './components/DomainSsl/index.vue'
import DomainDns from './components/DomainDns.vue'
import DomainIpSet from './components/DomainIpSet.vue'
import BlacklistDetection from './components/BlacklistDetection.vue'
import CheckLogs from './components/CheckLogs.vue'

const { t } = useI18n()

const { tableParams, tableList, loading, tableTotal, getTableData } = useTableData<MailDomain, MailDomainParams>({
	immediate: true,
	params: { page: 1, page_size: 10, keyword: '' },
	fetchFn: getDomainList,
})

const route = useRoute()
const router = useRouter()

function getSslClass(endtime: number) {
	const days = Math.floor((endtime * 1000 - Date.now()) / 86400000)
	return days < 0 ? 'badge-error' : days < 30 ? 'badge-warn' : 'badge-ok'
}

function getSslLabel(endtime: number) {
	const days = Math.floor((endtime * 1000 - Date.now()) / 86400000)
	return days < 0 ? 'Expired' : `${days}d left`
}

function handleBlacklistCheck(row: MailDomain) {
	blacklistModalApi.setState({ row })
	blacklistModalApi.open()
}

const [FormModal, formModalApi] = useModal({ component: DomainForm, state: { isEdit: false, refresh: getTableData } })

const handleAddDomain = () => {
	formModalApi.setState({ row: null, isEdit: false })
	formModalApi.open()
}

const [SslModal, sslModalApi] = useModal({ component: DomainSsl, state: { refresh: getTableData } })
const handleShowSsl = (row: MailDomain) => { sslModalApi.setState({ row }); sslModalApi.open() }

const [DnsModal, dnsModalApi] = useModal({ component: DomainDns, state: { refresh: getTableData } })
const handleDNSRecord = (row: MailDomain) => { dnsModalApi.setState({ row }); dnsModalApi.open() }

const handleSetDefault = (row: MailDomain) => {
	confirm({
		title: t('domain.setDefault.title', { domain: row.domain }),
		content: t('domain.setDefault.confirm'),
		onConfirm: async () => { await setDefaultDomain({ domain: row.domain }); getTableData() },
	})
}

const handleEdit = (row: MailDomain) => {
	router.push({ name: 'EditDomain', params: { domain: row.domain } })
}

const [DomainIpSetModal, domainIpSetModalApi] = useModal({ component: DomainIpSet })
const handleMultiIpDomains = (row: MailDomain) => { domainIpSetModalApi.setState({ row }); domainIpSetModalApi.open() }

const handleDeleteConfirm = (row: MailDomain) => {
	confirm({
		title: t('domain.delete.title'),
		content: t('domain.delete.confirm', { domain: row.domain }),
		confirmText: t('common.actions.delete'),
		confirmType: 'error',
		onConfirm: async () => { await deleteDomain({ domain: row.domain }); getTableData() },
	})
}

const [BlacklistModal, blacklistModalApi] = useModal({ component: BlacklistDetection })
const [CheckLogsModal, checkLogsModalApi] = useModal({ component: CheckLogs })

onMounted(() => {
	if (route.query.init === 'init-domain') handleAddDomain()
})
</script>

<style lang="scss" scoped>
.domains-page { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { margin: 0; font-size: 22px; font-weight: 700; color: #e2e8f0; }
.header-actions { display: flex; gap: 10px; align-items: center; }

.cards-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(360px, 1fr)); gap: 16px; }
.domain-card {
	background: #1a1d27; border: 1px solid #2e3142; border-radius: 12px;
	padding: 20px; transition: all 0.2s ease;
	&:hover { border-color: #6c5ce7; }
}
.card-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
.card-info { display: flex; align-items: center; gap: 8px; }
.card-domain { margin: 0; font-size: 16px; font-weight: 600; color: #e2e8f0; display: flex; align-items: center; }

.card-stats { display: flex; gap: 20px; margin-bottom: 14px; }
.stat { display: flex; align-items: center; gap: 6px; }
.stat-icon { font-size: 16px; color: #6b7280; }
.stat-val { font-size: 14px; font-weight: 600; color: #e2e8f0; }
.stat-lbl { font-size: 12px; color: #6b7280; }

.card-badges { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 14px; }
.badge-item {
	font-size: 12px; cursor: pointer; display: flex; align-items: center; gap: 4px;
	padding: 4px 8px; border-radius: 6px; background: rgba(255,255,255,0.04);
}
.badge-ok { color: #22c55e; i { color: #22c55e; } }
.badge-warn { color: #f59e0b; i { color: #f59e0b; } }
.badge-error { color: #ef4444; i { color: #ef4444; } }
.badge-neutral { color: #8892a8; i { color: #8892a8; } }

.card-actions { display: flex; gap: 4px; border-top: 1px solid #2e3142; padding-top: 12px; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 80px 20px; color: #6b7280;
	h3 { font-size: 18px; color: #8892a8; margin: 16px 0 8px; }
	p { font-size: 14px; margin: 0; }
}
.empty-icon { font-size: 48px; color: #2e3142; }
.pagination-bar { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
