<template>
	<div class="login-page">
		<div class="login-orb login-orb-1"></div>
		<div class="login-orb login-orb-2"></div>
		<div class="login-card">
			<div class="login-logo">
				<img src="@/assets/images/logo.png" alt="Logo" />
			</div>
			<h1 class="login-title">SelectCircle</h1>
			<p class="login-subtitle">Cold Email Platform</p>
			<n-form ref="formRef" size="large" :model="form" :rules="rules">
				<n-form-item :show-label="false" path="username">
					<n-input v-model:value="form.username" :placeholder="t('login.form.usernamePlaceholder')" />
				</n-form-item>
				<n-form-item :show-label="false" path="password">
					<n-input v-model:value="form.password" type="password" show-password-on="click"
						:placeholder="t('login.form.passwordPlaceholder')" @keyup.enter="handleLogin" />
				</n-form-item>
				<n-form-item v-if="isCode" :show-label="false" path="validate_code">
					<n-input v-model:value="form.validate_code" class="flex-1"
						:placeholder="t('login.form.captcha')" @keydown.enter="handleLogin" />
					<n-spin size="small" :show="codeLoading">
						<div class="code" @click="getCode()">
							<img class="w-full h-full" :src="codeUrl" :alt="t('login.form.captcha')" />
						</div>
					</n-spin>
				</n-form-item>
				<n-form-item :show-label="false" :show-feedback="false">
					<n-button type="primary" size="large" :loading="loading" :disabled="loading" block @click="handleLogin">
						{{ t('login.form.loginButton') }}
					</n-button>
				</n-form-item>
			</n-form>
		</div>
	</div>
</template>
<script lang="ts" setup>
import { useUserStore } from '@/store'
import { isObject } from '@/utils'
import { getValidateCode, login } from '@/api/modules/user'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const formRef = useTemplateRef('formRef')
const isCode = ref(false)
const codeUrl = ref('')
const codeLoading = ref(false)

const form = reactive({ username: '', password: '', validate_code: '', validate_code_id: '' })
const rules = {
	username: { required: true, message: t('login.validation.usernameRequired'), trigger: ['blur', 'input'] },
	password: { required: true, message: t('login.validation.passwordRequired'), trigger: ['blur', 'input'] },
	validate_code: { required: true, trigger: ['blur', 'input'], message: t('login.validation.captchaRequired') },
}

interface CodeResponse { mustValidateCode: boolean; validateCodeBase64: string; validateCodeId: string }
interface LoginResponse { token: string; refresh_token: string; ttl: number }

const getCode = async () => {
	try {
		codeLoading.value = true
		const res = await getValidateCode()
		if (isObject<CodeResponse>(res)) {
			isCode.value = res.mustValidateCode
			if (res.mustValidateCode) { codeUrl.value = res.validateCodeBase64; form.validate_code_id = res.validateCodeId }
		}
	} finally { codeLoading.value = false }
}

const loading = ref(false)
const handleLogin = async () => {
	try {
		await formRef.value?.validate()
		loading.value = true
		const res = await login(toRaw(form))
		if (isObject<LoginResponse>(res)) {
			userStore.setLoginInfo({ token: res.token, refresh_token: res.refresh_token, ttl: res.ttl })
			setTimeout(() => { router.push('/') }, 1000)
		}
	} catch { getCode() } finally { loading.value = false }
}
getCode()
</script>
<style scoped>
.login-page {
	min-height: 100vh; display: flex; align-items: center; justify-content: center;
	background: #0a0b10; position: relative; overflow: hidden;
}
.login-orb {
	position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.15;
}
.login-orb-1 { width: 500px; height: 500px; background: #6c5ce7; top: -10%; right: -10%; }
.login-orb-2 { width: 400px; height: 400px; background: #00d68f; bottom: -10%; left: -10%; }
.login-card {
	width: 100%; max-width: 400px; padding: 48px 32px 56px; border-radius: 16px; z-index: 1;
	background: rgba(26, 29, 39, 0.8); backdrop-filter: blur(24px);
	border: 1px solid rgba(46, 49, 66, 0.5); box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}
.login-logo { display: flex; justify-content: center; margin-bottom: 16px; }
.login-logo img { width: 48px; height: 48px; }
.login-title { text-align: center; font-size: 22px; font-weight: 700; color: #f1f3f7; margin: 0 0 4px; }
.login-subtitle { text-align: center; font-size: 14px; color: #6b7084; margin: 0 0 28px; }
.code { width: 120px; height: 40px; margin-left: 12px; border-radius: 6px; border: 1px solid #2e3142; overflow: hidden; cursor: pointer; }
</style>
