package deliverability

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/os/gtimer"
	"github.com/google/uuid"
)

type InboxPlacementService struct {
	seedList *SeedListService
	ctx      context.Context
}

func NewInboxPlacementService() *InboxPlacementService {
	return &InboxPlacementService{
		seedList: NewSeedListService(),
		ctx:      context.Background(),
	}
}

func (s *InboxPlacementService) StartTimer() {
	gtimer.Add(15*time.Minute, func() {
		s.RunPlacementTests(s.ctx)
	})
}

func (s *InboxPlacementService) RunPlacementTests(ctx context.Context) error {
	seeds, err := s.seedList.GetActiveSeeds(ctx)
	if err != nil || len(seeds) == 0 {
		return err
	}

	var senders []struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	err = g.DB().Model("bm_sender_config").
		Where("is_active", 1).
		Scan(&senders)
	if err != nil {
		return err
	}

	for _, sender := range senders {
		s.sendTestEmailsForSender(ctx, sender.Id, sender.Email, seeds)
	}
	return nil
}

func (s *InboxPlacementService) sendTestEmailsForSender(ctx context.Context, senderId int, senderEmail string, seeds []SeedEmail) error {
	for _, seed := range seeds {
		messageId := fmt.Sprintf("BM-INBOX-TEST-%s", uuid.New().String())

		_, err := g.DB().Model("bm_inbox_placement_results").
			Data(map[string]interface{}{
				"sender_config_id":     senderId,
				"seed_email":           seed.Email,
				"test_email_message_id": messageId,
				"sent_at":              time.Now(),
				"placement":            "none",
			}).
			Insert()
		if err != nil {
			continue
		}

		fmt.Printf("[InboxPlacement] TEST sender=%s seed=%s message_id=%s\n",
			senderEmail, seed.Email, messageId)
	}
	return nil
}

func (s *InboxPlacementService) UpdatePlacementFromReply(ctx context.Context, messageId string, receivedAt time.Time) error {
	if !strings.HasPrefix(messageId, "BM-INBOX-TEST-") {
		return nil
	}

	_, err := g.DB().Model("bm_inbox_placement_results").
		Where("test_email_message_id", messageId).
		Data(map[string]interface{}{
			"received_at": receivedAt,
			"placement":   "inbox",
		}).
		Update()
	return err
}

func (s *InboxPlacementService) CalculateScore(ctx context.Context, senderConfigId int, windowDays int) (int, string, error) {
	cutoff := time.Now().AddDate(0, 0, -windowDays)

	total, err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_config_id", senderConfigId).
		Where("sent_at > ?", cutoff).
		Count()
	if err != nil {
		return 0, "unknown", err
	}
	if total == 0 {
		return 0, "unknown", nil
	}

	inboxCount, err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_config_id", senderConfigId).
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

func (s *InboxPlacementService) GetAllScores(ctx context.Context) ([]map[string]interface{}, error) {
	var senders []struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	err := g.DB().Model("bm_sender_config").
		Where("is_active", 1).
		Scan(&senders)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for _, sender := range senders {
		rate, status, _ := s.CalculateScore(ctx, sender.Id, 30)
		results = append(results, map[string]interface{}{
			"sender_id":  sender.Id,
			"email":      sender.Email,
			"inbox_rate": rate,
			"status":     status,
			"last_test":  s.getLastTestTime(ctx, sender.Id),
		})
	}
	return results, nil
}

func (s *InboxPlacementService) getLastTestTime(ctx context.Context, senderConfigId int) *time.Time {
	var lastSent time.Time
	_, err := g.DB().Model("bm_inbox_placement_results").
		Where("sender_config_id", senderConfigId).
		Order("sent_at DESC").
		Value("sent_at", &lastSent)
	if err != nil {
		return nil
	}
	return &lastSent
}
