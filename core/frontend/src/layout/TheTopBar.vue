<template>
	<div class="topbar">
		<div class="topbar-left">
			<span class="section-name">{{ currentSection?.label || '' }}</span>
			<div v-if="subNav.length > 1" class="sub-nav">
				<router-link
					v-for="item in subNav" :key="item.path"
					:to="item.path"
					class="sub-nav-item"
					:class="{ active: route.path.startsWith(item.path) }">
					{{ item.label }}
				</router-link>
			</div>
		</div>
		<div class="topbar-right">
			<n-dropdown v-if="langOptions.length > 0" size="small" :options="langOptions" @select="handleLangAction">
				<button class="topbar-btn">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"/></svg>
				</button>
			</n-dropdown>
			<n-dropdown size="small" :options="userOptions" @select="handleUserAction">
				<button class="topbar-btn">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
				</button>
			</n-dropdown>
		</div>
	</div>
</template>
<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { DropdownOption } from 'naive-ui'
import { useGlobalStore, useUserStore } from '@/store'

const route = useRoute()
const { t } = useI18n()
const globalStore = useGlobalStore()
const userStore = useUserStore()
const { langList } = storeToRefs(globalStore)

const sectionMap = {
	dashboard: [
		{ label: 'Dashboard', path: '/cold-dashboard' },
		{ label: 'Overview', path: '/overview' },
	],
	campaigns: [
		{ label: 'Sequences', path: '/sequences' },
		{ label: 'Email Marketing', path: '/market' },
		{ label: 'AB Tests', path: '/abtest' },
	],
	templates: [
		{ label: 'Templates', path: '/template' },
	],
	contacts: [
		{ label: 'Contacts', path: '/contacts' },
		{ label: 'Lead Scoring', path: '/scoring' },
	],
	settings: [
		{ label: 'Settings', path: '/settings' },
		{ label: 'Domaines', path: '/domain' },
		{ label: 'Health', path: '/domain-health' },
		{ label: 'Mailboxes', path: '/mailbox' },
		{ label: 'SMTP', path: '/smtp' },
		{ label: 'API', path: '/api' },
		{ label: 'Logs', path: '/logs' },
	],
}

function getActiveKey() {
	const p = route.path
	if (p === '/' || p.startsWith('/cold-dashboard') || p.startsWith('/overview')) return 'dashboard'
	if (p.startsWith('/sequences') || p.startsWith('/market') || p.startsWith('/abtest')) return 'campaigns'
	if (p.startsWith('/template')) return 'templates'
	if (p.startsWith('/contacts') || p.startsWith('/scoring')) return 'contacts'
	return 'settings'
}

const currentSection = computed(() => {
	const key = getActiveKey()
	return { key, label: key.charAt(0).toUpperCase() + key.slice(1) }
})

const subNav = computed(() => {
	return sectionMap[getActiveKey()] || []
})

const langOptions = ref<DropdownOption[]>([])
const userOptions = ref<DropdownOption[]>([{ label: 'Logout', key: 'logout' }])

function handleLangAction(key: string) { globalStore.setLang(key).then(() => window.location.reload()) }
function handleUserAction(key: string) { if (key === 'logout') userStore.logout() }

onMounted(() => { langOptions.value = langList.value.map(i => ({ label: i.cn, key: i.name })) })
</script>
<style scoped>
.topbar { height: 48px; display: flex; justify-content: space-between; align-items: center; padding: 0 20px; background: #1a1d27; border-bottom: 1px solid #2e3142; flex-shrink: 0; }
.topbar-left { display: flex; align-items: center; gap: 16px; }
.section-name { font-size: 15px; font-weight: 600; color: #f1f3f7; }
.sub-nav { display: flex; gap: 2px; background: #232630; border-radius: 6px; padding: 2px; }
.sub-nav-item { padding: 4px 12px; border-radius: 4px; font-size: 13px; color: #6b7084; text-decoration: none; transition: all 0.2s ease; }
.sub-nav-item:hover { color: #a1a5b3; }
.sub-nav-item.active { color: #f1f3f7; background: #2a2d3a; }
.topbar-right { display: flex; align-items: center; gap: 8px; }
.topbar-btn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: none; background: transparent; color: #6b7084; cursor: pointer; transition: all 0.2s ease; }
.topbar-btn:hover { color: #a1a5b3; background: #232630; }
</style>
