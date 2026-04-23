import { instance } from '@/api'

export const getDashboardStats = () => {
	return instance.get('/cold-dashboard/stats')
}

export const getActiveSequences = () => {
	return instance.get('/cold-dashboard/active-sequences')
}

export const getAlerts = (params?: { limit?: number }) => {
	return instance.get('/cold-dashboard/alerts', { params })
}
