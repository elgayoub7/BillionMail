package deliverability

import (
	"context"

	v1 "billionmail-core/api/deliverability/v1"
)

type IDeliverabilityV1 interface {
	ListSeeds(ctx context.Context, req *v1.ListSeedsReq, res *v1.ListSeedsRes) error
	AddSeed(ctx context.Context, req *v1.AddSeedReq, res *v1.AddSeedRes) error
	RemoveSeed(ctx context.Context, req *v1.RemoveSeedReq, res *v1.RemoveSeedRes) error
	GetScores(ctx context.Context, req *v1.GetScoresReq, res *v1.GetScoresRes) error
	TriggerTest(ctx context.Context, req *v1.TriggerTestReq, res *v1.TriggerTestRes) error
}
