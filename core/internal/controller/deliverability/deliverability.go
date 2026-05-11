package deliverability

import (
	"context"

	v1 "billionmail-core/api/deliverability/v1"
	"billionmail-core/internal/service/deliverability"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) ListSeeds(ctx context.Context, req *v1.ListSeedsReq) (res *v1.ListSeedsRes, err error) {
	res = &v1.ListSeedsRes{}
	svc := deliverability.NewSeedListService()
	seeds, err := svc.GetSeeds(ctx)
	if err != nil {
		return nil, err
	}
	for _, seed := range seeds {
		res.Data.Seeds = append(res.Data.Seeds, v1.SeedEmail{
			Email:    seed.Email,
			Provider: seed.Provider,
			Active:   seed.Active,
		})
	}
	return res, nil
}

func (c *ControllerV1) AddSeed(ctx context.Context, req *v1.AddSeedReq) (res *v1.AddSeedRes, err error) {
	res = &v1.AddSeedRes{}
	svc := deliverability.NewSeedListService()
	err = svc.AddSeed(ctx, req.Email, req.Provider)
	if err != nil {
		return nil, err
	}
	res.Message = "Seed added"
	return res, nil
}

func (c *ControllerV1) RemoveSeed(ctx context.Context, req *v1.RemoveSeedReq) (res *v1.RemoveSeedRes, err error) {
	res = &v1.RemoveSeedRes{}
	svc := deliverability.NewSeedListService()
	err = svc.RemoveSeed(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	res.Message = "Seed removed"
	return res, nil
}

func (c *ControllerV1) GetScores(ctx context.Context, req *v1.GetScoresReq) (res *v1.GetScoresRes, err error) {
	res = &v1.GetScoresRes{}
	svc := deliverability.NewInboxPlacementService()
	scores, err := svc.GetAllScores(ctx)
	if err != nil {
		return nil, err
	}
	res.Data.Scores = scores
	return res, nil
}

func (c *ControllerV1) TriggerTest(ctx context.Context, req *v1.TriggerTestReq) (res *v1.TriggerTestRes, err error) {
	res = &v1.TriggerTestRes{}
	svc := deliverability.NewInboxPlacementService()
	err = svc.RunPlacementTests(ctx)
	if err != nil {
		return nil, err
	}
	res.Message = "Test triggered"
	return res, nil
}
