package test

import (
	"context"
	"testing"
	"time"

	"billionmail-core/internal/service/deliverability"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestInboxPlacementFullFlow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.GetIOContext()

		// 1. Add a seed
		svc := deliverability.NewSeedListService()
		testEmail := "test@gmail.com"
		err := svc.AddSeed(ctx, testEmail, "gmail")
		t.AssertNil(err)

		// 2. List seeds
		seeds, err := svc.GetSeeds(ctx)
		t.AssertNil(err)
		t.AssertGT(len(seeds), 0)

		// 3. Trigger test
		placementSvc := deliverability.NewInboxPlacementService()
		err = placementSvc.RunPlacementTests(ctx)
		t.AssertNil(err)

		// 4. Check results table has entry
		count, err := g.DB().Model("bm_inbox_placement_results").Count()
		t.AssertNil(err)
		t.AssertGT(count, 0)

		// 5. Remove seed
		err = svc.RemoveSeed(ctx, testEmail)
		t.AssertNil(err)
	})
}

func TestSeedListCRUD(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.GetIOContext()
		svc := deliverability.NewSeedListService()
		testEmail := "crud_test@example.com"

		// Clean up any existing seed
		_ = svc.RemoveSeed(ctx, testEmail)

		// Add
		err := svc.AddSeed(ctx, testEmail, "outlook")
		t.AssertNil(err)

		// Read
		seeds, err := svc.GetSeeds(ctx)
		t.AssertNil(err)
		found := false
		for _, s := range seeds {
			if s.Email == testEmail {
				found = true
				t.AssertEQ(s.Provider, "outlook")
			}
		}
		t.AssertEQ(found, true)

		// Remove
		err = svc.RemoveSeed(ctx, testEmail)
		t.AssertNil(err)

		// Verify removed
		seeds, err = svc.GetSeeds(ctx)
		t.AssertNil(err)
		for _, s := range seeds {
			if s.Email == testEmail {
				t.Fail()
			}
		}
	})
}

func TestInboxPlacementScoreCalculation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.GetIOContext()

		// Create a test sender entry first via direct DB insert
		now := time.Now()

		// Insert test sender (if not exists)
		var senderId int
		err := g.DB().Model("bm_sender_pool").Data(g.Map{
			"username":      "test_score_sender",
			"email":         "test_score@example.com",
			"status":        "active",
			"daily_limit":   100,
			"monthly_limit": 1000,
		}).Insert()
		t.AssertNil(err)

		// Get the sender ID
	row := g.DB().Model("bm_sender_pool").Where("username", "test_score_sender").Scan(&struct {
		Id int `json:"id"`
	}{})
		if row != nil {
			senderId = row.Id
		}

		// Skip if no sender available
		if senderId == 0 {
			t.Log("No active sender available, skipping score calculation test")
			return
		}

		// Insert test placement results
		for i := 0; i < 5; i++ {
			g.DB().Model("bm_inbox_placement_results").Data(g.Map{
				"sender_pool_id": senderId,
				"seed_email":     "test_score_seed@example.com",
				"test_email_message_id": "test-score-msg-" + time.Now().Format("20060102150405"),
				"sent_at":        now.Add(-time.Duration(i) * 24 * time.Hour),
				"placement":     "inbox",
			}).Insert()
		}

		// Calculate score
		placementSvc := deliverability.NewInboxPlacementService()
		rate, status, err := placementSvc.CalculateScore(ctx, senderId, 30)
		t.AssertNil(err)
		t.AssertGE(rate, 0)
		t.AssertNE(status, "")

		// Cleanup test data
		g.DB().Model("bm_inbox_placement_results").Where("sender_pool_id", senderId).Delete()
		g.DB().Model("bm_sender_pool").Where("id", senderId).Delete()
	})
}
