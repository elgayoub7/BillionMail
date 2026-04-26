package batch_mail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SenderPoolEntry represents a single sender in the rotation pool
type SenderPoolEntry struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// SenderRotation manages sender partition and daily limits
type SenderRotation struct {
	pool       []SenderPoolEntry
	dailyLimit int // per-sender daily limit (0 = unlimited)
	taskID     int
}

// NewSenderRotation creates a new sender rotation manager
func NewSenderRotation(poolJSON string, dailyLimit int, taskID int, startIndex int) *SenderRotation {
	sr := &SenderRotation{
		dailyLimit: dailyLimit,
		taskID:     taskID,
	}

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

// GetSenderForRecipient assigns a sender to a recipient based on index partition.
// Each sender gets a fixed slice of recipients: recipient[i] → pool[i % poolSize].
func (sr *SenderRotation) GetSenderForRecipient(ctx context.Context, recipientIndex int) (*SenderPoolEntry, error) {
	if len(sr.pool) == 0 {
		return nil, nil
	}

	idx := recipientIndex % len(sr.pool)
	sender := &sr.pool[idx]

	// Check daily limit
	if sr.dailyLimit > 0 {
		count, err := GetSenderDailyCount(ctx, sender.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check daily count for %s: %w", sender.Email, err)
		}
		if count >= sr.dailyLimit {
			return nil, fmt.Errorf("sender %s hit daily limit (%d/%d)", sender.Email, count, sr.dailyLimit)
		}
	}

	return sender, nil
}

// GetSenderDailyCount returns the number of emails sent by a sender today
func GetSenderDailyCount(ctx context.Context, email string) (int, error) {
	today := time.Now().Format("2006-01-02")

	countResult, err := g.DB().Ctx(ctx).Model("bm_sender_daily_stats").
		Where("sender_email", email).
		Where("stat_date", today).
		Fields("COALESCE(SUM(send_count), 0)").
		Value()
	if err != nil {
		return 0, err
	}
	count := countResult.Int()

	return count, nil
}

// IncrementSenderDailyCount atomically increments the daily send count for a sender
func IncrementSenderDailyCount(ctx context.Context, email string) error {
	today := time.Now().Format("2006-01-02")

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
