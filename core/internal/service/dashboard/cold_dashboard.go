package dashboard

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type ColdDashboardService struct{}

var service = &ColdDashboardService{}

func ColdDashboard() *ColdDashboardService {
	return service
}

type DashboardStats struct {
	EmailsSentToday int     `json:"emails_sent_today"`
	EmailsSentWeek  int     `json:"emails_sent_week"`
	EmailsSentMonth int     `json:"emails_sent_month"`
	OpenRate        float64 `json:"open_rate"`
	ClickRate       float64 `json:"click_rate"`
	ReplyRate       float64 `json:"reply_rate"`
	BounceRate      float64 `json:"bounce_rate"`
	ActiveSequences int     `json:"active_sequences"`
	ActiveWarmups   int     `json:"active_warmups"`
}

type ActiveSequenceMetrics struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Status      int     `json:"status"`
	Enrolled    int     `json:"enrolled"`
	Completed   int     `json:"completed"`
	ReplyRate   float64 `json:"reply_rate"`
	OpenRate    float64 `json:"open_rate"`
	BounceRate  float64 `json:"bounce_rate"`
	HasAbTest   bool    `json:"has_ab_test"`
}

type AlertItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int64  `json:`
}

// GetDashboardStats aggregates global stats
func (s *ColdDashboardService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	now := time.Now().Unix()
	dayAgo := now - 86400
	weekAgo := now - 604800
	monthAgo := now - 2592000

	g.DB().Model("recipient_info").Where("create_time >= ?", dayAgo).Count(&stats.EmailsSentToday)
	g.DB().Model("recipient_info").Where("create_time >= ?", weekAgo).Count(&stats.EmailsSentWeek)
	g.DB().Model("recipient_info").Where("create_time >= ?", monthAgo).Count(&stats.EmailsSentMonth)

	sentCount := 0
	g.DB().Model("recipient_info").Where("is_sent = 1").Where("create_time >= ?", monthAgo).Count(&sentCount)
	if sentCount > 0 {
		openCount := 0
		g.DB().Model("mailstat_opened").Where("log_time >= ?", monthAgo).Count(&openCount)
		clickCount := 0
		g.DB().Model("mailstat_clicked").Where("log_time >= ?", monthAgo).Count(&clickCount)
		bounceCount := 0
		g.DB().Model("mailstat_send_mails").Where("status = 'bounced'").Where("log_time >= ?", monthAgo).Count(&bounceCount)

		stats.OpenRate = float64(openCount) / float64(sentCount) * 100
		stats.ClickRate = float64(clickCount) / float64(sentCount) * 100
		stats.BounceRate = float64(bounceCount) / float64(sentCount) * 100
	}

	g.DB().Model("bm_sequences").Where("status = 1").Count(&stats.ActiveSequences)
	g.DB().Model("sender_ip_warmup").Where("warmup_period > 0 AND current_day < warmup_period").Count(&stats.ActiveWarmups)

	return stats, nil
}

// GetActiveSequencesWithMetrics returns active sequences with computed metrics
func (s *ColdDashboardService) GetActiveSequencesWithMetrics(ctx context.Context) ([]ActiveSequenceMetrics, error) {
	var sequences []struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	}
	g.DB().Model("bm_sequences").Where("status IN (0,1)").Order("id DESC").Scan(&sequences)

	result := []ActiveSequenceMetrics{}
	for _, seq := range sequences {
		metrics := ActiveSequenceMetrics{
			Id:     seq.Id,
			Name:   seq.Name,
			Status: seq.Status,
		}

		g.DB().Model("bm_sequence_enrollments").Where("sequence_id = ? AND status = 0", seq.Id).Count(&metrics.Enrolled)
		g.DB().Model("bm_sequence_enrollments").Where("sequence_id = ? AND status = 1", seq.Id).Count(&metrics.Completed)

		totalSent := 0
		g.DB().Model("bm_sequence_enrollments").Where("sequence_id = ?", seq.Id).
			Where("total_emails_sent > 0").Count(&totalSent)

		if totalSent > 0 {
			var openSum struct {
				Total int `json:"total"`
			}
			g.DB().Model("bm_sequence_enrollments").Where("sequence_id = ?", seq.Id).
				Fields("COALESCE(SUM(total_opens),0) as total").Scan(&openSum)
			metrics.OpenRate = float64(openSum.Total) / float64(totalSent) * 100

			var bounceSum struct {
				Total int `json:"total"`
			}
			g.DB().Model("bm_sequence_enrollments").Where("sequence_id = ?", seq.Id).
				Fields("COALESCE(SUM(total_bounces),0) as total").Scan(&bounceSum)
			metrics.BounceRate = float64(bounceSum.Total) / float64(totalSent) * 100
		}

		abTestCount, _ := g.DB().Model("bm_ab_tests").Where("sequence_id = ? AND status = 1", seq.Id).Count()
		metrics.HasAbTest = abTestCount > 0

		result = append(result, metrics)
	}

	return result, nil
}

// GetRecentAlerts compiles recent alerts from various subsystems
func (s *ColdDashboardService) GetRecentAlerts(ctx context.Context, limit int) ([]AlertItem, error) {
	alerts := []AlertItem{}

	// Hard bounces
	var recentBounces []struct {
		Email     string `json:"email"`
		CreatedAt int    `json:"created_at"`
	}
	g.DB().Model("bm_bounce_records").Where("bounce_type = 'hard'").Order("created_at DESC").Limit(5).Scan(&recentBounces)
	for _, b := range recentBounces {
		alerts = append(alerts, AlertItem{
			Type:    "hard_bounce",
			Message: "Hard bounce: " + b.Email,
			Time:    int64(b.CreatedAt),
		})
	}

	// Domain health issues
	var unhealthyDomains []struct {
		Domain     string `json:"domain"`
		SpfStatus  string `json:"spf_status"`
		DkimStatus string `json:"dkim_status"`
	}
	g.DB().Model("bm_domain_health").Where("spf_status = 'missing' OR dkim_status = 'missing'").Limit(5).Scan(&unhealthyDomains)
	for _, d := range unhealthyDomains {
		alerts = append(alerts, AlertItem{
			Type:    "domain_health",
			Message: "DNS issue on " + d.Domain,
			Time:    0,
		})
	}

	// AB test winners
	var winners []struct {
		Id          int `json:"id"`
		Winner      int `json:"winner_variant"`
		CompletedAt int `json:"completed_at"`
	}
	g.DB().Model("bm_ab_tests").Where("status = 2 AND completed_at > 0").Order("completed_at DESC").Limit(3).Scan(&winners)
	for _, w := range winners {
		variant := "A"
		if w.Winner == 1 {
			variant = "B"
		}
		alerts = append(alerts, AlertItem{
			Type:    "abtest_winner",
			Message: "AB test winner found: variant " + variant,
			Time:    int64(w.CompletedAt),
		})
	}

	if len(alerts) > limit {
		alerts = alerts[:limit]
	}

	return alerts, nil
}
