<template>
	<div class="app-layout">
		<TheSidebar />
		<div class="app-main">
			<TheTopBar />
			<div class="app-content">
				<router-view :key="key" v-slot="{ Component, route }">
					<component :is="Component" :key="route.path" />
				</router-view>
			</div>
		</div>
	</div>
</template>
<script lang="ts" setup>
import TheSidebar from './TheSidebar.vue'
import TheTopBar from './TheTopBar.vue'
const key = computed(() => Math.random())
</script>
<style scoped>
.app-layout { display: flex; height: 100vh; width: 100vw; overflow: hidden; }
.app-main { flex: 1; display: flex; flex-direction: column; overflow: hidden; background: #0f1117; }
.app-content { flex: 1; overflow-y: auto; }
</style>
