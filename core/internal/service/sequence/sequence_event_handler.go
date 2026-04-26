package sequence

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/service/scoring"
)

func RecordEvent(ctx context.Context, eventType, email string, campaignId int, messageId string) {
	var seqTasks []struct {
		SequenceId int `json:"sequence_id"`
		StepId     int `json:"step_id"`
	}
	g.DB().Model("bm_sequence_email_tasks").
		Where("email_task_id", campaignId).
		Fields("DISTINCT sequence_id, step_id").
		Scan(&seqTasks)

	if len(seqTasks) == 0 {
		return
	}

	stepIds := make([]int, len(seqTasks))
	for i, st := range seqTasks {
		stepIds[i] = st.StepId
	}

	enrollmentIdRows, err := g.DB().Model("bm_sequence_email_tasks").
		Fields("DISTINCT enrollment_id").
		Where("email_task_id", campaignId).
		Array()
	if err != nil || len(enrollmentIdRows) == 0 {
		return
	}
	enrollmentIds := make([]int, len(enrollmentIdRows))
	for i, v := range enrollmentIdRows {
		enrollmentIds[i] = v.Int()
	}

	switch eventType {
	case "open":
		g.DB().Model("bm_sequence_enrollments").
			WhereIn("id", enrollmentIds).
			Where("email", email).
			Where("status", 0).
			Data(g.Map{"total_opens": gdb.Raw("total_opens + 1")}).
			Update()
		g.DB().Model("bm_sequence_steps").
			WhereIn("id", stepIds).
			Data(g.Map{"opened_count": gdb.Raw("opened_count + 1")}).
			Update()
	case "click":
		g.DB().Model("bm_sequence_enrollments").
			WhereIn("id", enrollmentIds).
			Where("email", email).
			Where("status", 0).
			Data(g.Map{"total_clicks": gdb.Raw("total_clicks + 1")}).
			Update()
		g.DB().Model("bm_sequence_steps").
			WhereIn("id", stepIds).
			Data(g.Map{"clicked_count": gdb.Raw("clicked_count + 1")}).
			Update()
	}

	go scoring.Scoring().RecordEngagement(context.Background(), email, eventType)

	for _, st := range seqTasks {
		RecordAbTestEvent(ctx, st.SequenceId, st.StepId, email, eventType)
	}
}
