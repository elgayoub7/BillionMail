package batch_mail

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SenderPoolEntry represents a single sender in the rotation pool
type SenderPoolEntry struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// SenderRotation manages round-robin sender rotation and daily limits
type SenderRotation struct {
	currentIndex atomic.Int64
	pool         []SenderPoolEntry
	dailyLimit   int // per-sender daily limit (0 = unlimited)
	taskID       int
}

// NewSenderRotation creates a new sender rotation manager
func NewSenderRotation(poolJSON string, dailyLimit int, taskID int, startIndex int) *SenderRotation {
	sr := &SenderRotation{
		dailyLimit: dailyLimit,
		taskID:     taskID,
	}
	sr.currentIndex.Store(int64(startIndex))

	if poolJSON != "" && poolJSON != "[]" {
		if err := json.Unmarshal([]byte(poolJSON), &sr.pool); err != nil {
			g.Log().Warningf(context.Background(), "Failed to parse sender_pool JSON: %v", err)
		}
	}

	return sr
}

// PoolSize returns the number of senders in the pool
func (sr *SenderRotation) PoolSize() int {
	return len(sr.pool)
}

// GetNextSender returns the next available sender (round-robin with daily limit check)
// Returns nil if no rotation (single sender mode)
func (sr *SenderRotation) GetNextSender(ctx context.Context) (*SenderPoolEntry, error) {
	if len(sr.pool) == 0 {
		return nil, nil // no rotation, use default task sender
	}

	// Try all senders in the pool starting from current index
	startIdx := sr.currentIndex.Load()
	for i := 0; i < len(sr.pool); i++ {
		idx := (startIdx + int64(i)) % int64(len(sr.pool))
		sender := &sr.pool[idx]

		// Check daily limit
		if sr.dailyLimit > 0 {
			count, err := GetSenderDailyCount(ctx, sender.Email)
			if err != nil {
				g.Log().Warningf(ctx, "Failed to check daily count for %s: %v", sender.Email, err)
				continue
			}
			if count >= sr.dailyLimit {
				g.Log().Debugf(ctx, "Sender %s hit daily limit (%d/%d), skipping", sender.Email, count, sr.dailyLimit)
				continue
			}
		}

		// This sender is available — advance index
		nextIdx := (idx + 1) % int64(len(sr.pool))
		sr.currentIndex.Store(nextIdx)

		// Persist index in DB for crash recovery
		sr.persistIndex(ctx)

		return sender, nil
	}

	// All senders hit their daily limit
	return nil, fmt.Errorf("all %d sender(s) in pool have reached their daily limit (%d emails/day)", len(sr.pool), sr.dailyLimit)
}

// persistIndex saves current sender index to DB
func (sr *SenderRotation) persistIndex(ctx context.Context) {
	_, err := g.DB().Ctx(ctx).Model("email_tasks").
		Where("id", sr.taskID).
		Data(g.Map{"current_sender_index": sr.currentIndex.Load()}).
		Update()
	if err != nil {
		g.Log().Warningf(ctx, "Failed to persist sender index for task %d: %v", sr.taskID, err)
	}
}

// GetSenderDailyCount returns the number of emails sent by a sender today
func GetSenderDailyCount(ctx context.Context, email string) (int, error) {
	today := time.Now().Format("2006-01-02")

	var count int
	err := g.DB().Ctx(ctx).Model("bm_sender_daily_stats").
		Where("sender_email", email).
		Where("stat_date", today).
		Fields("COALESCE(SUM(send_count), 0)").
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// IncrementSenderDailyCount atomically increments the daily send count for a sender
func IncrementSenderDailyCount(ctx context.Context, email string) error {
	today := time.Now().Format("2006-01-02")

	// Upsert: insert or increment
	_, err := g.DB().Ctx(ctx).Model("bm_sender_daily_stats").
		Data(g.Map{
			"sender_email": email,
			"stat_date":    today,
			"send_count":   1,
		}).
		OnConflict("sender_email,stat_date").
		OnDuplicate(g.Map{
			"send_count": gdb.Raw("bm_sender_daily_stats.send_count + 1"),
		}).
		Save()

	return err
}
