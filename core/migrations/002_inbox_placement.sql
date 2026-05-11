-- ============================================================
-- Migration 002: Inbox Placement Testing
-- - Seed list storage
-- - Test result tracking
-- ============================================================

CREATE TABLE IF NOT EXISTS bm_inbox_placement_results (
    id SERIAL PRIMARY KEY,
    sender_pool_id INT REFERENCES bm_sender_pool(id) ON DELETE CASCADE,
    seed_email VARCHAR(255) NOT NULL,
    test_email_message_id VARCHAR(255) NOT NULL UNIQUE,
    sent_at TIMESTAMP NOT NULL,
    received_at TIMESTAMP,
    placement VARCHAR(20) DEFAULT 'none',  -- inbox | spam | none
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inbox_placement_sender ON bm_inbox_placement_results(sender_pool_id);
CREATE INDEX IF NOT EXISTS idx_inbox_placement_message_id ON bm_inbox_placement_results(test_email_message_id);
CREATE INDEX IF NOT EXISTS idx_inbox_placement_sent_at ON bm_inbox_placement_results(sent_at);

-- Add inbox_rate column to sender pool for quick access
ALTER TABLE bm_sender_pool ADD COLUMN IF NOT EXISTS inbox_rate INT DEFAULT 0;
ALTER TABLE bm_sender_pool ADD COLUMN IF NOT EXISTS inbox_placement_status VARCHAR(20) DEFAULT 'unknown';