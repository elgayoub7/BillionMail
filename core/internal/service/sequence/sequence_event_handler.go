package sequence

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/service/scoring"
)

func RecordEvent(ctx context.Context, eventType, email string, campaignId int, messageId string) {
	var enrollmentIds []int
	err := g.DB().Model("bm_sequence_email_tasks").
		Fields("DISTINCT enrollment_id").
		Where("email_task_id", campaignId).
		Scan(&enrollmentIds)
	if err != nil || len(enrollmentIds) == 0 {
		return
	}

	switch eventType {
	case "open":
		g.DB().Model("bm_sequence_enrollments").
			WhereIn("id", enrollmentIds).
			Where("email", email).
			Where("status", 0).
			Data(g.Map{"total_opens": gdb.Raw("total_opens + 1")}).
			Update()
	case "click":
		g.DB().Model("bm_sequence_enrollments").
			WhereIn("id", enrollmentIds).
			Where("email", email).
			Where("status", 0).
			Data(g.Map{"total_clicks": gdb.Raw("total_clicks + 1")}).
			Update()
	}

	// Lead scoring engagement tracking
	go scoring.Scoring().RecordEngagement(context.Background(), email, eventType)

	// AB test event tracking
	var seqTasks []struct {
		SequenceId int `json:"sequence_id"`
		StepId     int `json:"step_id"`
	}
	g.DB().Model("bm_sequence_email_tasks").
		Where("email_task_id", campaignId).
		Fields("DISTINCT sequence_id, step_id").
		Scan(&seqTasks)

	for _, st := range seqTasks {
		RecordAbTestEvent(ctx, st.SequenceId, st.StepId, email, eventType)
	}
}
