import { instance } from '@/api'
import i18n from '@/i18n'

const { t } = i18n.global

export const listDomainHealth = () => {
	return instance.get('/domainhealth/list')
}

export const checkDomain = (data: { domain: string }) => {
	return instance.post('/domainhealth/check', data, {
		fetchOptions: {
			successMessage: true,
		},
	})
}

export const checkAllDomains = () => {
	return instance.post('/domainhealth/check_all', {}, {
		fetchOptions: {
			successMessage: true,
		},
	})
}
