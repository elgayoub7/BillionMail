package sequence

import (
	"billionmail-core/api/sequence/v1"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	engineRunning int32
)

func ProcessSequenceEnrollments(ctx context.Context) {
	if !atomic.CompareAndSwapInt32(&engineRunning, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&engineRunning, 0)

	advanceReadyEnrollments(ctx)
	evaluateConditionSteps(ctx)
	createBatchStepEmails(ctx)
}

func advanceReadyEnrollments(ctx context.Context) {
	now := time.Now().Unix()

	type enrollmentRow struct {
		Id                   int `json:"id"`
		SequenceId           int `json:"sequence_id"`
		CurrentStep          int `json:"current_step"`
		CurrentStepEnteredAt int `json:"current_step_entered_at"`
	}

	var enrollments []enrollmentRow
	err := g.DB().Model("bm_sequence_enrollments e").
		Join("bm_sequences s", "s.id = e.sequence_id").
		Join("bm_sequence_steps st", "st.sequence_id = s.id AND st.step_order = e.current_step").
		Fields("e.id, e.sequence_id, e.current_step, e.current_step_entered_at").
		Where("e.status", 0).
		Where("s.status", 1).
		Where("st.step_type", "wait").
		Where(fmt.Sprintf("e.current_step_entered_at + (st.wait_days * 86400) + (st.wait_hours * 3600) <= %d", now)).
		Scan(&enrollments)

	if err != nil {
		g.Log().Error(ctx, "advanceReadyEnrollments query failed:", err)
		return
	}

	for _, enr := range enrollments {
		steps, err := GetSequenceSteps(ctx, enr.SequenceId)
		if err != nil {
			continue
		}

		nextOrder := findNextStepOrder(steps, enr.CurrentStep)
		if nextOrder == 0 {
			markEnrollmentCompleted(ctx, enr.Id, now)
			continue
		}

		_, err = g.DB().Model("bm_sequence_enrollments").Where("id", enr.Id).Data(g.Map{
			"current_step":            nextOrder,
			"current_step_entered_at": now,
		}).Update()
		if err != nil {
			g.Log().Warning(ctx, "Failed to advance enrollment", enr.Id, "to step", nextOrder)
		}
	}
}

func evaluateConditionSteps(ctx context.Context) {
	now := time.Now().Unix()

	type enrollmentRow struct {
		Id                   int    `json:"id"`
		SequenceId           int    `json:"sequence_id"`
		Email                string `json:"email"`
		CurrentStep          int    `json:"current_step"`
		CurrentStepEnteredAt int    `json:"current_step_entered_at"`
	}

	var enrollments []enrollmentRow
	err := g.DB().Model("bm_sequence_enrollments e").
		Join("bm_sequences s", "s.id = e.sequence_id").
		Join("bm_sequence_steps st", "st.sequence_id = s.id AND st.step_order = e.current_step").
		Fields("e.id, e.sequence_id, e.email, e.current_step, e.current_step_entered_at").
		Where("e.status", 0).
		Where("s.status", 1).
		Where("st.step_type", "condition").
		Scan(&enrollments)

	if err != nil {
		g.Log().Error(ctx, "evaluateConditionSteps query failed:", err)
		return
	}

	for _, enr := range enrollments {
		steps, err := GetSequenceSteps(ctx, enr.SequenceId)
		if err != nil {
			continue
		}

		currentStep := findStepByOrder(steps, enr.CurrentStep)
		if currentStep == nil {
			continue
		}

		result, err := checkCondition(ctx, currentStep, enr.Email, enr.SequenceId)
		if err != nil {
			g.Log().Warning(ctx, "Condition check failed for enrollment", enr.Id, err)
			continue
		}

		var nextOrder int
		if result {
			nextOrder = currentStep.OnTrueGoTo
		} else {
			nextOrder = currentStep.OnFalseGoTo
		}

		if nextOrder == 0 {
			nextOrder = findNextStepOrder(steps, enr.CurrentStep)
		}
		if nextOrder == 0 {
			markEnrollmentCompleted(ctx, enr.Id, now)
			continue
		}

		_, err = g.DB().Model("bm_sequence_enrollments").Where("id", enr.Id).Data(g.Map{
			"current_step":            nextOrder,
			"current_step_entered_at": now,
		}).Update()
		if err != nil {
			g.Log().Warning(ctx, "Failed to move enrollment", enr.Id, "after condition")
		}
	}
}

func checkCondition(ctx context.Context, step *v1.SequenceStepItem, email string, sequenceId int) (bool, error) {
	emailTaskIds, err := getEmailTaskIdsForEnrollmentStep(ctx, sequenceId, step.ConditionStepId, email)
	if err != nil || len(emailTaskIds) == 0 {
		return false, nil
	}

	switch step.ConditionType {
	case "opened":
		count, _ := g.DB().Model("mailstat_opened").
			WhereIn("campaign_id", emailTaskIds).
			Where("recipient", email).
			Count()
		return count > 0, nil

	case "not_opened":
		count, _ := g.DB().Model("mailstat_opened").
			WhereIn("campaign_id", emailTaskIds).
			Where("recipient", email).
			Count()
		return count == 0, nil

	case "clicked":
		count, _ := g.DB().Model("mailstat_clicked").
			WhereIn("campaign_id", emailTaskIds).
			Where("recipient", email).
			Count()
		return count > 0, nil

	case "not_clicked":
		count, _ := g.DB().Model("mailstat_clicked").
			WhereIn("campaign_id", emailTaskIds).
			Where("recipient", email).
			Count()
		return count == 0, nil

	case "bounced":
		count, _ := g.DB().Model("mailstat_send_mails").
			WhereIn("campaign_id", emailTaskIds).
			Where("status", "bounced").
			Count()
		return count > 0, nil
	}

	return false, nil
}

func getEmailTaskIdsForEnrollmentStep(ctx context.Context, sequenceId, stepOrder int, email string) ([]int, error) {
	var taskIds []int
	err := g.DB().Model("bm_sequence_email_tasks").
		Fields("DISTINCT email_task_id").
		Join("bm_sequence_enrollments e", "e.id = bm_sequence_email_tasks.enrollment_id").
		Join("bm_sequence_steps st", "st.id = bm_sequence_email_tasks.step_id").
		Where("bm_sequence_email_tasks.sequence_id", sequenceId).
		Where("st.step_order", stepOrder).
		Where("e.email", email).
		Where("bm_sequence_email_tasks.status", 1).
		Scan(&taskIds)
	return taskIds, err
}

func createBatchStepEmails(ctx context.Context) {
	type pendingGroup struct {
		SequenceId int    `json:"sequence_id"`
		StepId     int    `json:"step_id"`
		StepOrder  int    `json:"step_order"`
		Subject    string `json:"subject"`
		TemplateId int    `json:"template_id"`
		Addresser  string `json:"addresser"`
		FullName   string `json:"full_name"`
		TrackOpen  int    `json:"track_open"`
		TrackClick int    `json:"track_click"`
		Unsub      int    `json:"unsubscribe"`
	}

	var groups []pendingGroup
	err := g.DB().Model("bm_sequence_enrollments e").
		Join("bm_sequences s", "s.id = e.sequence_id").
		Join("bm_sequence_steps st", "st.sequence_id = s.id AND st.step_order = e.current_step").
		LeftJoin("bm_sequence_email_tasks set", "set.enrollment_id = e.id AND set.step_id = st.id AND set.status = 0").
		Fields(`DISTINCT s.id as sequence_id, st.id as step_id, st.step_order, st.subject, st.template_id,
			s.addresser, s.full_name, s.track_open, s.track_click, s.unsubscribe as unsub`).
		Where("e.status", 0).
		Where("s.status", 1).
		Where("st.step_type", "email").
		Where("set.id IS NULL").
		Scan(&groups)

	if err != nil {
		g.Log().Error(ctx, "createBatchStepEmails query failed:", err)
		return
	}

	now := time.Now().Unix()

	for _, grp := range groups {
		taskName := fmt.Sprintf("seq_%d_step_%d_%d", grp.SequenceId, grp.StepOrder, now)

		result, err := g.DB().Model("email_tasks").Insert(g.Map{
			"task_name":    taskName,
			"addresser":    grp.Addresser,
			"subject":      grp.Subject,
			"full_name":    grp.FullName,
			"template_id":  grp.TemplateId,
			"task_process": 0,
			"pause":        0,
			"threads":      5,
			"track_open":   grp.TrackOpen,
			"track_click":  grp.TrackClick,
			"unsubscribe":  grp.Unsub,
			"start_time":   now,
			"create_time":  now,
			"update_time":  now,
			"active":       1,
			"add_type":     3,
			"group_id":     0,
		})
		if err != nil {
			g.Log().Error(ctx, "Failed to create email_task for sequence step:", grp.SequenceId, grp.StepOrder, err)
			continue
		}

		taskId, _ := result.LastInsertId()

		var enrollments []struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}
		g.DB().Model("bm_sequence_enrollments").
			Where("sequence_id", grp.SequenceId).
			Where("current_step", grp.StepOrder).
			Where("status", 0).
			Fields("id, email").
			Scan(&enrollments)

		for _, enr := range enrollments {
			_, err := g.DB().Model("recipient_info").Insert(g.Map{
				"task_id":     taskId,
				"recipient":   enr.Email,
				"is_sent":     0,
				"sent_time":   0,
				"message_id":  "",
				"create_time": now,
			})
			if err != nil {
				g.Log().Warning(ctx, "Failed to create recipient_info for enrollment", enr.Id, err)
				continue
			}

			_, err = g.DB().Model("bm_sequence_email_tasks").Insert(g.Map{
				"sequence_id":   grp.SequenceId,
				"enrollment_id": enr.Id,
				"step_id":       grp.StepId,
				"email_task_id": taskId,
				"contact_email": enr.Email,
				"status":        0,
			})
			if err != nil {
				g.Log().Warning(ctx, "Failed to create sequence_email_task for enrollment", enr.Id, err)
				continue
			}

			g.DB().Model("bm_sequence_enrollments").Where("id", enr.Id).Data(g.Map{
				"last_email_sent_at": now,
			}).Update()
		}

		g.DB().Model("bm_sequence_steps").Where("id", grp.StepId).Data(g.Map{
			"sent_count": g.Raw("sent_count + 1"),
		}).Update()
	}
}

func UpdateEnrollmentsFromCompletedTasks(ctx context.Context) {
	now := time.Now().Unix()

	type completedTask struct {
		EmailTaskId int `json:"email_task_id"`
	}

	var tasks []completedTask
	err := g.DB().Model("bm_sequence_email_tasks").
		Fields("DISTINCT email_task_id").
		Join("email_tasks et", "et.id = bm_sequence_email_tasks.email_task_id").
		Where("et.task_process", 2).
		Where("bm_sequence_email_tasks.status", 0).
		Limit(100).
		Scan(&tasks)

	if err != nil {
		return
	}

	for _, t := range tasks {
		_, err := g.DB().Model("bm_sequence_email_tasks").
			Where("email_task_id", t.EmailTaskId).
			Where("status", 0).
			Data(g.Map{
				"status":  1,
				"sent_at": now,
			}).Update()
		if err != nil {
			continue
		}

		var enrollmentIds []int
		g.DB().Model("bm_sequence_email_tasks").
			Where("email_task_id", t.EmailTaskId).
			Fields("DISTINCT enrollment_id").
			Scan(&enrollmentIds)

		for _, enrId := range enrollmentIds {
			g.DB().Model("bm_sequence_enrollments").Where("id", enrId).Data(g.Map{
				"total_emails_sent": g.Raw("total_emails_sent + 1"),
			}).Update()
		}

		g.DB().Raw(fmt.Sprintf(`
			UPDATE bm_sequence_steps SET
				sent_count = sent_count + (SELECT COUNT(*) FROM bm_sequence_email_tasks WHERE step_id = bm_sequence_steps.id AND email_task_id = %d),
				update_time = %d
			WHERE id IN (SELECT DISTINCT step_id FROM bm_sequence_email_tasks WHERE email_task_id = %d)
		`, t.EmailTaskId, now, t.EmailTaskId)).Exec()

		// Auto-exit enrollments whose contacts bounced
		g.DB().Raw(fmt.Sprintf(`
			UPDATE bm_sequence_enrollments SET status = 3, completed_at = %d
			WHERE id IN (
				SELECT DISTINCT enrollment_id FROM bm_sequence_email_tasks
				WHERE email_task_id = %d
			)
			AND email IN (
				SELECT DISTINCT recipient FROM mailstat_send_mails
				WHERE campaign_id = %d AND status = 'bounced'
			)
		`, now, t.EmailTaskId, t.EmailTaskId)).Exec()
	}
}

func markEnrollmentCompleted(ctx context.Context, enrollmentId int, now int64) {
	_, err := g.DB().Model("bm_sequence_enrollments").Where("id", enrollmentId).Data(g.Map{
		"status":       1,
		"completed_at": now,
	}).Update()
	if err != nil {
		g.Log().Warning(ctx, "Failed to mark enrollment completed:", enrollmentId, err)
		return
	}

	var seqId int
	g.DB().Model("bm_sequence_enrollments").Where("id", enrollmentId).Fields("sequence_id").Scan(&seqId)
	if seqId > 0 {
		g.DB().Model("bm_sequences").Where("id", seqId).Data(g.Map{
			"total_completed": g.Raw("total_completed + 1"),
			"update_time":    now,
		}).Update()
	}
}

func findNextStepOrder(steps []v1.SequenceStepItem, currentOrder int) int {
	for _, s := range steps {
		if s.StepOrder > currentOrder {
			return s.StepOrder
		}
	}
	return 0
}

func findStepByOrder(steps []v1.SequenceStepItem, order int) *v1.SequenceStepItem {
	for i := range steps {
		if steps[i].StepOrder == order {
			return &steps[i]
		}
	}
	return nil
}
