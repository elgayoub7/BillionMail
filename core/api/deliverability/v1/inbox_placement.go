package v1

import (
	"context"

	"billionmail-core/internal/service/deliverability"
	"billionmail-core/utility/types/api_v1"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct{}

type SeedEmail struct {
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Active   bool   `json:"active"`
}

type ListSeedsReq struct {
	g.Meta `path:"/deliverability/seeds" method:"get" tags:"Deliverability" summary:"List seed emails"`
}

type ListSeedsRes struct {
	api_v1.StandardRes
	Data struct {
		Seeds []SeedEmail `json:"seeds"`
	}
}

type AddSeedReq struct {
	g.Meta   `path:"/deliverability/seeds" method:"post" tags:"Deliverability" summary:"Add seed email"`
	Email    string `json:"email" v:"required|email"`
	Provider string `json:"provider" v:"required|in:gmail,outlook,yahoo,icloud"`
}

type AddSeedRes struct {
	api_v1.StandardRes
	Message string `json:"message"`
}

type RemoveSeedReq struct {
	g.Meta `path:"/deliverability/seeds" method:"delete" tags:"Deliverability" summary:"Remove seed email"`
	Email  string `json:"email" v:"required|email"`
}

type RemoveSeedRes struct {
	api_v1.StandardRes
	Message string `json:"message"`
}

type GetScoresReq struct {
	g.Meta `path:"/deliverability/scores" method:"get" tags:"Deliverability" summary:"Get inbox placement scores"`
}

type GetScoresRes struct {
	api_v1.StandardRes
	Data struct {
		Scores []map[string]interface{} `json:"scores"`
	}
}

type TriggerTestReq struct {
	g.Meta `path:"/deliverability/trigger-test" method:"post" tags:"Deliverability" summary:"Trigger inbox placement test"`
}

type TriggerTestRes struct {
	api_v1.StandardRes
	Message string `json:"message"`
}

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) ListSeeds(ctx context.Context, req *ListSeedsReq) (res *ListSeedsRes, err error) {
	res = &ListSeedsRes{}
	svc := deliverability.NewSeedListService()
	seeds, err := svc.GetSeeds(ctx)
	if err != nil {
		return nil, err
	}
	for _, seed := range seeds {
		res.Data.Seeds = append(res.Data.Seeds, SeedEmail{
			Email:    seed.Email,
			Provider: seed.Provider,
			Active:   seed.Active,
		})
	}
	return res, nil
}

func (c *ControllerV1) AddSeed(ctx context.Context, req *AddSeedReq) (res *AddSeedRes, err error) {
	res = &AddSeedRes{}
	svc := deliverability.NewSeedListService()
	err = svc.AddSeed(ctx, req.Email, req.Provider)
	if err != nil {
		return nil, err
	}
	res.Message = "Seed added"
	return res, nil
}

func (c *ControllerV1) RemoveSeed(ctx context.Context, req *RemoveSeedReq) (res *RemoveSeedRes, err error) {
	res = &RemoveSeedRes{}
	svc := deliverability.NewSeedListService()
	err = svc.RemoveSeed(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	res.Message = "Seed removed"
	return res, nil
}

func (c *ControllerV1) GetScores(ctx context.Context, req *GetScoresReq) (res *GetScoresRes, err error) {
	res = &GetScoresRes{}
	svc := deliverability.NewInboxPlacementService()
	scores, err := svc.GetAllScores(ctx)
	if err != nil {
		return nil, err
	}
	res.Data.Scores = scores
	return res, nil
}

func (c *ControllerV1) TriggerTest(ctx context.Context, req *TriggerTestReq) (res *TriggerTestRes, err error) {
	res = &TriggerTestRes{}
	svc := deliverability.NewInboxPlacementService()
	err = svc.RunPlacementTests(ctx)
	if err != nil {
		return nil, err
	}
	res.Message = "Test triggered"
	return res, nil
}
