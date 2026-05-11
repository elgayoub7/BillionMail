package deliverability

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
)

const seedListOptionKey = "inbox_placement_seeds"

type SeedEmail struct {
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Active   bool   `json:"active"`
}

type SeedListService struct{}

func NewSeedListService() *SeedListService {
	return &SeedListService{}
}

func (s *SeedListService) GetSeeds(ctx context.Context) ([]SeedEmail, error) {
	val, err := g.DB().Model("bm_options").
		Where("name", seedListOptionKey).
		Value("value")
	if err != nil {
		return nil, err
	}
	if val.IsEmpty() {
		return []SeedEmail{}, nil
	}
	var seeds []SeedEmail
	if err := json.Unmarshal([]byte(val.String()), &seeds); err != nil {
		return nil, err
	}
	return seeds, nil
}

func (s *SeedListService) AddSeed(ctx context.Context, email, provider string) error {
	seeds, err := s.GetSeeds(ctx)
	if err != nil {
		return err
	}
	for i, seed := range seeds {
		if seed.Email == email {
			seeds[i].Provider = provider
			seeds[i].Active = true
			return s.saveSeeds(ctx, seeds)
		}
	}
	seeds = append(seeds, SeedEmail{Email: email, Provider: provider, Active: true})
	return s.saveSeeds(ctx, seeds)
}

func (s *SeedListService) RemoveSeed(ctx context.Context, email string) error {
	seeds, err := s.GetSeeds(ctx)
	if err != nil {
		return err
	}
	filtered := make([]SeedEmail, 0, len(seeds))
	for _, seed := range seeds {
		if seed.Email != email {
			filtered = append(filtered, seed)
		}
	}
	return s.saveSeeds(ctx, filtered)
}

func (s *SeedListService) GetActiveSeeds(ctx context.Context) ([]SeedEmail, error) {
	seeds, err := s.GetSeeds(ctx)
	if err != nil {
		return nil, err
	}
	active := make([]SeedEmail, 0, len(seeds))
	for _, seed := range seeds {
		if seed.Active {
			active = append(active, seed)
		}
	}
	return active, nil
}

func (s *SeedListService) saveSeeds(ctx context.Context, seeds []SeedEmail) error {
	data, err := json.Marshal(seeds)
	if err != nil {
		return err
	}
	_, err = g.DB().Model("bm_options").
		Data("name", seedListOptionKey).
		Data("value", string(data)).
		Replace()
	return err
}
