package v1

import (
	"context"

	"billionmail-core/internal/service/deliverability"
	"billionmail-core/utility/types/api_v1"
	"github.com/gogf/gf/v2/frame/g"
)

// SeedEmail represents a seed email for inbox placement testing
type SeedEmail struct {
	Email    string `json:"email"`
	Provider string `json:"provider"` // gmail, outlook, yahoo, icloud
	Active   bool   `json:"active"`
}

// --- Request types ---

type ListSeedsReq struct {
	g.Meta        `path:"/deliverability/seeds" method:"get" tags:"Deliverability" summary:"List seed emails"`
	Authorization string `json:"authorization" in:"header"`
}

type ListSeedsRes struct {
	api_v1.StandardRes
	Data struct {
		Seeds []SeedEmail `json:"seeds"`
	} `json:"data"`
}

type AddSeedReq struct {
	g.Meta        `path:"/deliverability/seeds" method:"post" tags:"Deliverability" summary:"Add seed email"`
	Authorization string `json:"authorization" in:"header"`
	Email         string `json:"email" v:"required|email" dc:"Seed email address"`
	Provider     string `json:"provider" v:"required|in:gmail,outlook,yahoo,icloud" dc:"Email provider"`
}

type AddSeedRes struct {
	api_v1.StandardRes
}

type RemoveSeedReq struct {
	g.Meta        `path:"/deliverability/seeds" method:"delete" tags:"Deliverability" summary:"Remove seed email"`
	Authorization string `json:"authorization" in:"header"`
	Email         string `json:"email" v:"required|email" dc:"Seed email address"`
}

type RemoveSeedRes struct {
	api_v1.StandardRes
}

type GetScoresReq struct {
	g.Meta        `path:"/deliverability/scores" method:"get" tags:"Deliverability" summary:"Get inbox placement scores"`
	Authorization string `json:"authorization" in:"header"`
}

type GetScoresRes struct {
	api_v1.StandardRes
	Data struct {
		Scores []map[string]interface{} `json:"scores"`
	} `json:"data"`
}

type TriggerTestReq struct {
	g.Meta        `path:"/deliverability/trigger-test" method:"post" tags:"Deliverability" summary:"Trigger inbox placement test"`
	Authorization string `json:"authorization" in:"header"`
}

type TriggerTestRes struct {
	api_v1.StandardRes
}

// --- Controller ---

type Controller struct{}

func (c *Controller) ListSeeds(ctx context.Context, req *ListSeedsReq, res *ListSeedsRes) error {
	svc := deliverability.NewSeedListService()
	seeds, err := svc.GetSeeds(ctx)
	if err != nil {
		return err
	}
	res.Data.Seeds = seeds
	return nil
}

func (c *Controller) AddSeed(ctx context.Context, req *AddSeedReq, res *AddSeedRes) error {
	svc := deliverability.NewSeedListService()
	err := svc.AddSeed(ctx, req.Email, req.Provider)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RemoveSeed(ctx context.Context, req *RemoveSeedReq, res *RemoveSeedRes) error {
	svc := deliverability.NewSeedListService()
	err := svc.RemoveSeed(ctx, req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetScores(ctx context.Context, req *GetScoresReq, res *GetScoresRes) error {
	svc := deliverability.NewInboxPlacementService()
	scores, err := svc.GetAllScores(ctx)
	if err != nil {
		return err
	}
	res.Data.Scores = scores
	return nil
}

func (c *Controller) TriggerTest(ctx context.Context, req *TriggerTestReq, res *TriggerTestRes) error {
	svc := deliverability.NewInboxPlacementService()
	err := svc.RunPlacementTests(ctx)
	if err != nil {
		return err
	}
	return nil
}
