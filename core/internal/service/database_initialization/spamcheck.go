package database_initialization

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	registerHandler(func() {
		var spamSQLList = []string{
			`CREATE TABLE IF NOT EXISTS bm_spam_scores (
				id SERIAL PRIMARY KEY,
				email_task_id INTEGER DEFAULT 0,
				sequence_step_id INTEGER DEFAULT 0,
				score REAL DEFAULT 0,
				action VARCHAR(20) DEFAULT '',
				symbols JSONB DEFAULT '[]',
				suggestions JSONB DEFAULT '[]',
				is_blocked INTEGER DEFAULT 0,
				checked_at INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_spam_scores_task ON bm_spam_scores(email_task_id)`,
			`CREATE TABLE IF NOT EXISTS bm_domain_health (
				id SERIAL PRIMARY KEY,
				domain VARCHAR(255) NOT NULL UNIQUE,
				spf_status VARCHAR(20) DEFAULT 'unknown',
				dkim_status VARCHAR(20) DEFAULT 'unknown',
				dmarc_status VARCHAR(20) DEFAULT 'unknown',
				mx_status VARCHAR(20) DEFAULT 'unknown',
				spf_record TEXT DEFAULT '',
				dkim_record TEXT DEFAULT '',
				dmarc_record TEXT DEFAULT '',
				mx_records TEXT DEFAULT '',
				issues JSONB DEFAULT '[]',
				last_checked INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_domain_health_domain ON bm_domain_health(domain)`,
		}

		for _, sql := range spamSQLList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute spamcheck SQL:", err, sql)
				return
			}
		}

		g.Log().Info(context.Background(), "Spam check + Domain health tables initialized successfully")
	})
}
