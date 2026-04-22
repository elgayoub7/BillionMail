package scoring

import (
	"context"

t"time"
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/model/entity"
)

type ScoringService struct{}

var service = &ScoringService{}

func Scoring() *ScoringService {
	return service
}

// RecordEngagement updates engagement counters and recalculates score
func (s *ScoringService) RecordEngagement(ctx context.Context, email string, eventType string) error {
	_, err := g.DB().Model("bm_lead_scores").Save(g.Map{
		"contact_email": email,
	})
	if err != nil {
		return err
	}

	var field string
	switch eventType {
	case "open":
		field = "total_opens"
	case "click":
		field = "total_clicks"
	case "reply":
		field = "total_replies"
	case "bounce":
		field = "total_bounces"
	default:
		return nil
	}

	now := time.Now().Unix()
	_, err = g.DB().Exec(ctx, `
		UPDATE bm_lead_scores
		SET `+field+` = `+field+` + 1,
		    last_engagement_at = ?,
		    last_scored_at = ?
		WHERE contact_email = ?
	`, now, now, email)
	if err != nil {
		return err
	}

	return s.RecalculateScore(ctx, email)
}

// RecalculateScore computes the engagement score and level for a contact
func (s *ScoringService) RecalculateScore(ctx context.Context, email string) error {
	var lead entity.LeadScore
	err := g.DB().Model("bm_lead_scores").Where("contact_email = ?", email).Scan(&lead)
	if err != nil {
		return err
	}

	score := s.CalculateScore(lead.TotalOpens, lead.TotalClicks, lead.TotalReplies, lead.TotalBounces)
	level := s.GetEngagementLevel(score)

	_, err = g.DB().Model("bm_lead_scores").Where("contact_email = ?", email).Update(g.Map{
		"score":             score,
		"engagement_level":  level,
	})
	return err
}

// CalculateScore: (opens×2) + (clicks×5) + (replies×10) - (bounces×5)
func (s *ScoringService) CalculateScore(opens, clicks, replies, bounces int) int {
	score := (opens * 2) + (clicks * 5) + (replies * 10) - (bounces * 5)
	if score < 0 {
		score = 0
	}
	return score
}

// GetEngagementLevel returns the level based on score thresholds
func (s *ScoringService) GetEngagementLevel(score int) string {
	switch {
	case score >= 51:
		return "converted"
	case score >= 21:
		return "hot"
	case score >= 6:
		return "warm"
	default:
		return "cold"
	}
}

// GetLeadsByLevel returns paginated leads filtered by engagement level
func (s *ScoringService) GetLeadsByLevel(ctx context.Context, level string, page, pageSize int) (int, []entity.LeadScore, error) {
	m := g.DB().Model("bm_lead_scores")
	if level != "" && level != "all" {
		m = m.Where("engagement_level = ?", level)
	}

	total, err := m.Count()
	if err != nil {
		return 0, nil, err
	}

	var leads []entity.LeadScore
	err = m.Page(page, pageSize).Order("score DESC").Scan(&leads)
	return total, leads, err
}

// GetTopLeads returns the highest scoring leads
func (s *ScoringService) GetTopLeads(ctx context.Context, limit int) ([]entity.LeadScore, error) {
	var leads []entity.LeadScore
	err := g.DB().Model("bm_lead_scores").Order("score DESC").Limit(limit).Scan(&leads)
	return leads, err
}

// GetStats returns aggregate scoring statistics
func (s *ScoringService) GetStats(ctx context.Context) (*entity.ScoringStats, error) {
	stats := &entity.ScoringStats{}

	g.DB().Model("bm_lead_scores").Where("engagement_level", "cold").Count(&stats.ColdCount)
	g.DB().Model("bm_lead_scores").Where("engagement_level", "warm").Count(&stats.WarmCount)
	g.DB().Model("bm_lead_scores").Where("engagement_level", "hot").Count(&stats.HotCount)
	g.DB().Model("bm_lead_scores").Where("engagement_level", "converted").Count(&stats.ConvertedCount)
	stats.TotalLeads = stats.ColdCount + stats.WarmCount + stats.HotCount + stats.ConvertedCount

	if stats.TotalLeads > 0 {
		var avg struct {
			Avg float64 `json:"avg"`
		}
		g.DB().Model("bm_lead_scores").Fields("AVG(score) as avg").Scan(&avg)
		stats.AverageScore = avg.Avg
	}

	return stats, nil
}

// GetOrCreateLeadScore ensures a lead score entry exists
func (s *ScoringService) GetOrCreateLeadScore(ctx context.Context, email string) (*entity.LeadScore, error) {
	var lead entity.LeadScore
	err := g.DB().Model("bm_lead_scores").Where("contact_email = ?", email).Scan(&lead)
	if err != nil {
		return nil, err
	}
	if lead.Id > 0 {
		return &lead, nil
	}

	_, err = g.DB().Model("bm_lead_scores").Insert(g.Map{
		"contact_email":    email,
		"score":            0,
		"engagement_level": "cold",
	})
	if err != nil {
		return nil, err
	}

	g.DB().Model("bm_lead_scores").Where("contact_email = ?", email).Scan(&lead)
	return &lead, nil
}

// RecalculateAll forces a full recalculation of all lead scores
func (s *ScoringService) RecalculateAll(ctx context.Context) (int, error) {
	var leads []entity.LeadScore
	err := g.DB().Model("bm_lead_scores").Scan(&leads)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, lead := range leads {
		err := s.RecalculateScore(ctx, lead.ContactEmail)
		if err == nil {
			count++
		}
	}
	return count, nil
}
