package entity

// AbTest represents an A/B test for a sequence step
type AbTest struct {
	Id                  int     `json:"id"                     dc:"ID"`
	SequenceId          int     `json:"sequence_id"            dc:"Sequence ID"`
	StepId              int     `json:"step_id"                dc:"Step ID"`
	Name                string  `json:"name"                   dc:"Test Name"`
	Status              int     `json:"status"                 dc:"Status: 0=draft, 1=running, 2=completed"`
	WinnerVariant       int     `json:"winner_variant"         dc:"Winner variant: -1=none, 0=A, 1=B"`
	WinnerCriteria      string  `json:"winner_criteria"        dc:"Winner criteria: open_rate, click_rate, reply_rate"`
	SplitRatio          float64 `json:"split_ratio"            dc:"Split ratio for variant A (0.0-1.0)"`
	ConfidenceThreshold float64 `json:"confidence_threshold"   dc:"Statistical confidence threshold"`
	CreatedAt           int     `json:"created_at"             dc:"Created At"`
	CompletedAt         int     `json:"completed_at"           dc:"Completed At"`
}

// AbTestVariant represents one variant (A or B) of an A/B test
type AbTestVariant struct {
	Id          int    `json:"id"            dc:"ID"`
	AbTestId    int    `json:"ab_test_id"    dc:"AB Test ID"`
	Variant     int    `json:"variant"       dc:"Variant: 0=A, 1=B"`
	Subject     string `json:"subject"       dc:"Email Subject"`
	BodyHtml    string `json:"body_html"     dc:"HTML Body"`
	BodyText    string `json:"body_text"     dc:"Plain Text Body"`
	SentCount   int    `json:"sent_count"    dc:"Sent Count"`
	OpenCount   int    `json:"open_count"    dc:"Open Count"`
	ClickCount  int    `json:"click_count"   dc:"Click Count"`
	ReplyCount  int    `json:"reply_count"   dc:"Reply Count"`
	BounceCount int    `json:"bounce_count"  dc:"Bounce Count"`
}

// AbTestAssignment tracks which variant each enrollment was assigned to
type AbTestAssignment struct {
	Id           int `json:"id"             dc:"ID"`
	AbTestId     int `json:"ab_test_id"     dc:"AB Test ID"`
	EnrollmentId int `json:"enrollment_id"  dc:"Enrollment ID"`
	Variant      int `json:"variant"        dc:"Assigned Variant: 0=A, 1=B"`
	AssignedAt   int `json:"assigned_at"    dc:"Assigned At"`
}

// AbTestResult represents the comparison results for an A/B test
type AbTestResult struct {
	TestId      int             `json:"test_id"`
	TestName    string          `json:"test_name"`
	Status      int             `json:"status"`
	Winner      int             `json:"winner"`
	VariantA    *VariantStats   `json:"variant_a"`
	VariantB    *VariantStats   `json:"variant_b"`
	IsSignificant bool          `json:"is_significant"`
}

// VariantStats holds computed statistics for one variant
type VariantStats struct {
	VariantId   int     `json:"variant_id"`
	Label       string  `json:"label"`
	SentCount   int     `json:"sent_count"`
	OpenRate    float64 `json:"open_rate"`
	ClickRate   float64 `json:"click_rate"`
	ReplyRate   float64 `json:"reply_rate"`
	BounceRate  float64 `json:"bounce_rate"`
}
