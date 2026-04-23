import { instance } from '@/api'
import i18n from '@/i18n'

const { t } = i18n.global

export interface VariantInput {
	subject: string
	body_html: string
	body_text: string
}

export const createAbTest = (data: {
	sequence_id: number
	step_id: number
	name?: string
	split_ratio?: number
	winner_criteria?: string
	variant_a: VariantInput
	variant_b: VariantInput
}) => {
	return instance.post('/abtest/create', data, {
		fetchOptions: {
			successMessage: true,
		},
	})
}

export const getAbTest = (params: { id: number }) => {
	return instance.get('/abtest/get', { params })
}

export const listAbTests = (params: { sequence_id: number }) => {
	return instance.get('/abtest/list', { params })
}

export const getAbTestResults = (params: { id: number }) => {
	return instance.get('/abtest/results', { params })
}

export const pickWinner = (data: { id: number; winner_variant: number }) => {
	return instance.post('/abtest/pick_winner', data, {
		fetchOptions: {
			successMessage: true,
		},
	})
}

export const deleteAbTest = (data: { id: number }) => {
	return instance.post('/abtest/delete', data, {
		fetchOptions: {
			successMessage: true,
		},
	})
}
