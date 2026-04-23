-- ============================================================
-- Migration 001: Cold Mail Features
-- - Sending delay, time windows, day scheduling
-- - Sender rotation, daily limits per sender
-- ============================================================

-- Add scheduling columns to email_tasks
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS send_delay INT DEFAULT 0;
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS schedule_start_hour INT DEFAULT 0;
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS schedule_end_hour INT DEFAULT 24;
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS schedule_days VARCHAR DEFAULT '[1,2,3,4,5,6,7]';
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS sender_pool TEXT DEFAULT '[]';
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS daily_limit_per_sender INT DEFAULT 0;
ALTER TABLE email_tasks ADD COLUMN IF NOT EXISTS current_sender_index INT DEFAULT 0;

-- Sender daily stats (track sends per sender per day)
CREATE TABLE IF NOT EXISTS bm_sender_daily_stats (
    id SERIAL PRIMARY KEY,
    sender_email VARCHAR(255) NOT NULL,
    stat_date DATE NOT NULL DEFAULT CURRENT_DATE,
    send_count INT NOT NULL DEFAULT 0,
    UNIQUE(sender_email, stat_date)
);

-- Sender configuration (display name, daily limit per sender)
CREATE TABLE IF NOT EXISTS bm_sender_config (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) DEFAULT '',
    daily_limit INT NOT NULL DEFAULT 100,
    is_active INT NOT NULL DEFAULT 1,
    create_time INT NOT NULL DEFAULT 0,
    update_time INT NOT NULL DEFAULT 0
);

-- Index for fast lookups
CREATE INDEX IF NOT EXISTS idx_sender_daily_stats_email_date ON bm_sender_daily_stats(sender_email, stat_date);
CREATE INDEX IF NOT EXISTS idx_sender_config_email ON bm_sender_config(email);
