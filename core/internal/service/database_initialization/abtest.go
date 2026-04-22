package database_initialization

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	registerHandler(func() {
		var abtestSQLList = []string{
			`CREATE TABLE IF NOT EXISTS bm_ab_tests (
				id SERIAL PRIMARY KEY,
				sequence_id INTEGER NOT NULL,
				step_id INTEGER NOT NULL DEFAULT 0,
				name VARCHAR(255) DEFAULT '',
				status INTEGER NOT NULL DEFAULT 0,
				winner_variant INTEGER DEFAULT -1,
				winner_criteria VARCHAR(50) DEFAULT 'open_rate',
				split_ratio REAL DEFAULT 0.5,
				confidence_threshold REAL DEFAULT 0.95,
				created_at INTEGER DEFAULT EXTRACT(EPOCH FROM NOW()),
				completed_at INTEGER DEFAULT 0
			)`,
			`CREATE TABLE IF NOT EXISTS bm_ab_test_variants (
				id SERIAL PRIMARY KEY,
				ab_test_id INTEGER NOT NULL REFERENCES bm_ab_tests(id) ON DELETE CASCADE,
				variant INTEGER NOT NULL,
				subject TEXT DEFAULT '',
				body_html TEXT DEFAULT '',
				body_text TEXT DEFAULT '',
				sent_count INTEGER DEFAULT 0,
				open_count INTEGER DEFAULT 0,
				click_count INTEGER DEFAULT 0,
				reply_count INTEGER DEFAULT 0,
				bounce_count INTEGER DEFAULT 0
			)`,
			`CREATE TABLE IF NOT EXISTS bm_ab_test_assignments (
				id SERIAL PRIMARY KEY,
				ab_test_id INTEGER NOT NULL,
				enrollment_id INTEGER NOT NULL,
				variant INTEGER NOT NULL,
				assigned_at INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_tests_sequence ON bm_ab_tests(sequence_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_tests_step ON bm_ab_tests(sequence_id, step_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_tests_status ON bm_ab_tests(status)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_test_variants_test ON bm_ab_test_variants(ab_test_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_test_assignments_test ON bm_ab_test_assignments(ab_test_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_ab_test_assignments_enrollment ON bm_ab_test_assignments(enrollment_id)`,
		}

		for _, sql := range abtestSQLList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute abtest SQL:", err, sql)
				return
			}
		}

		g.Log().Info(context.Background(), "AB Test tables initialized successfully")
	})
}
