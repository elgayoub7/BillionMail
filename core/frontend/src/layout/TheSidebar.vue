<template>
	<aside class="sidebar">
		<router-link to="/" class="sidebar-logo">
			<img src="@/assets/images/logo.png" alt="Logo" />
		</router-link>
		<nav class="sidebar-nav">
			<router-link v-for="section in sections" :key="section.key" :to="section.path"
				class="sidebar-item" :class="{ active: isActive(section) }">
				<div class="sidebar-icon" v-html="section.icon"></div>
				<span class="sidebar-tooltip">{{ section.label }}</span>
			</router-link>
		</nav>
		<div class="sidebar-bottom">
			<button class="sidebar-item" @click="toggleTheme">
				<svg v-if="isDark" class="sidebar-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"/></svg>
				<svg v-else class="sidebar-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"/></svg>
			</button>
			<button class="sidebar-item" @click="handleLogout">
				<svg class="sidebar-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9"/></svg>
			</button>
		</div>
	</aside>
</template>
<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useThemeStore, useUserStore } from '@/store'

const route = useRoute()
const themeStore = useThemeStore()
const userStore = useUserStore()
const { theme } = storeToRefs(themeStore)
const isDark = computed(() => theme.value === 'dark')

const sections = [
	{ key: 'dashboard', label: 'Dashboard', path: '/cold-dashboard', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"/></svg>' },
	{ key: 'campaigns', label: 'Campagnes', path: '/sequences', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"/></svg>' },
	{ key: 'templates', label: 'Templates', path: '/template', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>' },
	{ key: 'contacts', label: 'Contacts', path: '/contacts', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"/></svg>' },
	{ key: 'settings', label: 'Parametres', path: '/settings', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 010 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 010-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.281z"/></svg><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>' },
]

function isActive(section) {
	const p = route.path
	switch (section.key) {
		case 'dashboard': return p === '/' || p.startsWith('/cold-dashboard') || p.startsWith('/overview')
		case 'campaigns': return p.startsWith('/sequences') || p.startsWith('/market') || p.startsWith('/abtest')
		case 'templates': return p.startsWith('/template')
		case 'contacts': return p.startsWith('/contacts') || p.startsWith('/scoring')
		case 'settings': return p.startsWith('/settings') || p.startsWith('/domain') || p.startsWith('/mailbox') || p.startsWith('/smtp') || p.startsWith('/api') || p.startsWith('/logs')
		default: return false
	}
}

function toggleTheme() {
	themeStore.setTheme(isDark.value ? 'light' : 'dark')
	themeStore.changeTheme()
}
function handleLogout() { userStore.logout() }
</script>
<style scoped>
.sidebar { width: 60px; height: 100vh; display: flex; flex-direction: column; align-items: center; background: #0f1117; border-right: 1px solid #2e3142; z-index: 1010; flex-shrink: 0; }
.sidebar-logo { padding: 16px 0; margin-bottom: 8px; border-bottom: 1px solid #2e3142; width: 100%; display: flex; justify-content: center; }
.sidebar-logo img { width: 32px; height: 32px; }
.sidebar-nav { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 12px 0; overflow-y: auto; }
.sidebar-bottom { display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 12px 0; border-top: 1px solid #2e3142; }
.sidebar-item { position: relative; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; border-radius: 10px; color: #6b7084; transition: all 0.2s ease; cursor: pointer; text-decoration: none; border: none; background: transparent; }
.sidebar-item:hover { color: #a1a5b3; background: #1a1d27; }
.sidebar-item.active { color: #6c5ce7; background: rgba(108, 92, 231, 0.12); }
.sidebar-item.active::before { content: ''; position: absolute; left: -10px; top: 50%; transform: translateY(-50%); width: 3px; height: 20px; background: #6c5ce7; border-radius: 0 3px 3px 0; }
.sidebar-icon { width: 20px; height: 20px; }
.sidebar-tooltip { position: absolute; left: 54px; top: 50%; transform: translateY(-50%); background: #1a1d27; color: #f1f3f7; padding: 4px 10px; border-radius: 6px; font-size: 13px; white-space: nowrap; pointer-events: none; opacity: 0; transition: opacity 0.15s ease; border: 1px solid #2e3142; z-index: 100; }
.sidebar-item:hover .sidebar-tooltip { opacity: 1; }
</style>
