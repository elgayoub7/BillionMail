<template>
	<div class="domain-health-page">
		<div class="bt-title">Domain Health</div>

		<!-- Actions -->
		<div class="actions-row">
			<n-space>
				<n-button type="primary" @click="handleCheckAll" :loading="checkingAll">
					Check All Domains
				</n-button>
				<n-button @click="fetchDomains" :loading="loading">
					Refresh
				</n-button>
			</n-space>
		</div>

		<!-- Health Summary -->
		<div class="summary-row">
			<n-card class="summary-card">
				<div class="summary-value" style="color: #18a058">{{ healthyCount }}</div>
				<div class="summary-label">Healthy</div>
			</n-card>
			<n-card class="summary-card">
				<div class="summary-value" style="color: #e6a23c">{{ warningCount }}</div>
				<div class="summary-label">Warnings</div>
			</n-card>
			<n-card class="summary-card">
				<div class="summary-value" style="color: #d03050">{{ errorCount }}</div>
				<div class="summary-label">Issues</div>
			</n-card>
		</div>

		<!-- Domains Table -->
		<n-card>
			<n-data-table
				:columns="columns"
				:data="domains"
				:loading="loading"
				:bordered="false"
			/>
		</n-card>
	</div>
</template>

<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { NTag, NButton, useMessage } from 'naive-ui'
import { listDomainHealth, checkAllDomains, checkDomain } from '@/api/modules/domain-health'

interface DomainHealth {
	domain: string
	spf_status: string
	dkim_status: string
	dmarc_status: string
	mx_status: string
	issues: string[]
	last_checked: number
	spf_record: string
	dkim_record: string
	dmarc_record: string
	mx_records: string
}

const message = useMessage()
const domains = ref<DomainHealth[]>([])
const loading = ref(false)
const checkingAll = ref(false)

const healthyCount = ref(0)
const warningCount = ref(0)
const errorCount = ref(0)

function statusTag(status: string) {
	if (status === 'pass') return h(NTag, { size: 'small', type: 'success', bordered: false }, () => '✓ Pass')
	if (status === 'missing') return h(NTag, { size: 'small', type: 'error', bordered: false }, () => '✗ Missing')
	return h(NTag, { size: 'small', type: 'warning', bordered: false }, () => status)
}

const columns = [
	{
		title: 'Domain',
		key: 'domain',
		ellipsis: true,
		width: 200,
	},
	{
		title: 'SPF',
		key: 'spf_status',
		width: 100,
		render: (row: DomainHealth) => statusTag(row.spf_status),
	},
	{
		title: 'DKIM',
		key: 'dkim_status',
		width: 100,
		render: (row: DomainHealth) => statusTag(row.dkim_status),
	},
	{
		title: 'DMARC',
		key: 'dmarc_status',
		width: 100,
		render: (row: DomainHealth) => statusTag(row.dmarc_status),
	},
	{
		title: 'MX',
		key: 'mx_status',
		width: 100,
		render: (row: DomainHealth) => statusTag(row.mx_status),
	},
	{
		title: 'Issues',
		key: 'issues',
		ellipsis: true,
		render: (row: DomainHealth) => {
			const issues = row.issues
			if (!issues || issues.length === 0) return h('span', { style: 'color: var(--n-text-color-3)' }, '—')
			return h('span', { style: 'color: #d03050; font-size: 12px' }, issues.join('; '))
		},
	},
	{
		title: 'Actions',
		key: 'actions',
		width: 100,
		render: (row: DomainHealth) => {
			return h(NButton, {
				size: 'tiny',
				type: 'primary',
				quaternary: true,
				onClick: () => handleCheckSingle(row.domain),
			}, () => 'Recheck')
		},
	},
]

function updateSummary() {
	healthyCount.value = 0
	warningCount.value = 0
	errorCount.value = 0
	for (const d of domains.value) {
		const statuses = [d.spf_status, d.dkim_status, d.dmarc_status, d.mx_status]
		const missing = statuses.filter(s => s === 'missing').length
		if (missing === 0) healthyCount.value++
		else if (missing <= 1) warningCount.value++
		else errorCount.value++
	}
}

async function fetchDomains() {
	loading.value = true
	try {
		const res = await listDomainHealth()
		domains.value = Array.isArray(res) ? res : []
		updateSummary()
	} catch (e) {
		console.error('Failed to fetch domains', e)
	} finally {
		loading.value = false
	}
}

async function handleCheckAll() {
	checkingAll.value = true
	try {
		await checkAllDomains()
		message.success('All domains checked')
		await fetchDomains()
	} catch (e) {
		message.error('Domain check failed')
	} finally {
		checkingAll.value = false
	}
}

async function handleCheckSingle(domain: string) {
	try {
		await checkDomain({ domain })
		message.success(`${domain} checked`)
		await fetchDomains()
	} catch (e) {
		message.error(`Check failed for ${domain}`)
	}
}

onMounted(() => {
	fetchDomains()
})
</script>

<style lang="scss" scoped>
.domain-health-page {
	padding: 20px;
}

.actions-row {
	margin-bottom: 16px;
}

.summary-row {
	display: flex;
	gap: 16px;
	margin-bottom: 16px;
}

.summary-card {
	text-align: center;
	flex: 1;

	.summary-value {
		font-size: 28px;
		font-weight: 700;
	}

	.summary-label {
		font-size: 13px;
		color: var(--n-text-color-3);
		margin-top: 2px;
	}
}
</style>
