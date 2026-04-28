<template>
	<div class="sidebar" :class="{ collapsed: isCollapse }">
		<!-- Logo -->
		<div class="sidebar-logo">
			<router-link to="/overview" class="logo-link">
				<img class="logo-icon" src="@/assets/images/logo.png" alt="BillionMail" />
				<span v-show="!isCollapse" class="logo-text">BillionMail</span>
			</router-link>
		</div>

		<!-- Navigation -->
		<nav class="sidebar-nav">
			<template v-for="group in menuGroups" :key="group.label">
				<div class="nav-group">
					<div class="nav-group-label">
						<span v-if="isCollapse" class="nav-group-dot"></span>
						<span v-else>{{ group.label }}</span>
					</div>
					<router-link
						v-for="item in group.items"
						:key="item.key"
						:to="item.route"
						class="nav-item"
						:class="{ active: isActive(item.key) }"
						:title="isCollapse ? item.label : ''">
						<span class="nav-item-icon">
							<i :class="item.icon"></i>
							<span v-if="item.badge && item.badge > 0" class="nav-badge">{{ item.badge > 99 ? '99+' : item.badge }}</span>
						</span>
						<span v-show="!isCollapse" class="nav-item-label">{{ item.label }}</span>
					</router-link>
				</div>
			</template>
		</nav>

		<!-- Footer -->
		<div class="sidebar-footer">
			<router-link to="/settings" class="nav-item" title="Settings">
				<span class="nav-item-icon"><i class="i-mdi:cog-outline"></i></span>
				<span v-show="!isCollapse" class="nav-item-label">Settings</span>
			</router-link>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '@/store'

const route = useRoute()
const globalStore = useGlobalStore()
const { isCollapse } = storeToRefs(globalStore)

const unreadCount = ref(0)

const isActive = (key: string) => {
	return String(route.meta?.key || '').startsWith(key)
}

interface MenuItem {
	key: string
	label: string
	icon: string
	route: string
	badge?: number
}

interface MenuGroup {
	label: string
	items: MenuItem[]
}

const menuGroups = computed<MenuGroup[]>(() => [
	{
		label: 'DASHBOARD',
		items: [
			{ key: 'overview', label: 'Dashboard', icon: 'i-mdi:view-dashboard-outline', route: '/overview' },
		],
	},
	{
		label: 'OUTREACH',
		items: [
			{ key: 'sequences', label: 'Campaigns', icon: 'i-mdi:rocket-launch-outline', route: '/sequences' },
			{ key: 'inbox', label: 'Inbox', icon: 'i-mdi:email-outline', route: '/inbox', badge: unreadCount.value },
			{ key: 'template', label: 'Templates', icon: 'i-mdi:file-document-edit-outline', route: '/template' },
		],
	},
	{
		label: 'CONTACTS',
		items: [
			{ key: 'contacts', label: 'Contacts', icon: 'i-mdi:account-group-outline', route: '/contacts' },
		],
	},
	{
		label: 'ANALYTICS',
		items: [
			{ key: 'analytics', label: 'Analytics', icon: 'i-mdi:chart-areaspline', route: '/analytics' },
			{ key: 'logs', label: 'Activity Log', icon: 'i-mdi:text-box-outline', route: '/logs' },
			{ key: 'domain', label: 'Domains', icon: 'i-mdi:dns-outline', route: '/domain' },
			{ key: 'mailbox', label: 'Mailboxes', icon: 'i-mdi:email-fast-outline', route: '/mailbox' },
		],
	},
])

// Future: poll unread count for inbox badge
// onMounted(() => {
//   pollUnread()
// })
</script>

<style lang="scss" scoped>
.sidebar {
	width: 240px;
	height: 100vh;
	background: #0d0f16;
	border-right: 1px solid #1e2030;
	display: flex;
	flex-direction: column;
	transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
	flex-shrink: 0;
	overflow: hidden;

	&.collapsed {
		width: 64px;

		.sidebar-logo {
			padding: 16px 0;
			justify-content: center;
		}

		.nav-item {
			justify-content: center;
			padding: 10px 0;

			.nav-item-label {
				display: none;
			}
		}

		.nav-group-label {
			justify-content: center;
		}
	}
}

.sidebar-logo {
	padding: 16px 20px;
	border-bottom: 1px solid #1e2030;
	display: flex;
	align-items: center;
	transition: all 0.25s ease;

	.logo-link {
		display: flex;
		align-items: center;
		gap: 10px;
		text-decoration: none;
		white-space: nowrap;
		overflow: hidden;
	}

	.logo-icon {
		width: 32px;
		height: 32px;
		flex-shrink: 0;
	}

	.logo-text {
		font-size: 17px;
		font-weight: 700;
		color: #e2e8f0;
		letter-spacing: -0.3px;
	}
}

.sidebar-nav {
	flex: 1;
	overflow-y: auto;
	overflow-x: hidden;
	padding: 8px 0;

	&::-webkit-scrollbar {
		width: 0;
	}
}

.nav-group {
	margin-bottom: 4px;
}

.nav-group-label {
	padding: 12px 20px 6px;
	font-size: 11px;
	font-weight: 600;
	color: #4a5568;
	letter-spacing: 1px;
	text-transform: uppercase;
	white-space: nowrap;
	display: flex;
	align-items: center;

	.nav-group-dot {
		width: 4px;
		height: 4px;
		border-radius: 50%;
		background: #4a5568;
	}
}

.nav-item {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 10px 20px;
	color: #8892a8;
	text-decoration: none;
	font-size: 14px;
	font-weight: 500;
	transition: all 0.15s ease;
	cursor: pointer;
	white-space: nowrap;
	position: relative;
	margin: 1px 8px;
	border-radius: 8px;

	&:hover {
		color: #e2e8f0;
		background: rgba(108, 92, 231, 0.08);
	}

	&.active {
		color: #a78bfa;
		background: rgba(108, 92, 231, 0.15);

		.nav-item-icon i {
			color: #a78bfa;
		}
	}
}

.nav-item-icon {
	position: relative;
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;

	i {
		font-size: 20px;
		transition: color 0.15s ease;
	}
}

.nav-badge {
	position: absolute;
	top: -6px;
	right: -8px;
	min-width: 18px;
	height: 18px;
	padding: 0 5px;
	background: #ef4444;
	color: #fff;
	font-size: 10px;
	font-weight: 700;
	border-radius: 9px;
	display: flex;
	align-items: center;
	justify-content: center;
	line-height: 1;
}

.nav-item-label {
	overflow: hidden;
	text-overflow: ellipsis;
}

.sidebar-footer {
	border-top: 1px solid #1e2030;
	padding: 8px 0;

	.nav-item {
		margin: 0 8px;
	}
}
</style>
