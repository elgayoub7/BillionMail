<template>
	<div class="inbox-container">
		<div class="inbox-header">
			<h2>Inbox</h2>
			<n-button size="small" @click="openExternal">
				<template #icon><i class="i-mdi:open-in-new"></i></template>
				Open in new tab
			</n-button>
		</div>
		<div class="inbox-frame-wrapper">
			<iframe
				ref="iframeRef"
				src="/roundcube/"
				class="inbox-frame"
				allow="clipboard-read; clipboard-write"
				sandbox="allow-same-origin allow-scripts allow-forms allow-popups"
			></iframe>
		</div>
	</div>
</template>

<script lang="ts" setup>
const iframeRef = ref<HTMLIFrameElement>()

const openExternal = () => {
	window.open('/roundcube/', '_blank')
}
</script>

<style lang="scss" scoped>
.inbox-container {
	height: 100%;
	display: flex;
	flex-direction: column;
}

.inbox-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16px 20px;
	border-bottom: 1px solid #1e2030;

	h2 {
		margin: 0;
		font-size: 18px;
		font-weight: 600;
		color: #e2e8f0;
	}
}

.inbox-frame-wrapper {
	flex: 1;
	position: relative;
	overflow: hidden;
}

.inbox-frame {
	width: 100%;
	height: 100%;
	border: none;
	background: #1a1d27;
}
</style>
