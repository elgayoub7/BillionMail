package sequence

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
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
			Data(g.Map{"total_opens": g.Raw("total_opens + 1")}).
			Update()
	case "click":
		g.DB().Model("bm_sequence_enrollments").
			WhereIn("id", enrollmentIds).
			Where("email", email).
			Where("status", 0).
			Data(g.Map{"total_clicks": g.Raw("total_clicks + 1")}).
			Update()
	}
}
