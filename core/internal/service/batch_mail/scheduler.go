package batch_mail

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/internal/model/entity"
)

// ScheduleChecker handles time window and day-of-week scheduling
type ScheduleChecker struct{}

// NewScheduleChecker creates a new schedule checker
func NewScheduleChecker() *ScheduleChecker {
	return &ScheduleChecker{}
}

// IsWithinSchedule checks if current time is within the allowed schedule
func (sc *ScheduleChecker) IsWithinSchedule(task *entity.EmailTask) bool {
	now := time.Now()

	// Check time window (schedule_start_hour and schedule_end_hour)
	if task.ScheduleStartHour > 0 || task.ScheduleEndHour < 24 {
		hour := now.Hour()
		if task.ScheduleStartHour == task.ScheduleEndHour {
			// Both 0 or both same = no restriction
		} else if hour < task.ScheduleStartHour || hour >= task.ScheduleEndHour {
			return false
		}
	}

	// Check day of week
	if task.ScheduleDays != "" && task.ScheduleDays != "[1,2,3,4,5,6,7]" {
		var allowedDays []int
		if err := json.Unmarshal([]byte(task.ScheduleDays), &allowedDays); err == nil {
			// Convert Go weekday (Sunday=0) to ISO (Monday=1, Sunday=7)
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7 // Sunday = 7 in ISO
			}

			found := false
			for _, d := range allowedDays {
				if d == weekday {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

// SleepUntilNextWindow blocks until the next allowed schedule window
// Returns true if sleep completed, false if context was canceled
func (sc *ScheduleChecker) SleepUntilNextWindow(ctx context.Context, task *entity.EmailTask) bool {
	nextWindow := sc.CalculateNextWindow(task)
	sleepDuration := time.Until(nextWindow)

	if sleepDuration <= 0 {
		return true
	}

	g.Log().Infof(ctx, "Task %d: outside schedule window (hours %d-%d, days %s), sleeping until %s (%.1f hours)",
		task.Id,
		task.ScheduleStartHour, task.ScheduleEndHour,
		task.ScheduleDays,
		nextWindow.Format("2006-01-02 15:04:05"),
		sleepDuration.Hours(),
	)

	// Sleep in chunks to be responsive to context cancellation
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	deadline := time.After(sleepDuration)
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "Task %d: schedule sleep interrupted by context cancellation", task.Id)
			return false
		case <-ticker.C:
			// Check if schedule changed (task config reload)
			if sc.IsWithinSchedule(task) {
				g.Log().Infof(ctx, "Task %d: schedule now allows sending, waking up", task.Id)
				return true
			}
		case <-deadline:
			g.Log().Infof(ctx, "Task %d: schedule window opened, resuming", task.Id)
			return true
		}
	}
}

// CalculateNextWindow finds the next time when sending is allowed
func (sc *ScheduleChecker) CalculateNextWindow(task *entity.EmailTask) time.Time {
	now := time.Now()

	// Parse allowed days
	allowedDays := []int{1, 2, 3, 4, 5, 6, 7} // default: all days
	if task.ScheduleDays != "" {
		if err := json.Unmarshal([]byte(task.ScheduleDays), &allowedDays); err != nil {
			allowedDays = []int{1, 2, 3, 4, 5, 6, 7}
		}
	}

	startHour := task.ScheduleStartHour
	endHour := task.ScheduleEndHour

	// Default: no restriction
	if startHour == 0 && endHour == 24 {
		startHour = 0
	}

	// Try each day starting from today
	for dayOffset := 0; dayOffset <= 7; dayOffset++ {
		candidate := now.AddDate(0, 0, dayOffset)

		// Check if this day is allowed
		weekday := int(candidate.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		dayAllowed := false
		for _, d := range allowedDays {
			if d == weekday {
				dayAllowed = true
				break
			}
		}

		if !dayAllowed {
			continue
		}

		// Calculate the start time for this day
		targetTime := time.Date(
			candidate.Year(), candidate.Month(), candidate.Day(),
			startHour, 0, 0, 0, now.Location(),
		)

		// If this is today and we're already past startHour, try next window
		if dayOffset == 0 {
			if now.Hour() >= endHour && endHour > startHour {
				// Past end hour today, try next allowed day
				continue
			}
			if now.Hour() >= startHour {
				// Already within window or past start
				// This shouldn't happen if IsWithinSchedule returned false
				continue
			}
			// Before start hour today — return start hour today
			return targetTime
		}

		// Future day — return start hour on that day
		return targetTime
	}

	// Fallback: try again in 1 hour
	return now.Add(1 * time.Hour)
}
