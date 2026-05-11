package deliverability

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const seedListOptionKey = "inbox_placement_seeds"

// SeedEmail represents a seed email for inbox placement testing
type SeedEmail struct {
	Email    string `json:"email"`
	Provider string `json:"provider"` // gmail, outlook, yahoo, icloud
	Active   bool   `json:"active"`
}

// SeedListService handles seed email management for inbox placement
type SeedListService struct{}

// NewSeedListService creates a new SeedListService
func NewSeedListService() *SeedListService {
	return &SeedListService{}
}

// GetSeeds returns all seeds from bm_options
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
	if err := gconv.Structs(gconv.Map(gconv.StringToJson(val.String())), &seeds); err != nil {
		return nil, err
	}
	return seeds, nil
}

// AddSeed adds a seed, inserts if not exists
func (s *SeedListService) AddSeed(ctx context.Context, email, provider string) error {
	seeds, err := s.GetSeeds(ctx)
	if err != nil {
		return err
	}

	// Check if email already exists
	for i, seed := range seeds {
		if seed.Email == email {
			// Update existing
			seeds[i].Provider = provider
			seeds[i].Active = true
			return s.saveSeeds(ctx, seeds)
		}
	}

	// Add new seed
	seeds = append(seeds, SeedEmail{
		Email:    email,
		Provider: provider,
		Active:   true,
	})
	return s.saveSeeds(ctx, seeds)
}

// RemoveSeed removes a seed by email
func (s *SeedListService) RemoveSeed(ctx context.Context, email string) error {
	seeds, err := s.GetSeeds(ctx)
	if err != nil {
		return err
	}

	// Filter out the seed with the given email
	filtered := make([]SeedEmail, 0, len(seeds))
	for _, seed := range seeds {
		if seed.Email != email {
			filtered = append(filtered, seed)
		}
	}
	return s.saveSeeds(ctx, filtered)
}

// GetActiveSeeds returns only active seeds
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

// saveSeeds persists the seeds slice to bm_options
func (s *SeedListService) saveSeeds(ctx context.Context, seeds []SeedEmail) error {
	_, err := g.DB().Model("bm_options").
		Data("name", seedListOptionKey).
		Data("value", gconv.String(gconv.StringToJson(seeds))).
		Replace()
	return err
}