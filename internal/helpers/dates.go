package helpers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/models"
)

// Timestamped is an interface for items that have timestamps
type Timestamped interface {
	GetTimestamp() time.Time
}

// filterByDateRange is a generic function that filters items by date range
// Returns filtered list or original list if date string is empty or invalid
func filterByDateRange[T Timestamped](items []T, dateFilter string) []T {
	if dateFilter == "" {
		return items
	}

	start, end := ParseDateFilter(dateFilter)
	if start.IsZero() || end.IsZero() {
		return items
	}

	filtered := []T{}
	for _, item := range items {
		ts := item.GetTimestamp()
		if (ts.Equal(start) || ts.After(start)) &&
			(ts.Equal(end) || ts.Before(end)) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// ParseDateFilter parses date filter strings and returns start/end times
// Supported formats:
// - "today" - today's entries (00:00:00 to 23:59:59)
// - "yesterday" - yesterday's entries (00:00:00 to 23:59:59)
// - "last N days" - N days including today (e.g., "last 14 days" = 13 days ago through today)
// - "YYYY-MM-DD" - single date (00:00:00 to 23:59:59)
// - "YYYY-MM-DD to YYYY-MM-DD" - date range, both dates inclusive
func ParseDateFilter(input string) (start time.Time, end time.Time) {
	if input == "" {
		return time.Time{}, time.Time{}
	}

	input = strings.ToLower(strings.TrimSpace(input))
	now := time.Now()
	loc := now.Location()

	// Start and end of today
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)

	// Check for "today"
	if input == "today" {
		return startOfToday, endOfToday
	}

	// Check for "yesterday"
	if input == "yesterday" {
		startOfYesterday := startOfToday.AddDate(0, 0, -1)
		endOfYesterday := endOfToday.AddDate(0, 0, -1)
		return startOfYesterday, endOfYesterday
	}

	// Check for "last N days" (e.g., "last 14 days")
	// Design decision: "last N days" is INCLUSIVE of today and counts backwards.
	// Example: "last 7 days" on Monday Jan 8 includes Mon Jan 2 through Mon Jan 8 (7 days total).
	// Why: This matches user expectations - "last 7 days" should include today.
	// Implementation: N days total = go back (N-1) days from start of today, through end of today.
	lastDaysRegex := regexp.MustCompile(`^last\s+(\d+)\s+days?$`)
	if matches := lastDaysRegex.FindStringSubmatch(input); len(matches) > 1 {
		n, err := strconv.Atoi(matches[1])
		if err == nil && n > 0 {
			// N days including today = go back (N-1) days from start of today
			start = startOfToday.AddDate(0, 0, -(n - 1))
			end = endOfToday
			return start, end
		}
	}

	// Check for date range "YYYY-MM-DD to YYYY-MM-DD"
	if strings.Contains(input, " to ") {
		parts := strings.Split(input, " to ")
		if len(parts) == 2 {
			startDate, errStart := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[0]), loc)
			endDate, errEnd := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[1]), loc)

			if errStart == nil && errEnd == nil {
				// Both dates inclusive: start at 00:00:00, end at 23:59:59
				start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
				end = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)
				return start, end
			}
		}
	}

	// Check for single date "YYYY-MM-DD"
	singleDate, err := time.ParseInLocation("2006-01-02", input, loc)
	if err == nil {
		// Single date = that day only (00:00:00 to 23:59:59)
		start = time.Date(singleDate.Year(), singleDate.Month(), singleDate.Day(), 0, 0, 0, 0, loc)
		end = time.Date(singleDate.Year(), singleDate.Month(), singleDate.Day(), 23, 59, 59, 999999999, loc)
		return start, end
	}

	// Unrecognized format
	return time.Time{}, time.Time{}
}

// FilterEntriesByDateRange filters entries by date string
// Returns filtered list or original list if date string is empty or invalid
func FilterEntriesByDateRange(entries []models.Entry, dateFilter string) []models.Entry {
	return filterByDateRange(entries, dateFilter)
}

// FilterTodosByDateRange filters todos by date string
// Returns filtered list or original list if date string is empty or invalid
func FilterTodosByDateRange(todos []models.Todo, dateFilter string) []models.Todo {
	return filterByDateRange(todos, dateFilter)
}
