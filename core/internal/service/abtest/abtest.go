package abtest

import (
	"context"
	"fmt"
	"hash/fnv"
	"math"

	"github.com/gogf/gf/v2/frame/g"

	v1 "billionmail-core/internal/model/entity"
)

type AbTestService struct{}

var service = &AbTestService{}

func AbTest() *AbTestService {
	return service
}

// CreateAbTest creates a new A/B test for a sequence step
func (s *AbTestService) CreateAbTest(ctx context.Context, sequenceId, stepId int, name string, splitRatio float64, criteria string, variantA, variantB v1.AbTestVariant) (int, error) {
	if splitRatio <= 0 || splitRatio >= 1 {
		splitRatio = 0.5
	}
	if criteria == "" {
		criteria = "open_rate"
	}

	count, err := g.DB().Model("bm_ab_tests").Where("sequence_id = ? AND step_id = ? AND status IN (0,1)", sequenceId, stepId).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, fmt.Errorf("an active AB test already exists for this step")
	}

	// Insert the test
	result, err := g.DB().Model("bm_ab_tests").Insert(g.Map{
		"sequence_id":           sequenceId,
		"step_id":               stepId,
		"name":                  name,
		"status":                1, // running
		"winner_variant":        -1,
		"winner_criteria":       criteria,
		"split_ratio":           splitRatio,
		"confidence_threshold":  0.95,
	})
	if err != nil {
		return 0, err
	}

	testId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert variant A
	_, err = g.DB().Model("bm_ab_test_variants").Insert(g.Map{
		"ab_test_id": testId,
		"variant":    0,
		"subject":    variantA.Subject,
		"body_html":  variantA.BodyHtml,
		"body_text":  variantA.BodyText,
	})
	if err != nil {
		return 0, err
	}

	// Insert variant B
	_, err = g.DB().Model("bm_ab_test_variants").Insert(g.Map{
		"ab_test_id": testId,
		"variant":    1,
		"subject":    variantB.Subject,
		"body_html":  variantB.BodyHtml,
		"body_text":  variantB.BodyText,
	})
	if err != nil {
		return 0, err
	}

	g.Log().Infof(ctx, "AB Test created: id=%d, sequence=%d, step=%d", testId, sequenceId, stepId)
	return int(testId), nil
}

// GetVariantForRecipient assigns a variant using deterministic hash of email
func (s *AbTestService) GetVariantForRecipient(ctx context.Context, abTestId int, email string, enrollmentId int) (int, error) {
	var assignment v1.AbTestAssignment
	err := g.DB().Model("bm_ab_test_assignments").Where("ab_test_id = ? AND enrollment_id = ?", abTestId, enrollmentId).Scan(&assignment)
	if err == nil && assignment.Id > 0 {
		return assignment.Variant, nil
	}

	var test v1.AbTest
	err = g.DB().Model("bm_ab_tests").Where("id = ?", abTestId).Scan(&test)
	if err != nil {
		return 0, err
	}

	// Deterministic assignment using hash
	h := fnv.New32a()
	h.Write([]byte(email))
	hashVal := h.Sum32()

	var variant int
	if float64(hashVal%100) < test.SplitRatio*100 {
		variant = 0 // A
	} else {
		variant = 1 // B
	}

	// Save assignment
	_, err = g.DB().Model("bm_ab_test_assignments").Insert(g.Map{
		"ab_test_id":     abTestId,
		"enrollment_id":  enrollmentId,
		"variant":        variant,
	})
	if err != nil {
		return 0, err
	}

	return variant, nil
}

// GetActiveTestForStep returns the active AB test for a given sequence step
func (s *AbTestService) GetActiveTestForStep(ctx context.Context, sequenceId, stepId int) (*v1.AbTest, error) {
	var test v1.AbTest
	err := g.DB().Model("bm_ab_tests").Where("sequence_id = ? AND step_id = ? AND status = 1", sequenceId, stepId).Scan(&test)
	if err != nil {
		return nil, err
	}
	if test.Id == 0 {
		return nil, nil
	}
	return &test, nil
}

// GetTestVariants returns both variants for a test
func (s *AbTestService) GetTestVariants(ctx context.Context, abTestId int) ([]v1.AbTestVariant, error) {
	var variants []v1.AbTestVariant
	err := g.DB().Model("bm_ab_test_variants").Where("ab_test_id = ?", abTestId).Order("variant ASC").Scan(&variants)
	if err != nil {
		return nil, err
	}
	return variants, nil
}

// GetVariantContent returns the subject and body for a specific variant
func (s *AbTestService) GetVariantContent(ctx context.Context, abTestId, variant int) (subject, bodyHtml, bodyText string, err error) {
	var v v1.AbTestVariant
	err = g.DB().Model("bm_ab_test_variants").Where("ab_test_id = ? AND variant = ?", abTestId, variant).Scan(&v)
	if err != nil {
		return "", "", "", err
	}
	return v.Subject, v.BodyHtml, v.BodyText, nil
}

// UpdateVariantStats increments the stats for a variant
func (s *AbTestService) UpdateVariantStats(ctx context.Context, abTestId, variant int, eventType string) error {
	var field string
	switch eventType {
	case "sent":
		field = "sent_count"
	case "open":
		field = "open_count"
	case "click":
		field = "click_count"
	case "reply":
		field = "reply_count"
	case "bounce":
		field = "bounce_count"
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	_, err := g.DB().Exec(ctx, fmt.Sprintf(
		"UPDATE bm_ab_test_variants SET %s = %s + 1 WHERE ab_test_id = ? AND variant = ?",
		field, field,
	), abTestId, variant)

	return err
}

// CheckWinner evaluates if a statistically significant winner exists
func (s *AbTestService) CheckWinner(ctx context.Context, abTestId int) (bool, int, error) {
	test, err := s.GetTestById(ctx, abTestId)
	if err != nil || test == nil {
		return false, -1, err
	}
	if test.Status != 1 {
		return false, -1, nil
	}

	variants, err := s.GetTestVariants(ctx, abTestId)
	if err != nil || len(variants) < 2 {
		return false, -1, err
	}

	vA := variants[0]
	vB := variants[1]

	// Need minimum 50 sends per variant for statistical significance
	if vA.SentCount < 50 || vB.SentCount < 50 {
		return false, -1, nil
	}

	// Calculate rates based on winner criteria
	rateA := s.calcRate(test.WinnerCriteria, &vA)
	rateB := s.calcRate(test.WinnerCriteria, &vB)

	// Simple z-test for proportions
	significant := s.isSignificant(rateA, float64(vA.SentCount), rateB, float64(vB.SentCount), test.ConfidenceThreshold)

	if significant {
		winner := 0
		if rateB > rateA {
			winner = 1
		}
		// Auto-assign winner
		_, err = g.DB().Model("bm_ab_tests").Where("id = ?", abTestId).Update(g.Map{
			"winner_variant": winner,
			"status":         2,
			"completed_at":   g.DB().GetCore().Time().Unix(),
		})
		if err != nil {
			return false, -1, err
		}
		g.Log().Infof(ctx, "AB Test winner auto-selected: test=%d, winner=%c", abTestId, 'A'+winner)
		return true, winner, nil
	}

	return false, -1, nil
}

// PickWinner manually selects a winner
func (s *AbTestService) PickWinner(ctx context.Context, abTestId, winnerVariant int) error {
	_, err := g.DB().Model("bm_ab_tests").Where("id = ?", abTestId).Update(g.Map{
		"winner_variant": winnerVariant,
		"status":         2,
		"completed_at":   g.DB().GetCore().Time().Unix(),
	})
	return err
}

// GetTestById retrieves a test by ID
func (s *AbTestService) GetTestById(ctx context.Context, abTestId int) (*v1.AbTest, error) {
	var test v1.AbTest
	err := g.DB().Model("bm_ab_tests").Where("id = ?", abTestId).Scan(&test)
	if err != nil {
		return nil, err
	}
	if test.Id == 0 {
		return nil, nil
	}
	return &test, nil
}

// GetTestsForSequence returns all tests for a sequence
func (s *AbTestService) GetTestsForSequence(ctx context.Context, sequenceId int) ([]v1.AbTest, error) {
	var tests []v1.AbTest
	err := g.DB().Model("bm_ab_tests").Where("sequence_id = ?", sequenceId).Order("id DESC").Scan(&tests)
	return tests, err
}

// DeleteAbTest deletes an AB test and its variants
func (s *AbTestService) DeleteAbTest(ctx context.Context, abTestId int) error {
	_, err := g.DB().Model("bm_ab_tests").Where("id = ?", abTestId).Delete()
	return err
}

// GetTestResults returns full results for a test
func (s *AbTestService) GetTestResults(ctx context.Context, abTestId int) (*v1.AbTestResult, error) {
	test, err := s.GetTestById(ctx, abTestId)
	if err != nil || test == nil {
		return nil, err
	}

	variants, err := s.GetTestVariants(ctx, abTestId)
	if err != nil || len(variants) < 2 {
		return nil, fmt.Errorf("missing variants")
	}

	result := &v1.AbTestResult{
		TestId:      test.Id,
		TestName:    test.Name,
		Status:      test.Status,
		Winner:      test.WinnerVariant,
		VariantA:    s.buildVariantStats(&variants[0]),
		VariantB:    s.buildVariantStats(&variants[1]),
	}

	rateA := s.calcRate(test.WinnerCriteria, &variants[0])
	rateB := s.calcRate(test.WinnerCriteria, &variants[1])
	result.IsSignificant = s.isSignificant(rateA, float64(variants[0].SentCount), rateB, float64(variants[1].SentCount), test.ConfidenceThreshold)

	return result, nil
}

// ProcessRunningTests checks all running tests for auto-winner
func (s *AbTestService) ProcessRunningTests(ctx context.Context) {
	var tests []v1.AbTest
	err := g.DB().Model("bm_ab_tests").Where("status = 1").Scan(&tests)
	if err != nil {
		g.Log().Error(ctx, "Failed to fetch running AB tests:", err)
		return
	}

	for _, test := range tests {
		found, winner, err := s.CheckWinner(ctx, test.Id)
		if err != nil {
			g.Log().Warningf(ctx, "Error checking AB test %d: %v", test.Id, err)
			continue
		}
		if found {
			g.Log().Infof(ctx, "AB test %d auto-winner: variant %d", test.Id, winner)
		}
	}
}

// --- Helpers ---

func (s *AbTestService) calcRate(criteria string, v *v1.AbTestVariant) float64 {
	if v.SentCount == 0 {
		return 0
	}
	switch criteria {
	case "open_rate":
		return float64(v.OpenCount) / float64(v.SentCount)
	case "click_rate":
		return float64(v.ClickCount) / float64(v.SentCount)
	case "reply_rate":
		return float64(v.ReplyCount) / float64(v.SentCount)
	default:
		return float64(v.OpenCount) / float64(v.SentCount)
	}
}

func (s *AbTestService) buildVariantStats(v *v1.AbTestVariant) *v1.VariantStats {
	label := "A"
	if v.Variant == 1 {
		label = "B"
	}
	stats := &v1.VariantStats{
		VariantId:  v.Id,
		Label:      label,
		SentCount:  v.SentCount,
		BounceRate: 0,
	}
	if v.SentCount > 0 {
		stats.OpenRate = float64(v.OpenCount) / float64(v.SentCount) * 100
		stats.ClickRate = float64(v.ClickCount) / float64(v.SentCount) * 100
		stats.ReplyRate = float64(v.ReplyCount) / float64(v.SentCount) * 100
		stats.BounceRate = float64(v.BounceCount) / float64(v.SentCount) * 100
	}
	return stats
}

// isSignificant performs a two-proportion z-test
func (s *AbTestService) isSignificant(p1, n1, p2, n2 float64, threshold float64) bool {
	if n1 == 0 || n2 == 0 {
		return false
	}

	// Combined proportion
	pCombined := (p1*n1 + p2*n2) / (n1 + n2)
	if pCombined == 0 || pCombined == 1 {
		return false
	}

	// Standard error
	se := math.Sqrt(pCombined*(1-pCombined)*(1/n1 + 1/n2))
	if se == 0 {
		return false
	}

	// Z-score
	z := (p1 - p2) / se
	zAbs := math.Abs(z)

	// Critical value for 95% confidence (two-tailed) ≈ 1.96
	// For higher thresholds, scale accordingly
	criticalValue := 1.96
	if threshold > 0.95 {
		criticalValue = 2.576 // 99% confidence
	}

	return zAbs > criticalValue
}
