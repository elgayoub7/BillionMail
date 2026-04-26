export interface Template {
	id: number
	temp_name: string
	add_type: number
	html_content: string
	drag_data: string
	create_time: number
	update_time: number
	chat_id: string
	ab_enabled: boolean
	ab_split_ratio: number
	ab_winner_criteria: string
	variant_b_subject: string
	variant_b_html: string
	isEdit?: boolean
	edit_name?: string
}

export interface TemplateParams {
	page: number
	page_size: number
	keyword: string
}
