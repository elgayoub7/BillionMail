import { instance } from '@/api'
import i18n from '@/i18n'

const { t } = i18n.global

export const getScoringLeads = (params: { level?: string; page?: number; page_size?: number }) => {
	return instance.get('/scoring/leads', { params })
}

export const getScoringStats = () => {
	return instance.get('/scoring/stats')
}

export const getLeadScore = (params: { email: string }) => {
	return instance.get('/scoring/lead', { params })
}

export const recalculateScores = () => {
	return instance.post('/scoring/recalculate', {}, {
		fetchOptions: {
			successMessage: true,
		},
	})
}
