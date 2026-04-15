package entity

// Sequence represents an automated follow-up email sequence
type Sequence struct {
	Id                int    `json:"id"                dc:"Sequence ID"`
	Name              string `json:"name"              dc:"Sequence Name"`
	Description       string `json:"description"       dc:"Description"`
	Status            int    `json:"status"            dc:"Status: 0=draft, 1=active, 2=paused, 3=archived"`
	Addresser         string `json:"addresser"         dc:"Sender Email"`
	FullName          string `json:"full_name"         dc:"Sender Display Name"`
	GroupId           int    `json:"group_id"          dc:"Contact Group ID"`
	TagIds            string `json:"tag_ids"           dc:"Tag IDs (JSON)"`
	TagLogic          string `json:"tag_logic"         dc:"Tag Logic: AND/OR/NOT"`
	TrackOpen         int    `json:"track_open"        dc:"Track Opens"`
	TrackClick        int    `json:"track_click"       dc:"Track Clicks"`
	Unsubscribe       int    `json:"unsubscribe"       dc:"Allow Unsubscribe"`
	TotalEnrolled     int    `json:"total_enrolled"    dc:"Total Enrolled"`
	TotalCompleted    int    `json:"total_completed"   dc:"Total Completed"`
	TotalUnsubscribed int    `json:"total_unsubscribed" dc:"Total Unsubscribed"`
	TotalBounced      int    `json:"total_bounced"     dc:"Total Bounced"`
	CreateTime        int    `json:"create_time"       dc:"Create Time"`
	UpdateTime        int    `json:"update_time"       dc:"Update Time"`
}

// SequenceStep represents a single step within a sequence
type SequenceStep struct {
	Id             int    `json:"id"               dc:"Step ID"`
	SequenceId     int    `json:"sequence_id"      dc:"Sequence ID"`
	StepOrder      int    `json:"step_order"       dc:"Step Order"`
	StepType       string `json:"step_type"        dc:"Step Type: email/wait/condition"`
	Subject        string `json:"subject"          dc:"Email Subject"`
	TemplateId     int    `json:"template_id"      dc:"Template ID"`
	WaitDays       int    `json:"wait_days"        dc:"Wait Days"`
	WaitHours      int    `json:"wait_hours"       dc:"Wait Hours"`
	ConditionType  string `json:"condition_type"   dc:"Condition: opened/not_opened/clicked/not_clicked/bounced"`
	ConditionStepId int   `json:"condition_step_id" dc:"Reference Step ID for Condition"`
	OnTrueGoTo     int    `json:"on_true_go_to"    dc:"Jump to Step Order if True"`
	OnFalseGoTo    int    `json:"on_false_go_to"   dc:"Jump to Step Order if False"`
	SentCount      int    `json:"sent_count"       dc:"Sent Count"`
	OpenedCount    int    `json:"opened_count"     dc:"Opened Count"`
	ClickedCount   int    `json:"clicked_count"    dc:"Clicked Count"`
	BouncedCount   int    `json:"bounced_count"    dc:"Bounced Count"`
	CreateTime     int    `json:"create_time"      dc:"Create Time"`
	UpdateTime     int    `json:"update_time"      dc:"Update Time"`
	TemplateName   string `json:"template_name"    dc:"Template Name (joined)"`
}

// SequenceEnrollment represents a contact enrolled in a sequence
type SequenceEnrollment struct {
	Id                   int    `json:"id"                     dc:"Enrollment ID"`
	SequenceId           int    `json:"sequence_id"            dc:"Sequence ID"`
	ContactId            int    `json:"contact_id"             dc:"Contact ID"`
	Email                string `json:"email"                  dc:"Contact Email"`
	GroupId              int    `json:"group_id"               dc:"Group ID"`
	CurrentStep          int    `json:"current_step"           dc:"Current Step Order (1-indexed)"`
	Status               int    `json:"status"                 dc:"Status: 0=active, 1=completed, 2=paused, 3=exited"`
	EnrolledAt           int    `json:"enrolled_at"            dc:"Enrolled At"`
	CurrentStepEnteredAt int    `json:"current_step_entered_at" dc:"Current Step Entered At"`
	LastEmailSentAt      int    `json:"last_email_sent_at"     dc:"Last Email Sent At"`
	CompletedAt          int    `json:"completed_at"           dc:"Completed At"`
	TotalEmailsSent      int    `json:"total_emails_sent"      dc:"Total Emails Sent"`
	TotalOpens           int    `json:"total_opens"            dc:"Total Opens"`
	TotalClicks          int    `json:"total_clicks"           dc:"Total Clicks"`
}

// SequenceEmailTask links a sequence step execution to an actual email_task
type SequenceEmailTask struct {
	Id           int    `json:"id"            dc:"ID"`
	SequenceId   int    `json:"sequence_id"   dc:"Sequence ID"`
	EnrollmentId int    `json:"enrollment_id" dc:"Enrollment ID"`
	StepId       int    `json:"step_id"       dc:"Step ID"`
	EmailTaskId  int    `json:"email_task_id" dc:"Email Task ID"`
	ContactEmail string `json:"contact_email" dc:"Contact Email"`
	Status       int    `json:"status"        dc:"Status: 0=pending, 1=sent, 2=failed"`
	SentAt       int    `json:"sent_at"       dc:"Sent At"`
	MessageId    string `json:"message_id"    dc:"Message ID"`
}
