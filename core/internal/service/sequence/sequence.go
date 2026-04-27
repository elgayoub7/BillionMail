package sequence

import (
	"billionmail-core/api/sequence/v1"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type CreateSequenceArgs struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Addresser   string        `json:"addresser"`
	FullName    string        `json:"full_name"`
	GroupId     int           `json:"group_id"`
	TagIds      []int         `json:"tag_ids"`
	TagLogic    string        `json:"tag_logic"`
	TrackOpen   int           `json:"track_open"`
	TrackClick  int           `json:"track_click"`
	Unsubscribe int           `json:"unsubscribe"`
	SenderPool  string        `json:"sender_pool"`
	DailyLimit  int           `json:"daily_limit_per_sender"`
	SendDelay   int           `json:"send_delay"`
	ScheduleStart int         `json:"schedule_start_hour"`
	ScheduleEnd  int           `json:"schedule_end_hour"`
	ScheduleDays string        `json:"schedule_days"`
	Steps       []v1.StepInput `json:"steps"`
}

type UpdateSequenceArgs struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Addresser   string        `json:"addresser"`
	FullName    string        `json:"full_name"`
	GroupId     int           `json:"group_id"`
	TagIds      []int         `json:"tag_ids"`
	TagLogic    string        `json:"tag_logic"`
	TrackOpen   int           `json:"track_open"`
	TrackClick  int           `json:"track_click"`
	Unsubscribe int           `json:"unsubscribe"`
	SenderPool  string        `json:"sender_pool"`
	DailyLimit  int           `json:"daily_limit_per_sender"`
	SendDelay   int           `json:"send_delay"`
	ScheduleStart int         `json:"schedule_start_hour"`
	ScheduleEnd  int           `json:"schedule_end_hour"`
	ScheduleDays string        `json:"schedule_days"`
	Steps       []v1.StepInput `json:"steps"`
}

// CreateSequence creates a new sequence with its steps in a transaction
func CreateSequence(ctx context.Context, args CreateSequenceArgs) (int, error) {
	var sequenceId int64

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := time.Now().Unix()
		tagIdsJson, _ := json.Marshal(args.TagIds)

		// Default values
		// Auto-fill full_name from sender_pool if empty
		if args.FullName == "" && args.SenderPool != "" {
			type senderEntry struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			}
			var senders []senderEntry
			if json.Unmarshal([]byte(args.SenderPool), &senders) == nil && len(senders) > 0 {
				args.FullName = senders[0].Name
			}
		}
		trackOpen := args.TrackOpen
		if trackOpen == 0 {
			trackOpen = 1
		}
		trackClick := args.TrackClick
		if trackClick == 0 {
			trackClick = 1
		}
		unsubscribe := args.Unsubscribe
		if unsubscribe == 0 {
			unsubscribe = 1
		}
		tagLogic := args.TagLogic
		if tagLogic == "" {
			tagLogic = "AND"
		}

		result, err := tx.Ctx(ctx).Model("bm_sequences").Insert(g.Map{
			"name":         args.Name,
			"description":  args.Description,
			"status":       0, // draft
			"addresser":    args.Addresser,
			"full_name":    args.FullName,
			"group_id":     args.GroupId,
			"tag_ids":      string(tagIdsJson),
			"tag_logic":    tagLogic,
			"track_open":   trackOpen,
			"track_click":  trackClick,
			"unsubscribe":  unsubscribe,
			"sender_pool":  args.SenderPool,
			"daily_limit_per_sender": args.DailyLimit,
			"send_delay":   args.SendDelay,
			"schedule_start_hour": args.ScheduleStart,
			"schedule_end_hour":   args.ScheduleEnd,
			"schedule_days":       args.ScheduleDays,
			"create_time":  now,
			"update_time":  now,
		})
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		sequenceId = id

		// Insert steps
		for _, step := range args.Steps {
			_, err := tx.Ctx(ctx).Model("bm_sequence_steps").Insert(g.Map{
				"sequence_id":      sequenceId,
				"step_order":       step.StepOrder,
				"step_type":        step.StepType,
				"subject":          step.Subject,
				"template_id":      step.TemplateId,
				"wait_days":        step.WaitDays,
				"wait_hours":       step.WaitHours,
				"condition_type":   step.ConditionType,
				"condition_step_id": step.ConditionStepId,
				"on_true_go_to":    step.OnTrueGoTo,
				"on_false_go_to":   step.OnFalseGoTo,
				"create_time":      now,
				"update_time":      now,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return int(sequenceId), nil
}

// GetSequencesWithPage returns paginated list of sequences
func GetSequencesWithPage(ctx context.Context, page, pageSize int, keyword string, status int) (int, []v1.SequenceListItem, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	model := g.DB().Model("bm_sequences s").Safe()

	if keyword != "" {
		model = model.WhereLike("s.name", "%"+keyword+"%")
	}
	if status > -1 {
		model = model.Where("s.status", status)
	}

	total, err := model.Count()
	if err != nil {
		return 0, nil, err
	}

	type seqRow struct {
		Id             int    `json:"id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Status         int    `json:"status"`
		TotalEnrolled  int    `json:"total_enrolled"`
		TotalCompleted int    `json:"total_completed"`
		TotalBounced   int    `json:"total_bounced"`
		GroupId        int    `json:"group_id"`
		CreateTime     int    `json:"create_time"`
		UpdateTime     int    `json:"update_time"`
		StepCount      int    `json:"step_count"`
		GroupName      string `json:"group_name"`
	}

	var rows []seqRow
	err = model.Page(page, pageSize).
		Fields(`s.*, (SELECT COUNT(*) FROM bm_sequence_steps WHERE sequence_id = s.id) as step_count,
			(SELECT name FROM bm_contact_groups WHERE id = s.group_id) as group_name`).
		Order("s.create_time DESC").
		Scan(&rows)

	if err != nil {
		return 0, nil, err
	}

	list := make([]v1.SequenceListItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, v1.SequenceListItem{
			Id:             row.Id,
			Name:           row.Name,
			Description:    row.Description,
			Status:         row.Status,
			StepCount:      row.StepCount,
			TotalEnrolled:  row.TotalEnrolled,
			TotalCompleted: row.TotalCompleted,
			TotalBounced:   row.TotalBounced,
			GroupName:      row.GroupName,
			CreateTime:     row.CreateTime,
			UpdateTime:     row.UpdateTime,
		})
	}

	return total, list, nil
}

// GetSequenceDetail returns full sequence detail with steps
func GetSequenceDetail(ctx context.Context, id int) (*v1.SequenceDetail, error) {
	type seqRow struct {
		Id                int    `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		Status            int    `json:"status"`
		Addresser         string `json:"addresser"`
		FullName          string `json:"full_name"`
		GroupId           int    `json:"group_id"`
		TagIds            string `json:"tag_ids"`
		TagLogic          string `json:"tag_logic"`
		TrackOpen         int    `json:"track_open"`
		TrackClick        int    `json:"track_click"`
		Unsubscribe       int    `json:"unsubscribe"`
		TotalEnrolled     int    `json:"total_enrolled"`
		TotalCompleted    int    `json:"total_completed"`
		TotalUnsubscribed int    `json:"total_unsubscribed"`
		TotalBounced      int    `json:"total_bounced"`
		GroupName         string `json:"group_name"`
		CreateTime        int    `json:"create_time"`
		UpdateTime        int    `json:"update_time"`
	}

	var row seqRow
	err := g.DB().Model("bm_sequences s").
		Fields(`s.*, (SELECT name FROM bm_contact_groups WHERE id = s.group_id) as group_name`).
		Where("s.id", id).
		Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, gerror.New("Sequence not found")
	}

	var tagIds []int
	_ = json.Unmarshal([]byte(row.TagIds), &tagIds)

	type stepRow struct {
		Id              int    `json:"id"`
		StepOrder       int    `json:"step_order"`
		StepType        string `json:"step_type"`
		Subject         string `json:"subject"`
		TemplateId      int    `json:"template_id"`
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
	RepliedCount    int    `json:"replied_count"`
		TemplateName    string `json:"template_name"`
	}

	var steps []stepRow
	err = g.DB().Model("bm_sequence_steps st").
		LeftJoin("email_templates et", "et.id = st.template_id").
		Fields(`st.*, et.temp_name as template_name`).
		Where("st.sequence_id", id).
		Order("st.step_order ASC").
		Scan(&steps)
	if err != nil {
		g.Log().Warningf(ctx, "Failed to query steps for sequence %d: %v", id, err)
	}

	stepItems := make([]v1.SequenceStepItem, 0, len(steps))
	for _, s := range steps {
		stepItems = append(stepItems, v1.SequenceStepItem{
			Id:              s.Id,
			StepOrder:       s.StepOrder,
			StepType:        s.StepType,
			Subject:         s.Subject,
			TemplateId:      s.TemplateId,
			TemplateName:    s.TemplateName,
			WaitDays:        s.WaitDays,
			WaitHours:       s.WaitHours,
			ConditionType:   s.ConditionType,
			ConditionStepId: s.ConditionStepId,
			OnTrueGoTo:      s.OnTrueGoTo,
			OnFalseGoTo:     s.OnFalseGoTo,
			SentCount:       s.SentCount,
			OpenedCount:     s.OpenedCount,
			ClickedCount:    s.ClickedCount,
			BouncedCount:    s.BouncedCount,
			RepliedCount:    s.RepliedCount,
		})
	}

	return &v1.SequenceDetail{
		Id:                row.Id,
		Name:              row.Name,
		Description:       row.Description,
		Status:            row.Status,
		Addresser:         row.Addresser,
		FullName:          row.FullName,
		GroupId:           row.GroupId,
		TagIds:            tagIds,
		TagLogic:          row.TagLogic,
		TrackOpen:         row.TrackOpen,
		TrackClick:        row.TrackClick,
		Unsubscribe:       row.Unsubscribe,
		TotalEnrolled:     row.TotalEnrolled,
		TotalCompleted:    row.TotalCompleted,
		TotalUnsubscribed: row.TotalUnsubscribed,
		TotalBounced:      row.TotalBounced,
		GroupName:         row.GroupName,
		CreateTime:        row.CreateTime,
		UpdateTime:        row.UpdateTime,
		Steps:             stepItems,
	}, nil
}

// UpdateSequence updates a sequence and replaces its steps
func UpdateSequence(ctx context.Context, args UpdateSequenceArgs) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := time.Now().Unix()
		tagIdsJson, _ := json.Marshal(args.TagIds)

		data := g.Map{
			"update_time": now,
		}
		if args.Name != "" {
			data["name"] = args.Name
		}
		if args.Description != "" {
			data["description"] = args.Description
		}
		if args.Addresser != "" {
			data["addresser"] = args.Addresser
		}
		if args.FullName != "" {
			data["full_name"] = args.FullName
		}
		if args.GroupId > 0 {
			data["group_id"] = args.GroupId
		}
		if args.TagLogic != "" {
			data["tag_logic"] = args.TagLogic
		}
		data["tag_ids"] = string(tagIdsJson)
		data["sender_pool"] = args.SenderPool
		data["daily_limit_per_sender"] = args.DailyLimit
		data["send_delay"] = args.SendDelay
		data["schedule_start_hour"] = args.ScheduleStart
		data["schedule_end_hour"] = args.ScheduleEnd
		data["schedule_days"] = args.ScheduleDays

		_, err := tx.Ctx(ctx).Model("bm_sequences").Where("id", args.Id).Data(data).Update()
		if err != nil {
			return err
		}

		if len(args.Steps) > 0 {
			_, err = tx.Ctx(ctx).Model("bm_sequence_steps").Where("sequence_id", args.Id).Delete()
			if err != nil {
				return err
			}
			for _, step := range args.Steps {
				_, err := tx.Ctx(ctx).Model("bm_sequence_steps").Insert(g.Map{
					"sequence_id":      args.Id,
					"step_order":       step.StepOrder,
					"step_type":        step.StepType,
					"subject":          step.Subject,
					"template_id":      step.TemplateId,
					"wait_days":        step.WaitDays,
					"wait_hours":       step.WaitHours,
					"condition_type":   step.ConditionType,
					"condition_step_id": step.ConditionStepId,
					"on_true_go_to":    step.OnTrueGoTo,
					"on_false_go_to":   step.OnFalseGoTo,
					"create_time":      now,
					"update_time":      now,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteSequence deletes a sequence and all related data
func DeleteSequence(ctx context.Context, id int) error {
	// Get email_task_ids from sequence_email_tasks before deleting
	var taskIds []int
	g.DB().Model("bm_sequence_email_tasks").
		Where("sequence_id", id).
		Fields("DISTINCT email_task_id").
		Scan(&taskIds)

	// Delete recipient_info for these tasks
	if len(taskIds) > 0 {
		g.DB().Model("recipient_info").WhereIn("task_id", taskIds).Delete()
	}

	// Delete sequence_email_tasks
	g.DB().Model("bm_sequence_email_tasks").Where("sequence_id", id).Delete()

	// Delete sequence_enrollments
	g.DB().Model("bm_sequence_enrollments").Where("sequence_id", id).Delete()

	// Delete email_tasks created by this sequence
	if len(taskIds) > 0 {
		g.DB().Model("email_tasks").WhereIn("id", taskIds).Delete()
	}

	// Delete sequence_steps
	g.DB().Model("bm_sequence_steps").Where("sequence_id", id).Delete()

	// Delete the sequence itself
	_, err := g.DB().Model("bm_sequences").Where("id", id).Delete()
	return err
}

// ActivateSequence changes status from draft (0) to active (1)
func ActivateSequence(ctx context.Context, id int) error {
	var seq struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	}
	err := g.DB().Model("bm_sequences").Where("id", id).Scan(&seq)
	if err != nil {
		return err
	}
	if seq.Id == 0 {
		return gerror.New("Sequence not found")
	}
	if seq.Status != 0 && seq.Status != 2 {
		return gerror.New("Only draft or paused sequences can be activated")
	}

	_, err = g.DB().Model("bm_sequences").Where("id", id).Data(g.Map{
		"status":      1,
		"update_time": time.Now().Unix(),
	}).Update()
	return err
}

// PauseSequence changes status from active (1) to paused (2)
func PauseSequence(ctx context.Context, id int) error {
	_, err := g.DB().Model("bm_sequences").Where("id", id).Where("status", 1).Data(g.Map{
		"status":      2,
		"update_time": time.Now().Unix(),
	}).Update()
	if err != nil {
		return gerror.New("Failed to pause sequence")
	}
	return nil
}

// ResumeSequence changes status from paused (2) to active (1)
func ResumeSequence(ctx context.Context, id int) error {
	_, err := g.DB().Model("bm_sequences").Where("id", id).Where("status", 2).Data(g.Map{
		"status":      1,
		"update_time": time.Now().Unix(),
	}).Update()
	if err != nil {
		return gerror.New("Failed to resume sequence")
	}
	return nil
}

// EnrollContacts enrolls contacts into a sequence. If contactIds is empty, enrolls all from the sequence's group.
func EnrollContacts(ctx context.Context, sequenceId int, contactIds []int) (enrolled int, skipped int, err error) {
	var seq struct {
		Id      int    `json:"id"`
		Status  int    `json:"status"`
		GroupId int    `json:"group_id"`
		TagIds  string `json:"tag_ids"`
		TagLogic string `json:"tag_logic"`
	}
	err = g.DB().Model("bm_sequences").Where("id", sequenceId).Scan(&seq)
	if err != nil {
		return 0, 0, err
	}
	if seq.Id == 0 {
		return 0, 0, gerror.New("Sequence not found")
	}
	if seq.Status != 1 {
		return 0, 0, gerror.New("Sequence must be active to enroll contacts")
	}

	now := time.Now().Unix()

	var contacts []struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	model := g.DB().Model("bm_contacts").Where("active", 1).Where("status", 1)
	if len(contactIds) > 0 {
		model = model.WhereIn("id", contactIds)
	} else if seq.GroupId > 0 {
		model = model.Where("group_id", seq.GroupId)
	}
	// group_id=0 means all groups - no filter applied
	err = model.Fields("id, email").Scan(&contacts)
	if err != nil {
		return 0, 0, err
	}

	const batchSize = 500
	for i := 0; i < len(contacts); i += batchSize {
		end := i + batchSize
		if end > len(contacts) {
			end = len(contacts)
		}
		batch := contacts[i:end]

		for _, c := range batch {
			count, _ := g.DB().Model("bm_sequence_enrollments").
				Where("sequence_id", sequenceId).
				Where("email", c.Email).
				Count()
			if count > 0 {
				skipped++
				continue
			}

			_, err = g.DB().Model("bm_sequence_enrollments").Insert(g.Map{
				"sequence_id":             sequenceId,
				"contact_id":              c.Id,
				"email":                   c.Email,
				"group_id":                seq.GroupId,
				"current_step":            1,
				"status":                  0,
				"enrolled_at":             now,
				"current_step_entered_at": now,
			})
			if err != nil {
				g.Log().Warning(ctx, "Failed to enroll contact:", c.Email, err)
				skipped++
				continue
			}
			enrolled++
		}
	}

	// Update sequence total_enrolled to reflect actual count
	g.DB().Model("bm_sequences").Where("id", sequenceId).Data(g.Map{
		"total_enrolled": gdb.Raw(fmt.Sprintf("(SELECT COUNT(*) FROM bm_sequence_enrollments WHERE sequence_id = %d)", sequenceId)),
		"update_time":    now,
	}).Update()

	return enrolled, skipped, nil
}

// GetEnrollments returns paginated enrollments for a sequence
func GetEnrollments(ctx context.Context, sequenceId int, page, pageSize int, status int) (int, []v1.EnrollmentListItem, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	model := g.DB().Model("bm_sequence_enrollments").Where("sequence_id", sequenceId)
	if status > -1 {
		model = model.Where("status", status)
	}

	total, err := model.Count()
	if err != nil {
		return 0, nil, err
	}

	var rows []v1.EnrollmentListItem
	err = model.Page(page, pageSize).
		Fields("id, sequence_id, contact_id, email, current_step, status, enrolled_at, last_email_sent_at, total_emails_sent, total_opens, total_clicks, total_replies").
		Order("enrolled_at DESC").
		Scan(&rows)

	return total, rows, err
}

// RemoveEnrollment removes a contact from a sequence
func RemoveEnrollment(ctx context.Context, enrollmentId int) error {
	// Get enrollment info
	var enrollment struct {
		Id         int    `json:"id"`
		SequenceId int    `json:"sequence_id"`
		Email      string `json:"email"`
	}
	err := g.DB().Model("bm_sequence_enrollments").Where("id", enrollmentId).Scan(&enrollment)
	if err != nil {
		return err
	}
	if enrollment.Id == 0 {
		return gerror.New("Enrollment not found")
	}

	_, err = g.DB().Model("bm_sequence_enrollments").Where("id", enrollmentId).Data(g.Map{
		"status":       3, // exited
		"completed_at": time.Now().Unix(),
	}).Update()
	if err != nil {
		return err
	}

	g.DB().Model("bm_sequences").Where("id", enrollment.SequenceId).Data(g.Map{
		"total_unsubscribed": gdb.Raw(fmt.Sprintf(
			"(SELECT COUNT(*) FROM bm_sequence_enrollments WHERE sequence_id = %d AND status = 3)",
			enrollment.SequenceId)),
		"update_time": time.Now().Unix(),
	}).Update()

	return nil
}

// SendTestStep sends a test email for a specific step
func SendTestStep(ctx context.Context, sequenceId, stepId int, testEmail string) error {
	var step struct {
		Id         int    `json:"id"`
		SequenceId int    `json:"sequence_id"`
		Subject    string `json:"subject"`
		TemplateId int    `json:"template_id"`
	}
	err := g.DB().Model("bm_sequence_steps").
		Where("id", stepId).Where("sequence_id", sequenceId).Scan(&step)
	if err != nil {
		return err
	}
	if step.Id == 0 {
		return gerror.New("Step not found")
	}

	var seq struct {
		Addresser string `json:"addresser"`
		FullName  string `json:"full_name"`
	}
	err = g.DB().Model("bm_sequences").Where("id", sequenceId).Scan(&seq)
	if err != nil {
		return err
	}

	// Import the test recipient using the existing batch_mail infrastructure
	// For now, create a minimal email_task for testing
	now := time.Now().Unix()
	taskName := fmt.Sprintf("seq_test_%d_step_%d_%d", sequenceId, stepId, now)

	result, err := g.DB().Model("email_tasks").Insert(g.Map{
		"task_name":    taskName,
		"addresser":    seq.Addresser,
		"subject":      step.Subject,
		"full_name":    seq.FullName,
		"template_id":  step.TemplateId,
		"recipient_count": 0,
		"task_process": 0,
		"pause":        0,
		"threads":      1,
		"track_open":   1,
		"track_click":  1,
		"unsubscribe":  1,
		"start_time":   now,
		"create_time":  now,
		"update_time":  now,
		"active":       1,
		"add_type":     3, // sequence test
		"group_id":     0,
	})
	if err != nil {
		return gerror.New("Failed to create test task")
	}

	taskId, _ := result.LastInsertId()

	// Insert the test recipient
	_, err = g.DB().Model("recipient_info").Insert(g.Map{
		"task_id":     taskId,
		"recipient":   testEmail,
		"is_sent":     0,
		"sent_time":   0,
		"message_id":  "",
		"create_time": now,
	})
	if err != nil {
		return gerror.New("Failed to add test recipient")
	}

	g.Log().Info(ctx, "Test email task created for sequence step:", sequenceId, stepId, "task:", taskId)
	return nil
}

// Helper: check if a sequence exists and is active
func GetActiveSequence(ctx context.Context, id int) (bool, error) {
	count, err := g.DB().Model("bm_sequences").Where("id", id).Where("status", 1).Count()
	return count > 0, err
}

// Helper: get steps for a sequence ordered by step_order
func GetSequenceSteps(ctx context.Context, sequenceId int) ([]v1.SequenceStepItem, error) {
	type stepRow struct {
		Id              int    `json:"id"`
		StepOrder       int    `json:"step_order"`
		StepType        string `json:"step_type"`
		Subject         string `json:"subject"`
		TemplateId      int    `json:"template_id"`
		WaitDays        int    `json:"wait_days"`
		WaitHours       int    `json:"wait_hours"`
		ConditionType   string `json:"condition_type"`
		ConditionStepId int    `json:"condition_step_id"`
		OnTrueGoTo      int    `json:"on_true_go_to"`
		OnFalseGoTo     int    `json:"on_false_go_to"`
	}

	var rows []stepRow
	err := g.DB().Model("bm_sequence_steps").
		Where("sequence_id", sequenceId).
		Order("step_order ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	items := make([]v1.SequenceStepItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.SequenceStepItem{
			Id:              r.Id,
			StepOrder:       r.StepOrder,
			StepType:        r.StepType,
			Subject:         r.Subject,
			TemplateId:      r.TemplateId,
			WaitDays:        r.WaitDays,
			WaitHours:       r.WaitHours,
			ConditionType:   r.ConditionType,
			ConditionStepId: r.ConditionStepId,
			OnTrueGoTo:      r.OnTrueGoTo,
			OnFalseGoTo:     r.OnFalseGoTo,
		})
	}
	return items, nil
}

// Helper: ensure template ID is always passed as proper type
func ensureInt(val interface{}) int {
	if val == nil {
		return 0
	}
	return gconv.Int(val)
}
