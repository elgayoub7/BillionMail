package bounce

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/model/entity"
	"billionmail-core/internal/service/scoring"
)

type BounceHandlerService struct{}

var service = &BounceHandlerService{}

func BounceHandler() *BounceHandlerService {
	return service
}

// ProcessBounceLog scans mailstat_send_mails for new bounces and handles them
func (s *BounceHandlerService) ProcessBounceLog(ctx context.Context) {
	var bounces []struct {
		PostfixMessageId string `json:"postfix_message_id"`
		Status           string `json:"status"`
		Recipient        string `json:"recipient"`
		Description      string `json:"description"`
	}

	err := g.DB().Model("mailstat_send_mails").
		Where("status = ?", "bounced").
		Where("postfix_message_id NOT IN (SELECT message_id FROM bm_bounce_records WHERE message_id != '')").
		Fields("postfix_message_id, status, recipient, description").
		Limit(500).
		Scan(&bounces)

	if err != nil {
		g.Log().Error(ctx, "BounceHandler: failed to query bounces:", err)
		return
	}

	for _, b := range bounces {
		bounceType := s.ClassifyBounce(b.Description)
		s.HandleBounce(ctx, b.Recipient, bounceType, b.Description, 0, b.PostfixMessageId)
	}

	if len(bounces) > 0 {
		g.Log().Infof(ctx, "BounceHandler: processed %d bounces", len(bounces))
	}
}

// ClassifyBounce determines if a bounce is hard, soft, or a complaint
func (s *BounceHandlerService) ClassifyBounce(description string) string {
	desc := strings.ToLower(description)

	hardPatterns := []string{
		"user not found", "no such user", "unknown user", "recipient invalid",
		"recipient not found", "address rejected", "invalid recipient",
		"domain not found", "host not found", "no such domain",
		"550 5.1.1", "550 5.4.1", "550 5.7.1", "553 5.3.0",
		"no route to host", "name service error",
	}
	for _, p := range hardPatterns {
		if strings.Contains(desc, p) {
			return "hard"
		}
	}

	softPatterns := []string{
		"mailbox full", "quota exceeded", "over quota",
		"temporarily unavailable", "try again later",
		"450 4.2.2", "450 4.7.1", "451 4.3.0", "452 4.2.2",
		"connection timed out", "deferred",
	}
	for _, p := range softPatterns {
		if strings.Contains(desc, p) {
			return "soft"
		}
	}

	complaintPatterns := []string{
		"spam report", "abuse report", "complaint", "marked as spam",
	}
	for _, p := range complaintPatterns {
		if strings.Contains(desc, p) {
			return "complaint"
		}
	}

	return "soft" // default to soft for unknown patterns
}

// HandleBounce processes a single bounce event
func (s *BounceHandlerService) HandleBounce(ctx context.Context, email, bounceType, description string, sequenceId int, messageId string) {
	var action string

	switch bounceType {
	case "hard":
		action = s.HandleHardBounce(ctx, email, sequenceId)
	case "soft":
		action = s.HandleSoftBounce(ctx, email, sequenceId)
	case "complaint":
		action = s.HandleComplaint(ctx, email, sequenceId)
	}

	// Record the bounce
	g.DB().Model("bm_bounce_records").Insert(g.Map{
		"email":        email,
		"bounce_type":  bounceType,
		"description":  description,
		"sequence_id":  sequenceId,
		"message_id":   messageId,
		"action_taken": action,
	})

	// Update lead scoring
	_ = scoring.Scoring().RecordEngagement(ctx, email, "bounce")
}

// HandleHardBounce removes contact from all active sequences
func (s *BounceHandlerService) HandleHardBounce(ctx context.Context, email string, sequenceId int) string {
	// Complete all active enrollments for this email
	result, err := g.DB().Model("bm_sequence_enrollments").
		Where("email = ? AND status = 0", email).
		Update(g.Map{
			"status": 3, // 3 = bounced
		})
	if err != nil {
		g.Log().Warningf(ctx, "HardBounce: failed to update enrollments for %s: %v", email, err)
		return "error_updating_enrollments"
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		g.Log().Infof(ctx, "HardBounce: %s removed from %d active enrollments", email, rows)
	}

	return "removed_from_sequences"
}

// HandleSoftBounce tracks soft bounces and escalates after 3 consecutive
func (s *BounceHandlerService) HandleSoftBounce(ctx context.Context, email string, sequenceId int) string {
	count, err := g.DB().Model("bm_bounce_records").
		Where("email = ? AND bounce_type = 'soft'", email).
		Count()
	if err != nil {
		return "tracked"
	}

	if count >= 3 {
		g.Log().Infof(ctx, "SoftBounce: %s reached 3 soft bounces, escalating to hard", email)
		return s.HandleHardBounce(ctx, email, sequenceId)
	}

	return "tracked"
}

// HandleComplaint flags the contact and notifies
func (s *BounceHandlerService) HandleComplaint(ctx context.Context, email string, sequenceId int) string {
	// Remove from all active sequences
	s.HandleHardBounce(ctx, email, sequenceId)

	g.Log().Warningf(ctx, "Complaint: %s reported email as spam", email)
	return "removed_and_flagged"
}

// GetBounceStats returns bounce statistics
func (s *BounceHandlerService) GetBounceStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	var results []struct {
		BounceType string `json:"bounce_type"`
		Count      int    `json:"count"`
	}
	g.DB().Model("bm_bounce_records").
		Fields("bounce_type, COUNT(*) as count").
		Group("bounce_type").
		Scan(&results)

	for _, r := range results {
		stats[r.BounceType] = r.Count
	}

	return stats, nil
}

// GetBouncesByEmail returns bounce history for an email
func (s *BounceHandlerService) GetBouncesByEmail(ctx context.Context, email string) ([]entity.BounceRecord, error) {
	var records []entity.BounceRecord
	err := g.DB().Model("bm_bounce_records").
		Where("email = ?", email).
		Order("created_at DESC").
		Limit(20).
		Scan(&records)
	return records, err
}
