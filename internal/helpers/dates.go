package helpers

import (
	"time"

	"github.com/apodacaa/amos/internal/models"
)

// Date filter presets
const (
	DateFilterToday       = "TODAY"
	DateFilterYesterday   = "YESTERDAY"
	DateFilterLast7Days   = "LAST_7_DAYS"
	DateFilterLast30Days  = "LAST_30_DAYS"
	DateFilterLast60Days  = "LAST_60_DAYS"
	DateFilterLast90Days  = "LAST_90_DAYS"
	DateFilterLast365Days = "LAST_365_DAYS"
)

// GetDatePresets returns all available date filter presets in order
func GetDatePresets() []string {
	return []string{
		DateFilterToday,
		DateFilterYesterday,
		DateFilterLast7Days,
		DateFilterLast30Days,
		DateFilterLast60Days,
		DateFilterLast90Days,
		DateFilterLast365Days,
	}
}

// GetDateRange returns start and end times for a given preset
// End time is always "now", start time varies by preset
func GetDateRange(preset string) (start time.Time, end time.Time) {
	now := time.Now()
	end = now

	// Start of today (midnight)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch preset {
	case DateFilterToday:
		start = startOfToday
	case DateFilterYesterday:
		start = startOfToday.AddDate(0, 0, -1)
		end = startOfToday // End at start of today
	case DateFilterLast7Days:
		start = startOfToday.AddDate(0, 0, -7)
	case DateFilterLast30Days:
		start = startOfToday.AddDate(0, 0, -30)
	case DateFilterLast60Days:
		start = startOfToday.AddDate(0, 0, -60)
	case DateFilterLast90Days:
		start = startOfToday.AddDate(0, 0, -90)
	case DateFilterLast365Days:
		start = startOfToday.AddDate(0, 0, -365)
	default:
		// Unknown preset = no filtering
		return time.Time{}, now
	}

	return start, end
}

// FilterEntriesByDateRange filters entries by date preset
// Returns filtered list or original list if preset is empty
func FilterEntriesByDateRange(entries []models.Entry, preset string) []models.Entry {
	if preset == "" {
		return entries
	}

	start, end := GetDateRange(preset)
	if start.IsZero() {
		return entries
	}

	filtered := []models.Entry{}
	for _, entry := range entries {
		if (entry.Timestamp.Equal(start) || entry.Timestamp.After(start)) &&
			(entry.Timestamp.Equal(end) || entry.Timestamp.Before(end)) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

// FilterTodosByDateRange filters todos by date preset
// Returns filtered list or original list if preset is empty
func FilterTodosByDateRange(todos []models.Todo, preset string) []models.Todo {
	if preset == "" {
		return todos
	}

	start, end := GetDateRange(preset)
	if start.IsZero() {
		return todos
	}

	filtered := []models.Todo{}
	for _, todo := range todos {
		if (todo.CreatedAt.Equal(start) || todo.CreatedAt.After(start)) &&
			(todo.CreatedAt.Equal(end) || todo.CreatedAt.Before(end)) {
			filtered = append(filtered, todo)
		}
	}

	return filtered
}

// FormatDatePreset returns a human-readable label for a preset
func FormatDatePreset(preset string) string {
	switch preset {
	case DateFilterToday:
		return "today"
	case DateFilterYesterday:
		return "yesterday"
	case DateFilterLast7Days:
		return "last 7 days"
	case DateFilterLast30Days:
		return "last 30 days"
	case DateFilterLast60Days:
		return "last 60 days"
	case DateFilterLast90Days:
		return "last 90 days"
	case DateFilterLast365Days:
		return "last 365 days"
	default:
		return ""
	}
}
