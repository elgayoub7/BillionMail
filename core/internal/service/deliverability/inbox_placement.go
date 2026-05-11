package deliverability

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/google/uuid"
)

type InboxPlacementService struct {
	seedList *SeedListService
	ctx      context.Context
}

func NewInboxPlacementService() *InboxPlacementService {
	return &InboxPlacementService{
		seedList: NewSeedListService(),
		ctx:      gctx.GetIOContext(),
	}
}

// StartTimer registers the 15-minute inbox placement test timer
func (s *InboxPlacementService) StartTimer() {
	gtimer.Add(15*time.Minute, func(ctx context.Context) {
		s.RunPlacementTests(ctx)
	})
}

// RunPlacementTests sends test emails to all active seeds from all active senders
func (s *InboxPlacementService) RunPlacementTests(ctx context.Context) error {
	seeds, err := s.seedList.GetActiveSeeds(ctx)
	if err != nil || len(seeds) == 0 {
		return err
	}

	var senders []struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	err = g.DB().Model("bm_sender_pool").
		Where("status", "active").
		Scan(&senders)
	if err != nil {
		return err
	}

	for _, sender := range senders {
		s.sendTestEmailsForSender(ctx, sender.Id, sender.Username, sender.Email, seeds)
	}
	return nil
}

// sendTestEmailsForSender sends one test email per seed for a given sender
// message_id format: BM-INBOX-TEST-{uuid}
func (s *InboxPlacementService) sendTestEmailsForSender(ctx context.Context, senderId int, senderUsername, senderEmail string, seeds []SeedEmail) error {
	for _, seed := range seeds {
		messageId := fmt.Sprintf("BM-INBOX-TEST-%s", uuid.New().String())
		subject := fmt.Sprintf("BM-INBOX-TEST-%s", uuid.New().String()[:8])

		_, err := g.DB().Model("bm_inbox_placement_results").
			Data(g.Map{
				"sender_pool_id":       senderId,
				"seed_email":           seed.Email,
				"test_email_message_id": messageId,
				"sent_at":               time.Now(),
				"placement":            "none",
			}).
			Insert()
		if err != nil {
			continue
		}

		// Log the test email that would be sent
		fmt.Printf("[InboxPlacement] TEST sender=%s seed=%s message_id=%s\n",
			senderUsername, seed.Email, messageId)
	}
	return nil
}

// UpdatePlacementFromReply processes a reply from Maildir scanner and updates placement
func (s *InboxPlacementService) UpdatePlacementFromReply(ctx context.Context, messageId string, receivedAt time.Time) error {
	if !strings.HasPrefix(messageId, "BM-INBOX-TEST-") {
		return nil
	}

	_, err := g.DB().Model("bm_inbox_placement_results").
		Where("test_email_message_id", messageId).
		Data(g.Map{
			"received_at": receivedAt,
			"placement":   "inbox",
		}).
		Update()
	return err
}

// CalculateScore computes inbox_rate for a sender over a rolling window (default 30 days)
func (s *InboxPlacementService) CalculateScore(ctx context.Context, senderPoolId int, windowDays int) (int, string, error) {
	cutoff := time.Now().AddDate(0, 0, -windowDays)

	total, err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_pool_id", senderPoolId).
		Where("sent_at > ?", cutoff).
		Count()
	if err != nil {
		return 0, "unknown", err
	}
	if total == 0 {
		return 0, "unknown", nil
	}

	inboxCount, err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_pool_id", senderPoolId).
		Where("sent_at > ?", cutoff).
		Where("placement", "inbox").
		Count()
	if err != nil {
		return 0, "unknown", err
	}

	rate := (inboxCount * 100) / total
	status := "good"
	if rate < 50 {
		status = "poor"
	} else if rate < 80 {
		status = "warning"
	}

	return rate, status, nil
}

// GetAllScores returns scores for all active senders
func (s *InboxPlacementService) GetAllScores(ctx context.Context) ([]map[string]interface{}, error) {
	var senders []struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	err := g.DB().Model("bm_sender_pool").
		Where("status", "active").
		Scan(&senders)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for _, sender := range senders {
		rate, status, _ := s.CalculateScore(ctx, sender.Id, 30)
		results = append(results, map[string]interface{}{
			"sender_id":  sender.Id,
			"username":   sender.Username,
			"email":      sender.Email,
			"inbox_rate": rate,
			"status":     status,
			"last_test":  s.getLastTestTime(ctx, sender.Id),
		})
	}
	return results, nil
}

func (s *InboxPlacementService) getLastTestTime(ctx context.Context, senderPoolId int) *time.Time {
	var lastSent time.Time
	err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_pool_id", senderPoolId).
		Order("sent_at DESC").
		Value("sent_at", &lastSent)
	if err != nil {
		return nil
	}
	return &lastSent
}