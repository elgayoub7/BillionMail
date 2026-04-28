<template>
	<div class="app-header" :style="{ top: `${top}px` }">
		<div class="header-left">
			<n-button class="icon-btn" :bordered="false" @click="$emit('toggleSidebar')">
				<i class="icon" :class="isCollapse ? 'i-mdi:menu-open' : 'i-mdi:menu-close'"></i>
			</n-button>
		</div>

		<div class="header-right">
			<n-button class="icon-btn" :bordered="false" @click="handleSetTheme">
				<i class="icon" :class="theme === 'light' ? 'i-ri:sun-line' : 'i-ri:moon-line'"></i>
			</n-button>
			<InstanceSwitcher />
			<n-dropdown size="large" :options="userOptions" @select="handleUserAction">
				<n-button class="icon-btn" :bordered="false">
					<i class="icon i-mdi:account-circle-outline"></i>
				</n-button>
			</n-dropdown>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { DropdownOption } from 'naive-ui'
import { useUserStore, useGlobalStore, useThemeStore } from '@/store'
import InstanceSwitcher from './InstanceSwitcher.vue'

defineProps({
	top: {
		type: Number,
		default: 0,
	},
})

defineEmits(['toggleSidebar'])

const { t } = useI18n()

const userStore = useUserStore()
const globalStore = useGlobalStore()
const { isCollapse } = storeToRefs(globalStore)

const themeStore = useThemeStore()
const { theme } = storeToRefs(themeStore)

const handleSetTheme = () => {
	themeStore.setTheme(theme.value === 'dark' ? 'light' : 'dark')
}

const userOptions = ref<DropdownOption[]>([
	{
		label: t('layout.menu.logout'),
		key: 'logout',
	},
])

const handleUserAction = (key: string) => {
	if (key === 'logout') {
		userStore.logout()
	}
}
</script>

<style lang="scss" scoped>
.app-header {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 48px;
	padding: 0 20px;
	background: #0d0f16;
	border-bottom: 1px solid #1e2030;
	z-index: 1000;
}

.header-left,
.header-right {
	display: flex;
	align-items: center;
	gap: 8px;
}

.icon-btn {
	--n-width: 36px;
	--n-height: 36px;
	--n-padding: 0;
	--n-font-size: 20px;
	--n-text-color: #8892a8;
	--n-ripple-color: none;

	&:hover {
		--n-text-color: #e2e8f0;
	}
}
</style>
