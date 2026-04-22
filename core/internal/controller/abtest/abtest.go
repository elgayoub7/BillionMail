package abtest

import (
	"context"

	v1 "billionmail-core/api/abtest/v1"
	"billionmail-core/internal/model/entity"
	"billionmail-core/internal/service/abtest"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateAbTest(ctx context.Context, req *v1.CreateAbTestReq) (res *v1.CreateAbTestRes, err error) {
	res = &v1.CreateAbTestRes{}

	id, err := abtest.AbTest().CreateAbTest(ctx, req.SequenceId, req.StepId, req.Name, req.SplitRatio, req.WinnerCriteria,
		entity.AbTestVariant{
			Subject:  req.VariantA.Subject,
			BodyHtml: req.VariantA.BodyHtml,
			BodyText: req.VariantA.BodyText,
		},
		entity.AbTestVariant{
			Subject:  req.VariantB.Subject,
			BodyHtml: req.VariantB.BodyHtml,
			BodyText: req.VariantB.BodyText,
		},
	)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data.Id = id
	res.SetSuccess("AB test created successfully")
	return
}

func (c *ControllerV1) GetAbTest(ctx context.Context, req *v1.GetAbTestReq) (res *v1.GetAbTestRes, err error) {
	res = &v1.GetAbTestRes{}

	test, err := abtest.AbTest().GetTestById(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	variants, err := abtest.AbTest().GetTestVariants(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = test
	_ = variants
	return
}

func (c *ControllerV1) ListAbTests(ctx context.Context, req *v1.ListAbTestsReq) (res *v1.ListAbTestsRes, err error) {
	res = &v1.ListAbTestsRes{}

	tests, err := abtest.AbTest().GetTestsForSequence(ctx, req.SequenceId)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = tests
	return
}

func (c *ControllerV1) GetAbTestResults(ctx context.Context, req *v1.GetAbTestResultsReq) (res *v1.GetAbTestResultsRes, err error) {
	res = &v1.GetAbTestResultsRes{}

	results, err := abtest.AbTest().GetTestResults(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = results
	return
}

func (c *ControllerV1) PickWinner(ctx context.Context, req *v1.PickWinnerReq) (res *v1.PickWinnerRes, err error) {
	res = &v1.PickWinnerRes{}

	err = abtest.AbTest().PickWinner(ctx, req.Id, req.WinnerVariant)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("Winner selected successfully")
	return
}

func (c *ControllerV1) DeleteAbTest(ctx context.Context, req *v1.DeleteAbTestReq) (res *v1.DeleteAbTestRes, err error) {
	res = &v1.DeleteAbTestRes{}

	err = abtest.AbTest().DeleteAbTest(ctx, req.Id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("AB test deleted successfully")
	return
}
