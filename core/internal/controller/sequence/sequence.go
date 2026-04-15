package sequence

import (
	"billionmail-core/api/sequence/v1"
	"billionmail-core/internal/consts"
	"billionmail-core/internal/service/public"
	"billionmail-core/internal/service/sequence"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateSequence(ctx context.Context, req *v1.CreateSequenceReq) (res *v1.CreateSequenceRes, err error) {
	res = &v1.CreateSequenceRes{}

	id, err := sequence.CreateSequence(ctx, sequence.CreateSequenceArgs{
		Name:        req.Name,
		Description: req.Description,
		Addresser:   req.Addresser,
		FullName:    req.FullName,
		GroupId:     req.GroupId,
		TagIds:      req.TagIds,
		TagLogic:    req.TagLogic,
		TrackOpen:   req.TrackOpen,
		TrackClick:  req.TrackClick,
		Unsubscribe: req.Unsubscribe,
		Steps:       req.Steps,
	})
	if err != nil {
		res.SetError(err)
		return
	}

	_ = public.WriteLog(ctx, public.LogParams{
		Type: consts.LOGTYPE.Task,
		Log:  "Create sequence: " + req.Name,
		Data: req,
	})

	res.Data.Id = id
	res.SetSuccess("Sequence created successfully")
	return
}

func (c *ControllerV1) UpdateSequence(ctx context.Context, req *v1.UpdateSequenceReq) (res *v1.UpdateSequenceRes, err error) {
	res = &v1.UpdateSequenceRes{}

	if req.Id == 0 {
		res.Code = 400
		res.SetError(gerror.New("Sequence ID is required"))
		return
	}

	err = sequence.UpdateSequence(ctx, sequence.UpdateSequenceArgs{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Addresser:   req.Addresser,
		FullName:    req.FullName,
		GroupId:     req.GroupId,
		TagIds:      req.TagIds,
		TagLogic:    req.TagLogic,
		TrackOpen:   req.TrackOpen,
		TrackClick:  req.TrackClick,
		Unsubscribe: req.Unsubscribe,
		Steps:       req.Steps,
	})
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Sequence updated successfully")
	return
}

func (c *ControllerV1) DeleteSequence(ctx context.Context, req *v1.DeleteSequenceReq) (res *v1.DeleteSequenceRes, err error) {
	res = &v1.DeleteSequenceRes{}

	err = sequence.DeleteSequence(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Sequence deleted successfully")
	return
}

func (c *ControllerV1) ListSequences(ctx context.Context, req *v1.ListSequencesReq) (res *v1.ListSequencesRes, err error) {
	res = &v1.ListSequencesRes{}

	total, list, err := sequence.GetSequencesWithPage(ctx, req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess("")
	return
}

func (c *ControllerV1) FindSequence(ctx context.Context, req *v1.FindSequenceReq) (res *v1.FindSequenceRes, err error) {
	res = &v1.FindSequenceRes{}

	detail, err := sequence.GetSequenceDetail(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = detail
	res.SetSuccess("")
	return
}

func (c *ControllerV1) ActivateSequence(ctx context.Context, req *v1.ActivateSequenceReq) (res *v1.ActivateSequenceRes, err error) {
	res = &v1.ActivateSequenceRes{}

	err = sequence.ActivateSequence(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Sequence activated successfully")
	return
}

func (c *ControllerV1) PauseSequence(ctx context.Context, req *v1.PauseSequenceReq) (res *v1.PauseSequenceRes, err error) {
	res = &v1.PauseSequenceRes{}

	err = sequence.PauseSequence(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Sequence paused successfully")
	return
}

func (c *ControllerV1) ResumeSequence(ctx context.Context, req *v1.ResumeSequenceReq) (res *v1.ResumeSequenceRes, err error) {
	res = &v1.ResumeSequenceRes{}

	err = sequence.ResumeSequence(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Sequence resumed successfully")
	return
}

func (c *ControllerV1) EnrollContacts(ctx context.Context, req *v1.EnrollContactsReq) (res *v1.EnrollContactsRes, err error) {
	res = &v1.EnrollContactsRes{}

	enrolled, skipped, err := sequence.EnrollContacts(ctx, req.SequenceId, req.ContactIds)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data.EnrolledCount = enrolled
	res.Data.SkippedCount = skipped
	res.SetSuccess("Contacts enrolled successfully")
	return
}

func (c *ControllerV1) ListEnrollments(ctx context.Context, req *v1.ListEnrollmentsReq) (res *v1.ListEnrollmentsRes, err error) {
	res = &v1.ListEnrollmentsRes{}

	total, list, err := sequence.GetEnrollments(ctx, req.SequenceId, req.Page, req.PageSize, req.Status)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess("")
	return
}

func (c *ControllerV1) RemoveEnrollment(ctx context.Context, req *v1.RemoveEnrollmentReq) (res *v1.RemoveEnrollmentRes, err error) {
	res = &v1.RemoveEnrollmentRes{}

	err = sequence.RemoveEnrollment(ctx, req.EnrollmentId)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Enrollment removed successfully")
	return
}

func (c *ControllerV1) SendTestStep(ctx context.Context, req *v1.SendTestStepReq) (res *v1.SendTestStepRes, err error) {
	res = &v1.SendTestStepRes{}

	err = sequence.SendTestStep(ctx, req.SequenceId, req.StepId, req.TestEmail)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Test email sent successfully")
	return
}
