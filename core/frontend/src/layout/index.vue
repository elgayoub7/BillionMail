<template>
	<div class="app-layout">
		<sidebar></sidebar>
		<div class="app-main-wrapper" :style="{ marginLeft: isCollapse ? '64px' : '240px' }">
			<div class="app-content-scroll" :style="{ paddingTop: '48px' }" @scroll="handleScroll">
				<app-header :top="scrollTop" @toggle-sidebar="handleToggleSidebar"></app-header>
				<app-main></app-main>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { Sidebar, AppHeader, AppMain } from './components'
import { useGlobalStore } from '@/store'

const globalStore = useGlobalStore()
const { isCollapse } = storeToRefs(globalStore)

const scrollTop = ref(0)

const handleScroll = (e: Event) => {
	scrollTop.value = (e.target as HTMLElement).scrollTop || 0
}

const handleToggleSidebar = () => {
	globalStore.setCollapse()
}
</script>

<style lang="scss" scoped>
.app-layout {
	display: flex;
	width: 100%;
	height: 100%;
	overflow: hidden;
}

.app-main-wrapper {
	flex: 1;
	height: 100%;
	transition: margin-left 0.25s cubic-bezier(0.4, 0, 0.2, 1);
	overflow: hidden;
}

.app-content-scroll {
	height: 100%;
	overflow-y: auto;
	overflow-x: hidden;
}
</style>
