package v1

import (
	"billionmail-core/utility/types/api_v1"

	"github.com/gogf/gf/v2/frame/g"
)

// StepInput represents a single step when creating/updating a sequence
type StepInput struct {
	StepOrder       int    `json:"step_order"        v:"required|min:1"     dc:"Step order"`
	StepType        string `json:"step_type"         v:"required|in:email,wait,condition" dc:"Step type"`
	Subject         string `json:"subject"           dc:"Email subject"`
	TemplateId      int    `json:"template_id"       dc:"Template ID"`
	WaitDays        int    `json:"wait_days"         dc:"Wait days"`
	WaitHours       int    `json:"wait_hours"        dc:"Wait hours"`
	ConditionType   string `json:"condition_type"    v:"in:opened,not_opened,clicked,not_clicked,bounced" dc:"Condition type"`
	ConditionStepId int    `json:"condition_step_id" dc:"Reference step order for condition"`
	OnTrueGoTo      int    `json:"on_true_go_to"     dc:"Step order to jump if condition true"`
	OnFalseGoTo     int    `json:"on_false_go_to"    dc:"Step order to jump if condition false"`
}

// --- Request types ---

type CreateSequenceReq struct {
	g.Meta        `path:"/sequence/create" method:"post" tags:"Sequence" summary:"Create sequence"`
	Authorization string      `json:"authorization" in:"header"`
	Name          string      `json:"name"          v:"required" dc:"Sequence name"`
	Description   string      `json:"description"   dc:"Description"`
	Addresser     string      `json:"addresser"     v:"required" dc:"Sender email"`
	FullName      string      `json:"full_name"     dc:"Sender display name"`
	GroupId       int         `json:"group_id"      v:"required|min:1" dc:"Contact group ID"`
	TagIds        []int       `json:"tag_ids"       dc:"Tag IDs"`
	TagLogic      string      `json:"tag_logic"     v:"in:AND,OR,NOT" dc:"Tag filter logic"`
	TrackOpen     int         `json:"track_open"    dc:"Track opens (0/1)"`
	TrackClick    int         `json:"track_click"   dc:"Track clicks (0/1)"`
	Unsubscribe   int         `json:"unsubscribe"   dc:"Allow unsubscribe (0/1)"`
	Steps         []StepInput `json:"steps"         v:"required" dc:"Sequence steps"`
}

type CreateSequenceRes struct {
	api_v1.StandardRes
	Data struct {
		Id int `json:"id"`
	} `json:"data"`
}

type UpdateSequenceReq struct {
	g.Meta        `path:"/sequence/update" method:"post" tags:"Sequence" summary:"Update sequence"`
	Authorization string      `json:"authorization" in:"header"`
	Id            int         `json:"id"            v:"required|min:1" dc:"Sequence ID"`
	Name          string      `json:"name"          dc:"Sequence name"`
	Description   string      `json:"description"   dc:"Description"`
	Addresser     string      `json:"addresser"     dc:"Sender email"`
	FullName      string      `json:"full_name"     dc:"Sender display name"`
	GroupId       int         `json:"group_id"      dc:"Contact group ID"`
	TagIds        []int       `json:"tag_ids"       dc:"Tag IDs"`
	TagLogic      string      `json:"tag_logic"     v:"in:AND,OR,NOT" dc:"Tag filter logic"`
	TrackOpen     int         `json:"track_open"    dc:"Track opens (0/1)"`
	TrackClick    int         `json:"track_click"   dc:"Track clicks (0/1)"`
	Unsubscribe   int         `json:"unsubscribe"   dc:"Allow unsubscribe (0/1)"`
	Steps         []StepInput `json:"steps"         dc:"Sequence steps"`
}

type UpdateSequenceRes struct {
	api_v1.StandardRes
}

type DeleteSequenceReq struct {
	g.Meta        `path:"/sequence/delete" method:"post" tags:"Sequence" summary:"Delete sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id"            v:"required|min:1" dc:"Sequence ID"`
}

type DeleteSequenceRes struct {
	api_v1.StandardRes
}

type ListSequencesReq struct {
	g.Meta        `path:"/sequence/list" method:"get" tags:"Sequence" summary:"List sequences"`
	Authorization string `json:"authorization" in:"header"`
	Page          int    `json:"page"          v:"required|min:1" dc:"Page number"`
	PageSize      int    `json:"page_size"     v:"required|min:1|max:100" dc:"Page size"`
	Keyword       string `json:"keyword"       dc:"Search keyword"`
	Status        int    `json:"status"        dc:"Filter by status (-1 = all)"`
}

type ListSequencesRes struct {
	api_v1.StandardRes
	Data struct {
		Total int                `json:"total"`
		List  []SequenceListItem `json:"list"`
	} `json:"data"`
}

type SequenceListItem struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Status         int    `json:"status"`
	StepCount      int    `json:"step_count"`
	TotalEnrolled  int    `json:"total_enrolled"`
	TotalCompleted int    `json:"total_completed"`
	TotalBounced   int    `json:"total_bounced"`
	GroupName      string `json:"group_name"`
	CreateTime     int    `json:"create_time"`
	UpdateTime     int    `json:"update_time"`
}

type FindSequenceReq struct {
	g.Meta        `path:"/sequence/find" method:"get" tags:"Sequence" summary:"Find sequence detail"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id"            v:"required|min:1" dc:"Sequence ID"`
}

type FindSequenceRes struct {
	api_v1.StandardRes
	Data *SequenceDetail `json:"data"`
}

type SequenceDetail struct {
	Id                int              `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Status            int              `json:"status"`
	Addresser         string           `json:"addresser"`
	FullName          string           `json:"full_name"`
	GroupId           int              `json:"group_id"`
	TagIds            []int            `json:"tag_ids"`
	TagLogic          string           `json:"tag_logic"`
	TrackOpen         int              `json:"track_open"`
	TrackClick        int              `json:"track_click"`
	Unsubscribe       int              `json:"unsubscribe"`
	TotalEnrolled     int              `json:"total_enrolled"`
	TotalCompleted    int              `json:"total_completed"`
	TotalUnsubscribed int              `json:"total_unsubscribed"`
	TotalBounced      int              `json:"total_bounced"`
	GroupName         string           `json:"group_name"`
	CreateTime        int              `json:"create_time"`
	UpdateTime        int              `json:"update_time"`
	Steps             []SequenceStepItem `json:"steps"`
}

type SequenceStepItem struct {
	Id              int    `json:"id"`
	StepOrder       int    `json:"step_order"`
	StepType        string `json:"step_type"`
	Subject         string `json:"subject"`
	TemplateId      int    `json:"template_id"`
	TemplateName    string `json:"template_name"`
	WaitDays        int    `json:"wait_days"`
	WaitHours       int    `json:"wait_hours"`
	ConditionType   string `json:"condition_type"`
	ConditionStepId int    `json:"condition_step_id"`
	OnTrueGoTo      int    `json:"on_true_go_to"`
	OnFalseGoTo     int    `json:"on_false_go_to"`
	SentCount       int    `json:"sent_count"`
	OpenedCount     int    `json:"opened_count"`
	ClickedCount    int    `json:"clicked_count"`
	BouncedCount    int    `json:"bounced_count"`
}

// --- Sequence lifecycle ---

type ActivateSequenceReq struct {
	g.Meta        `path:"/sequence/activate" method:"post" tags:"Sequence" summary:"Activate sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id"            v:"required|min:1" dc:"Sequence ID"`
}

type ActivateSequenceRes struct {
	api_v1.StandardRes
}

type PauseSequenceReq struct {
	g.Meta        `path:"/sequence/pause" method:"post" tags:"Sequence" summary:"Pause sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id"            v:"required|min:1" dc:"Sequence ID"`
}

type PauseSequenceRes struct {
	api_v1.StandardRes
}

type ResumeSequenceReq struct {
	g.Meta        `path:"/sequence/resume" method:"post" tags:"Sequence" summary:"Resume sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id"            v:"required|min:1" dc:"Sequence ID"`
}

type ResumeSequenceRes struct {
	api_v1.StandardRes
}

// --- Enrollment ---

type EnrollContactsReq struct {
	g.Meta        `path:"/sequence/enroll" method:"post" tags:"Sequence" summary:"Enroll contacts"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id"   v:"required|min:1" dc:"Sequence ID"`
	ContactIds    []int  `json:"contact_ids"   dc:"Contact IDs (empty = all from group/tags)"`
}

type EnrollContactsRes struct {
	api_v1.StandardRes
	Data struct {
		EnrolledCount int `json:"enrolled_count"`
		SkippedCount  int `json:"skipped_count"`
	} `json:"data"`
}

type ListEnrollmentsReq struct {
	g.Meta        `path:"/sequence/enrollments" method:"get" tags:"Sequence" summary:"List enrollments"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id"    v:"required|min:1" dc:"Sequence ID"`
	Page          int    `json:"page"           v:"required|min:1" dc:"Page number"`
	PageSize      int    `json:"page_size"      v:"required|min:1|max:100" dc:"Page size"`
	Status        int    `json:"status"         dc:"Filter by status (-1 = all)"`
}

type ListEnrollmentsRes struct {
	api_v1.StandardRes
	Data struct {
		Total int                  `json:"total"`
		List  []EnrollmentListItem `json:"list"`
	} `json:"data"`
}

type EnrollmentListItem struct {
	Id              int    `json:"id"`
	SequenceId      int    `json:"sequence_id"`
	ContactId       int    `json:"contact_id"`
	Email           string `json:"email"`
	CurrentStep     int    `json:"current_step"`
	Status          int    `json:"status"`
	EnrolledAt      int    `json:"enrolled_at"`
	LastEmailSentAt int    `json:"last_email_sent_at"`
	TotalEmailsSent int    `json:"total_emails_sent"`
	TotalOpens      int    `json:"total_opens"`
	TotalClicks     int    `json:"total_clicks"`
}

type RemoveEnrollmentReq struct {
	g.Meta        `path:"/sequence/remove_enrollment" method:"post" tags:"Sequence" summary:"Remove enrollment"`
	Authorization string `json:"authorization" in:"header"`
	EnrollmentId  int    `json:"enrollment_id"  v:"required|min:1" dc:"Enrollment ID"`
}

type RemoveEnrollmentRes struct {
	api_v1.StandardRes
}

// --- Test ---

type SendTestStepReq struct {
	g.Meta        `path:"/sequence/send_test" method:"post" tags:"Sequence" summary:"Send test email for a step"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id"   v:"required|min:1" dc:"Sequence ID"`
	StepId        int    `json:"step_id"       v:"required|min:1" dc:"Step ID"`
	TestEmail     string `json:"test_email"    v:"required|email" dc:"Test recipient email"`
}

type SendTestStepRes struct {
	api_v1.StandardRes
}
