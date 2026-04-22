package database_initialization

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	registerHandler(func() {
		var scoringSQLList = []string{
			`CREATE TABLE IF NOT EXISTS bm_lead_scores (
				id SERIAL PRIMARY KEY,
				contact_email TEXT NOT NULL UNIQUE,
				score INTEGER DEFAULT 0,
				engagement_level VARCHAR(20) DEFAULT 'cold',
				total_opens INTEGER DEFAULT 0,
				total_clicks INTEGER DEFAULT 0,
				total_replies INTEGER DEFAULT 0,
				total_bounces INTEGER DEFAULT 0,
				last_engagement_at INTEGER DEFAULT 0,
				last_scored_at INTEGER DEFAULT EXTRACT(EPOCH FROM NOW()),
				metadata JSONB DEFAULT '{}'
			)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_lead_scores_email ON bm_lead_scores(contact_email)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_lead_scores_level ON bm_lead_scores(engagement_level)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_lead_scores_score ON bm_lead_scores(score DESC)`,
			`CREATE TABLE IF NOT EXISTS bm_bounce_records (
				id SERIAL PRIMARY KEY,
				email TEXT NOT NULL,
				bounce_type VARCHAR(20) NOT NULL,
				description TEXT DEFAULT '',
				sequence_id INTEGER DEFAULT 0,
				message_id TEXT DEFAULT '',
				action_taken VARCHAR(50) DEFAULT '',
				created_at INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_bounce_records_email ON bm_bounce_records(email)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_bounce_records_type ON bm_bounce_records(bounce_type)`,
		}

		for _, sql := range scoringSQLList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute scoring SQL:", err, sql)
				return
			}
		}

		g.Log().Info(context.Background(), "Scoring + Bounce tables initialized successfully")
	})
}
