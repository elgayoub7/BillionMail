package entity

type LeadScore struct {
	Id               int    `json:"id"                  dc:"ID"`
	ContactEmail     string `json:"contact_email"       dc:"Contact email"`
	Score            int    `json:"score"               dc:"Engagement score"`
	EngagementLevel  string `json:"engagement_level"    dc:"cold, warm, hot, converted"`
	TotalOpens       int    `json:"total_opens"         dc:"Total opens"`
	TotalClicks      int    `json:"total_clicks"        dc:"Total clicks"`
	TotalReplies     int    `json:"total_replies"       dc:"Total replies"`
	TotalBounces     int    `json:"total_bounces"       dc:"Total bounces"`
	LastEngagementAt int    `json:"last_engagement_at"  dc:"Last engagement timestamp"`
	LastScoredAt     int    `json:"last_scored_at"      dc:"Last scoring timestamp"`
}

type BounceRecord struct {
	Id           int    `json:"id"            dc:"ID"`
	Email        string `json:"email"         dc:"Bounced email"`
	BounceType   string `json:"bounce_type"   dc:"hard, soft, complaint"`
	Description  string `json:"description"   dc:"Bounce description"`
	SequenceId   int    `json:"sequence_id"   dc:"Source sequence ID"`
	MessageId    string `json:"message_id"    dc:"Message ID"`
	ActionTaken  string `json:"action_taken"  dc:"Action taken"`
	CreatedAt    int    `json:"created_at"    dc:"Created At"`
}

type ScoringStats struct {
	ColdCount     int `json:"cold_count"`
	WarmCount     int `json:"warm_count"`
	HotCount      int `json:"hot_count"`
	ConvertedCount int `json:"converted_count"`
	AverageScore  float64 `json:"average_score"`
	TotalLeads    int `json:"total_leads"`
}
