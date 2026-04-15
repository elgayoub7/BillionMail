export interface Sequence {
  id: number
  name: string
  description: string
  status: number
  step_count: number
  total_enrolled: number
  total_completed: number
  total_bounced: number
  group_name: string
  create_time: number
  update_time: number
}

export interface SequenceParams {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

export interface SequenceDetail {
  id: number
  name: string
  description: string
  status: number
  addresser: string
  full_name: string
  group_id: number
  tag_ids: number[]
  tag_logic: string
  track_open: number
  track_click: number
  unsubscribe: number
  total_enrolled: number
  total_completed: number
  total_unsubscribed: number
  total_bounced: number
  group_name: string
  create_time: number
  update_time: number
  steps: SequenceStepItem[]
}

export interface SequenceStepItem {
  id: number
  step_order: number
  step_type: 'email' | 'wait' | 'condition'
  subject: string
  template_id: number
  template_name: string
  wait_days: number
  wait_hours: number
  condition_type: string
  condition_step_id: number
  on_true_go_to: number
  on_false_go_to: number
  sent_count: number
  opened_count: number
  clicked_count: number
  bounced_count: number
}

export interface StepInput {
  step_order: number
  step_type: 'email' | 'wait' | 'condition'
  subject?: string
  template_id?: number
  wait_days?: number
  wait_hours?: number
  condition_type?: string
  condition_step_id?: number
  on_true_go_to?: number
  on_false_go_to?: number
}

export interface Enrollment {
  id: number
  sequence_id: number
  contact_id: number
  email: string
  current_step: number
  status: number
  enrolled_at: number
  last_email_sent_at: number
  total_emails_sent: number
  total_opens: number
  total_clicks: number
}

export interface EnrollmentParams {
  sequence_id: number
  page: number
  page_size: number
  status?: number
}
