package sequence

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/service/abtest"
)

// GetAbTestSubjectAndBody checks if an AB test is active for the step
// and returns the appropriate variant content for the given enrollment.
// If no AB test exists, returns the original subject and empty body (use template).
func GetAbTestSubjectAndBody(ctx context.Context, sequenceId, stepId int, enrollmentId int, email string, defaultSubject string) (subject, bodyHtml, bodyText string, hasAbTest bool) {
	test, err := abtest.AbTest().GetActiveTestForStep(ctx, sequenceId, stepId)
	if err != nil || test == nil {
		return defaultSubject, "", "", false
	}

	// If test has a winner, use winner's content for all new sends
	if test.Status == 2 && test.WinnerVariant >= 0 {
		subj, html, txt, err := abtest.AbTest().GetVariantContent(ctx, test.Id, test.WinnerVariant)
		if err != nil {
			g.Log().Warningf(ctx, "AB test %d: failed to get winner variant %d: %v", test.Id, test.WinnerVariant, err)
			return defaultSubject, "", "", false
		}
		return subj, html, txt, true
	}

	// Assign variant to this enrollment
	variant, err := abtest.AbTest().GetVariantForRecipient(ctx, test.Id, email, enrollmentId)
	if err != nil {
		g.Log().Warningf(ctx, "AB test %d: failed to assign variant for %s: %v", test.Id, email, err)
		return defaultSubject, "", "", false
	}

	subj, html, txt, err := abtest.AbTest().GetVariantContent(ctx, test.Id, variant)
	if err != nil {
		g.Log().Warningf(ctx, "AB test %d: failed to get variant %d content: %v", test.Id, variant, err)
		return defaultSubject, "", "", false
	}

	// Track sent
	_ = abtest.AbTest().UpdateVariantStats(ctx, test.Id, variant, "sent")

	return subj, html, txt, true
}

// RecordAbTestEvent records an engagement event for AB test tracking
func RecordAbTestEvent(ctx context.Context, sequenceId, stepId int, email string, eventType string) {
	test, err := abtest.AbTest().GetActiveTestForStep(ctx, sequenceId, stepId)
	if err != nil || test == nil {
		return
	}

	// Find the assignment for this email
	var assignment struct {
		Variant int `json:"variant"`
	}
	err = g.DB().Model("bm_ab_test_assignments a").
		Join("bm_sequence_enrollments e", "e.id = a.enrollment_id").
		Where("a.ab_test_id = ? AND e.email = ?", test.Id, email).
		Fields("a.variant").
		Scan(&assignment)
	if err != nil || assignment.Variant < 0 {
		return
	}

	err = abtest.AbTest().UpdateVariantStats(ctx, test.Id, assignment.Variant, eventType)
	if err != nil {
		g.Log().Warningf(ctx, "AB test %d: failed to update stats for variant %d, event %s: %v", test.Id, assignment.Variant, eventType, err)
	}
}

// GetAbTestInfoForStep returns AB test info for display purposes
func GetAbTestInfoForStep(ctx context.Context, sequenceId, stepId int) (hasTest bool, testId int, status int, winner int) {
	test, err := abtest.AbTest().GetActiveTestForStep(ctx, sequenceId, stepId)
	if err != nil || test == nil {
		return false, 0, 0, -1
	}
	return true, test.Id, test.Status, test.WinnerVariant
}

// BuildAbTestTaskNameSuffix creates a task name suffix indicating AB variant
func BuildAbTestTaskNameSuffix(ctx context.Context, sequenceId, stepId int, taskId int64) string {
	test, err := abtest.AbTest().GetActiveTestForStep(ctx, sequenceId, stepId)
	if err != nil || test == nil {
		return ""
	}
	return fmt.Sprintf("_ab%d", test.Id)
}
