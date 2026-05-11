package deliverability

import (
	"context"

	v1 "billionmail-core/api/deliverability/v1"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) ListSeeds(ctx context.Context, req *v1.ListSeedsReq, res *v1.ListSeedsRes) error {
	return v1.ListSeeds(ctx, req, res)
}

func (c *ControllerV1) AddSeed(ctx context.Context, req *v1.AddSeedReq, res *v1.AddSeedRes) error {
	return v1.AddSeed(ctx, req, res)
}

func (c *ControllerV1) RemoveSeed(ctx context.Context, req *v1.RemoveSeedReq, res *v1.RemoveSeedRes) error {
	return v1.RemoveSeed(ctx, req, res)
}

func (c *ControllerV1) GetScores(ctx context.Context, req *v1.GetScoresReq, res *v1.GetScoresRes) error {
	return v1.GetScores(ctx, req, res)
}

func (c *ControllerV1) TriggerTest(ctx context.Context, req *v1.TriggerTestReq, res *v1.TriggerTestRes) error {
	return v1.TriggerTest(ctx, req, res)
}