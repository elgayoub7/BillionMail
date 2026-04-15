package database_initialization

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {

	registerHandler(func() {
		sequenceSQLList := []string{

			`CREATE TABLE IF NOT EXISTS bm_sequences (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT DEFAULT '',
				status SMALLINT NOT NULL DEFAULT 0,
				addresser VARCHAR(320) NOT NULL DEFAULT '',
				full_name VARCHAR(255) DEFAULT '',
				group_id INTEGER NOT NULL DEFAULT 0,
				tag_ids TEXT DEFAULT '',
				tag_logic VARCHAR(10) DEFAULT 'AND',
				track_open SMALLINT NOT NULL DEFAULT 1,
				track_click SMALLINT NOT NULL DEFAULT 1,
				unsubscribe SMALLINT NOT NULL DEFAULT 1,
				total_enrolled INTEGER NOT NULL DEFAULT 0,
				total_completed INTEGER NOT NULL DEFAULT 0,
				total_unsubscribed INTEGER NOT NULL DEFAULT 0,
				total_bounced INTEGER NOT NULL DEFAULT 0,
				create_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				UNIQUE(name)
			)`,

			`CREATE TABLE IF NOT EXISTS bm_sequence_steps (
				id SERIAL PRIMARY KEY,
				sequence_id INTEGER NOT NULL,
				step_order INTEGER NOT NULL DEFAULT 1,
				step_type VARCHAR(50) NOT NULL DEFAULT 'email',
				subject TEXT DEFAULT '',
				template_id INTEGER NOT NULL DEFAULT 0,
				wait_days INTEGER NOT NULL DEFAULT 0,
				wait_hours INTEGER NOT NULL DEFAULT 0,
				condition_type VARCHAR(50) DEFAULT '',
				condition_step_id INTEGER DEFAULT 0,
				on_true_go_to INTEGER DEFAULT 0,
				on_false_go_to INTEGER DEFAULT 0,
				sent_count INTEGER NOT NULL DEFAULT 0,
				opened_count INTEGER NOT NULL DEFAULT 0,
				clicked_count INTEGER NOT NULL DEFAULT 0,
				bounced_count INTEGER NOT NULL DEFAULT 0,
				create_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				FOREIGN KEY (sequence_id) REFERENCES bm_sequences(id) ON DELETE CASCADE,
				UNIQUE(sequence_id, step_order)
			)`,

			`CREATE TABLE IF NOT EXISTS bm_sequence_enrollments (
				id SERIAL PRIMARY KEY,
				sequence_id INTEGER NOT NULL,
				contact_id INTEGER NOT NULL DEFAULT 0,
				email VARCHAR(320) NOT NULL,
				group_id INTEGER NOT NULL DEFAULT 0,
				current_step INTEGER NOT NULL DEFAULT 1,
				status SMALLINT NOT NULL DEFAULT 0,
				enrolled_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				current_step_entered_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				last_email_sent_at INTEGER DEFAULT 0,
				completed_at INTEGER DEFAULT 0,
				total_emails_sent INTEGER NOT NULL DEFAULT 0,
				total_opens INTEGER NOT NULL DEFAULT 0,
				total_clicks INTEGER NOT NULL DEFAULT 0,
				FOREIGN KEY (sequence_id) REFERENCES bm_sequences(id) ON DELETE CASCADE,
				UNIQUE(sequence_id, email)
			)`,

			`CREATE TABLE IF NOT EXISTS bm_sequence_email_tasks (
				id SERIAL PRIMARY KEY,
				sequence_id INTEGER NOT NULL,
				enrollment_id INTEGER NOT NULL,
				step_id INTEGER NOT NULL,
				email_task_id INTEGER NOT NULL,
				contact_email VARCHAR(320) NOT NULL,
				status SMALLINT NOT NULL DEFAULT 0,
				sent_at INTEGER DEFAULT 0,
				message_id TEXT DEFAULT '',
				FOREIGN KEY (sequence_id) REFERENCES bm_sequences(id) ON DELETE CASCADE,
				FOREIGN KEY (enrollment_id) REFERENCES bm_sequence_enrollments(id) ON DELETE CASCADE,
				FOREIGN KEY (step_id) REFERENCES bm_sequence_steps(id) ON DELETE CASCADE
			)`,

			// Indexes
			`CREATE INDEX IF NOT EXISTS idx_bm_sequences_status ON bm_sequences(status)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_steps_sequence_id ON bm_sequence_steps(sequence_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_steps_sequence_order ON bm_sequence_steps(sequence_id, step_order)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_enrollments_sequence_status ON bm_sequence_enrollments(sequence_id, status)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_enrollments_step_status ON bm_sequence_enrollments(sequence_id, current_step, status)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_enrollments_email ON bm_sequence_enrollments(email)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_email_tasks_enrollment ON bm_sequence_email_tasks(enrollment_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_email_tasks_step ON bm_sequence_email_tasks(step_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_sequence_email_tasks_task_id ON bm_sequence_email_tasks(email_task_id)`,
		}

		for _, sql := range sequenceSQLList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute sequence SQL:", err, sql)
				return
			}
		}
	})
}
