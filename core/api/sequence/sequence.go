package sequence

import (
	"context"

	"billionmail-core/api/sequence/v1"
)

type ISequenceV1 interface {
	CreateSequence(ctx context.Context, req *v1.CreateSequenceReq) (res *v1.CreateSequenceRes, err error)
	UpdateSequence(ctx context.Context, req *v1.UpdateSequenceReq) (res *v1.UpdateSequenceRes, err error)
	DeleteSequence(ctx context.Context, req *v1.DeleteSequenceReq) (res *v1.DeleteSequenceRes, err error)
	ListSequences(ctx context.Context, req *v1.ListSequencesReq) (res *v1.ListSequencesRes, err error)
	FindSequence(ctx context.Context, req *v1.FindSequenceReq) (res *v1.FindSequenceRes, err error)
	ActivateSequence(ctx context.Context, req *v1.ActivateSequenceReq) (res *v1.ActivateSequenceRes, err error)
	PauseSequence(ctx context.Context, req *v1.PauseSequenceReq) (res *v1.PauseSequenceRes, err error)
	ResumeSequence(ctx context.Context, req *v1.ResumeSequenceReq) (res *v1.ResumeSequenceRes, err error)
	EnrollContacts(ctx context.Context, req *v1.EnrollContactsReq) (res *v1.EnrollContactsRes, err error)
	ListEnrollments(ctx context.Context, req *v1.ListEnrollmentsReq) (res *v1.ListEnrollmentsRes, err error)
	RemoveEnrollment(ctx context.Context, req *v1.RemoveEnrollmentReq) (res *v1.RemoveEnrollmentRes, err error)
	SendTestStep(ctx context.Context, req *v1.SendTestStepReq) (res *v1.SendTestStepRes, err error)
}
