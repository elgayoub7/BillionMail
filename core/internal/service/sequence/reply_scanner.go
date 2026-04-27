package sequence

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// ScanReplyInboxes scans sender maildirs for replies to sequence emails
func ScanReplyInboxes(ctx context.Context) {
	senders, err := getActiveSenderAddresses(ctx)
	if err != nil {
		g.Log().Warning(ctx, "ScanReplyInboxes: failed to get senders:", err)
		return
	}
	if len(senders) == 0 {
		return
	}

	for _, sender := range senders {
		scanMaildir(ctx, sender)
	}
}

// getActiveSenderAddresses returns distinct sender addresses from running sequences
func getActiveSenderAddresses(ctx context.Context) ([]string, error) {
	result, err := g.DB().Model("bm_sequences").
		Where("status = 1").
		Fields("DISTINCT addresser").
		Array()
	if err != nil {
		return nil, err
	}

	var senders []string
	for _, v := range result {
		s := strings.TrimSpace(v.String())
		if s != "" {
			senders = append(senders, s)
		}
	}
	return senders, nil
}

// scanMaildir reads the sender's Maildir for new/unseen messages and checks In-Reply-To
func scanMaildir(ctx context.Context, emailAddr string) {
	// Parse emailAddr -> domain/local_part
	parts := strings.SplitN(emailAddr, "@", 2)
	if len(parts) != 2 {
		return
	}
	localPart := parts[0]
	domain := parts[1]

	// Maildir base path: /opt/billionmail/vmail-data/<domain>/<local>/
	maildirBase := filepath.Join("/opt/billionmail/vmail-data", domain, localPart)

	// Check cur/ and new/ directories
	for _, subdir := range []string{"cur", "new"} {
		dir := filepath.Join(maildirBase, subdir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Check if already processed (flag :2, contains 'R' = replied or 'S' = seen)
			name := entry.Name()
			if strings.Contains(name, ":2,") {
				flags := name[strings.Index(name, ":2,")+3:]
				if strings.Contains(flags, "S") || strings.Contains(flags, "R") {
					continue
				}
			}

			// Read message headers
			inReplyTo, err := readInReplyTo(filepath.Join(dir, name))
			if err != nil || inReplyTo == "" {
				continue
			}

			cleanID := strings.Trim(inReplyTo, " <>")
			if cleanID == "" {
				continue
			}

			// Match against bm_sequence_email_tasks
			enrollmentId, err := g.DB().Model("bm_sequence_email_tasks").
				Where("message_id = ?", cleanID).
				Fields("enrollment_id").
				Value()
			if err != nil || enrollmentId.Int() == 0 {
				continue
			}

			// Update enrollment replies counter
			_, err = g.DB().Exec(ctx,
				"UPDATE bm_sequence_enrollments SET total_replies = total_replies + 1 WHERE id = ?",
				enrollmentId.Int())
			if err != nil {
				g.Log().Errorf(ctx, "ScanReplyInboxes: failed to update replies for enrollment %d: %v", enrollmentId.Int(), err)
			} else {
				g.Log().Infof(ctx, "ScanReplyInboxes: detected reply for enrollment %d (in-reply-to: %s)", enrollmentId.Int(), cleanID)

				// Also increment step replied_count
				stepId, err2 := g.DB().Model("bm_sequence_email_tasks").
					Where("message_id = ?", cleanID).
					Fields("step_id").
					Value()
				if err2 == nil && stepId.Int() > 0 {
					_, _ = g.DB().Exec(ctx,
						"UPDATE bm_sequence_steps SET replied_count = replied_count + 1 WHERE id = ?",
						stepId.Int())
				}
			}

			// Mark as seen by renaming with :S flag
			markSeen(dir, name)
		}
	}
}

// readInReplyTo reads an email file and extracts the In-Reply-To header
func readInReplyTo(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Empty line = end of headers
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "In-Reply-To:") {
			value := strings.TrimPrefix(line, "In-Reply-To:")
			return strings.TrimSpace(value), nil
		}
	}
	return "", nil
}

// markSeen renames a Maildir file to add the :S (seen) flag
func markSeen(dir, name string) {
	// If already has flags, add S
	if idx := strings.Index(name, ":2,"); idx >= 0 {
		flags := name[idx+3:]
		if !strings.Contains(flags, "S") {
			newName := name[:idx+3] + "S" + flags
			os.Rename(filepath.Join(dir, name), filepath.Join(dir, newName))
		}
	} else {
		// No flags yet, add :2,S
		os.Rename(filepath.Join(dir, name), filepath.Join(dir, name+":2,S"))
	}
}
