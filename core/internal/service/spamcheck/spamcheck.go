package spamcheck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

t"time"
	"github.com/gogf/gf/v2/frame/g"
)

type SpamCheckService struct{}

var service = &SpamCheckService{}

func SpamCheck() *SpamCheckService {
	return service
}

type RspamdRequest struct {
	From    string `json:"from,omitempty"`
	To      string `json:"rcpt,omitempty"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
	Helo    string `json:"helo,omitempty"`
	IP      string `json:"ip,omitempty"`
}

type RspamdSymbol struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

type RspamdResponse struct {
	Score         float64       `json:"score"`
	RequiredScore float64       `json:"required_score"`
	Action        string        `json:"action"`
	Symbols       []RspamdSymbol `json:"symbols"`
	MessageID     string        `json:"message-id"`
}

type SpamResult struct {
	Score        float64  `json:"score"`
	Action       string   `json:"action"`
	IsBlocked    bool     `json:"is_blocked"`
	Symbols      []string `json:"symbols"`
	Suggestions  []string `json:"suggestions"`
}

// CheckEmailContent sends content to Rspamd for spam scoring
func (s *SpamCheckService) CheckEmailContent(ctx context.Context, from, subject, bodyHtml, bodyText string) (*SpamResult, error) {
	body := bodyHtml
	if body == "" {
		body = bodyText
	}

	reqBody := RspamdRequest{
		From:    from,
		Subject: subject,
		Body:    body,
		Helo:    "mail.selectcircle.fr",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	host := "127.0.0.1:11333"
	if g.Cfg().MustGet(ctx, "server.address").String() != "" {
		host = "rspamd:11333"
	}

	url := fmt.Sprintf("http://%s/checkv2", host)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rspamd connection failed: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rspamdResp RspamdResponse
	if err := json.Unmarshal(bodyResp, &rspamdResp); err != nil {
		return nil, fmt.Errorf("rspamd response parse error: %w", err)
	}

	result := &SpamResult{
		Score:   rspamdResp.Score,
		Action:  rspamdResp.Action,
		Symbols: make([]string, 0),
	}

	for _, sym := range rspamdResp.Symbols {
		result.Symbols = append(result.Symbols, sym.Name)
	}

	// Scoring thresholds
	result.IsBlocked = rspamdResp.Score >= 6.0
	result.Suggestions = s.GetSuggestions(result.Symbols)

	// Save to DB
	s.SaveResult(ctx, 0, 0, result)

	return result, nil
}

// CheckBeforeSend validates email before sending, returns true if OK to send
func (s *SpamCheckService) CheckBeforeSend(ctx context.Context, taskId, stepId int, from, subject, bodyHtml, bodyText string) (bool, *SpamResult, error) {
	result, err := s.CheckEmailContent(ctx, from, subject, bodyHtml, bodyText)
	if err != nil {
		g.Log().Warningf(ctx, "Spam check failed for task %d: %v", taskId, err)
		return true, nil, nil // Allow sending if check fails
	}

	s.SaveResult(ctx, taskId, stepId, result)

	if result.IsBlocked {
		g.Log().Warningf(ctx, "Email blocked by spam check: task=%d, score=%.1f", taskId, result.Score)
		return false, result, nil
	}

	return true, result, nil
}

// SaveResult persists spam check result
func (s *SpamCheckService) SaveResult(ctx context.Context, taskId, stepId int, result *SpamResult) {
	now := time.Now().Unix()
	isBlocked := 0
	if result.IsBlocked {
		isBlocked = 1
	}

	symbols, _ := json.Marshal(result.Symbols)
	suggestions, _ := json.Marshal(result.Suggestions)

	g.DB().Model("bm_spam_scores").Insert(g.Map{
		"email_task_id":      taskId,
		"sequence_step_id":   stepId,
		"score":              result.Score,
		"action":             result.Action,
		"symbols":            string(symbols),
		"suggestions":        string(suggestions),
		"is_blocked":         isBlocked,
		"checked_at":         now,
	})
}

// GetSuggestions maps Rspamd symbols to human-readable suggestions
func (s *SpamCheckService) GetSuggestions(symbols []string) []string {
	suggestions := []string{}
	symbolMap := map[string]string{
		"BAYES_SPAM":       "Content looks spammy — rewrite with more natural language",
		"R_DKIM_NA":        "DKIM signature is missing — configure DKIM for this domain",
		"R_SPF_FAIL":       "SPF record check failed — verify SPF configuration",
		"R_DKIM_REJECT":    "DKIM validation failed — check DKIM key and selector",
		"DMARC_POLICY_REJECT": "DMARC policy rejection — update DMARC record",
		"MISSING_FROM":     "From header is missing",
		"MISSING_SUBJECT":  "Subject is missing",
		"EMPTY_MESSAGE":    "Message body is empty",
		"HTML_IMAGE_ONLY":  "Email is image-only — add text content",
		"TRACKING_PIXEL":   "Tracking pixel detected — may trigger spam filters",
		"FORGED_RECIPIENTS":"Recipient address looks forged",
		"MANY_URLS":        "Too many URLs — reduce link count",
		"SUSPICIOUS_URL":   "Suspicious URL detected",
	}

	for _, sym := range symbols {
		if suggestion, ok := symbolMap[sym]; ok {
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

// GetRecentBlocked returns recently blocked emails
func (s *SpamCheckService) GetRecentBlocked(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := g.DB().Model("bm_spam_scores").
		Where("is_blocked = 1").
		Order("checked_at DESC").
		Limit(limit).
		Scan(&results)
	return results, err
}

// HasHighSpamSymbols checks if any critical spam symbols are present
func (s *SpamCheckService) HasHighSpamSymbols(symbols []string) bool {
	critical := map[string]bool{
		"R_SPF_FAIL": true, "R_DKIM_REJECT": true, "DMARC_POLICY_REJECT": true,
	}
	for _, sym := range symbols {
		if strings.HasPrefix(sym, "R_SPF_") || strings.HasPrefix(sym, "R_DKIM_") {
			return true
		}
		if critical[sym] {
			return true
		}
	}
	return false
}
