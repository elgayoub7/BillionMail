package batch_mail

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

type SubjectScore struct {
	Score    int      `json:"score"`
	Warnings []string `json:"warnings"`
}

var spamWords = []string{
	"free", "winner", "win", "urgent", "act now", "limited time",
	"buy direct", "click here", "order now", "risk free", "guarantee",
	"no obligation", "no risk", "special offer", "deal", "discount",
	"cash", "money", "earn", "income", "profit", "subscribe",
	"unsubscribe", "bonus", "prize", "congratulations", "credit",
	"loan", "debt", "affordable", "cheap", "save big", "million",
	"billion", "opportunity", "miracle", "hidden", "secret",
	"password", "verify", "account suspended", "action required",
	"100% free", "extra income", "make money", "online pharmacy",
	"overdue", "invoice attached", "weight loss",
}

var (
	urlRegex     = regexp.MustCompile(`(?i)https?://[^\s]+|www\.[^\s]+`)
	excessiveRe  = regexp.MustCompile(`[!?]{3,}`)
	allCapsRe    = regexp.MustCompile(`^[^a-z]*$`)
)

func ScoreSubject(subject string) *SubjectScore {
	result := &SubjectScore{Score: 100, Warnings: nil}
	s := strings.TrimSpace(subject)
	if s == "" {
		result.Score = 0
		result.Warnings = append(result.Warnings, "Subject is empty")
		return result
	}

	// ALL CAPS check
	letters := 0
	upper := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters++
			if unicode.IsUpper(r) {
				upper++
			}
		}
	}
	if letters > 0 && float64(upper)/float64(letters) > 0.7 {
		deduct(result, 15, "Subject is mostly ALL CAPS")
	}

	// Spam words
	spamHits := 0
	lower := strings.ToLower(s)
	for _, word := range spamWords {
		if strings.Contains(lower, word) {
			spamHits++
			if spamHits <= 3 {
				deduct(result, 10, "Spam word detected: \""+word+"\"")
			}
		}
	}
	if spamHits > 3 {
		extra := (spamHits - 3) * 5
		deduct(result, extra, "More spam words detected")
	}

	// Excessive punctuation
	matches := excessiveRe.FindAllString(s, -1)
	if len(matches) > 0 {
		deduct(result, 10, "Excessive punctuation (multiple ! or ?)")
	}

	// Length check
	if len(s) > 78 {
		over := len(s) - 78
		deduct(result, int(math.Min(float64(over), 10)), "Subject too long ("+fmt.Sprintf("%d", len(s))+" chars, max 78)")
	}

	// URLs
	if urlRegex.MatchString(s) {
		deduct(result, 15, "URL detected in subject")
	}

	// Repeated words (manual check since Go regexp does not support backreferences)
	words := strings.Fields(lower)
	seen := make(map[string]bool)
	for _, w := range words {
		clean := strings.Trim(w, ".,!?:;")
		if clean == "" {
			continue
		}
		if seen[clean] {
			deduct(result, 5, "Repeated word detected")
			break
		}
		seen[clean] = true
	}

	if result.Score < 0 {
		result.Score = 0
	}
	return result
}

func deduct(s *SubjectScore, points int, reason string) {
	s.Score -= points
	s.Warnings = append(s.Warnings, reason)
}
