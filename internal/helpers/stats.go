package helpers

import (
	"fmt"
	"sort"
	"time"

	"github.com/apodacaa/amos/internal/models"
)

// WeekStats holds statistics for a single week
type WeekStats struct {
	Year       int
	Week       int
	WeekLabel  string // Format: "W43"
	EntryCount int
	TodoCount  int
}

// GetISOWeek returns the ISO 8601 week number and year for a given time
func GetISOWeek(t time.Time) (year, week int) {
	year, week = t.ISOWeek()
	return year, week
}

// AggregateByWeek aggregates entries and todos by ISO 8601 week
// Returns the last N weeks of statistics, sorted from oldest to newest
func AggregateByWeek(entries []models.Entry, todos []models.Todo, lastNWeeks int) []WeekStats {
	// Create map to track counts by year-week
	weekMap := make(map[string]*WeekStats)

	// Aggregate entries
	for _, entry := range entries {
		year, week := GetISOWeek(entry.Timestamp)
		key := weekKey(year, week)

		if weekMap[key] == nil {
			weekMap[key] = &WeekStats{
				Year:      year,
				Week:      week,
				WeekLabel: weekLabel(week),
			}
		}
		weekMap[key].EntryCount++
	}

	// Aggregate todos
	for _, todo := range todos {
		year, week := GetISOWeek(todo.CreatedAt)
		key := weekKey(year, week)

		if weekMap[key] == nil {
			weekMap[key] = &WeekStats{
				Year:      year,
				Week:      week,
				WeekLabel: weekLabel(week),
			}
		}
		weekMap[key].TodoCount++
	}

	// Convert map to slice
	var weeks []WeekStats
	for _, stats := range weekMap {
		weeks = append(weeks, *stats)
	}

	// Sort by year and week (oldest first)
	sort.Slice(weeks, func(i, j int) bool {
		if weeks[i].Year != weeks[j].Year {
			return weeks[i].Year < weeks[j].Year
		}
		return weeks[i].Week < weeks[j].Week
	})

	// Get last N weeks, filling in gaps with zero counts
	return getLastNWeeks(weeks, lastNWeeks)
}

// weekKey generates a unique key for year-week combination
func weekKey(year, week int) string {
	return fmt.Sprintf("%d-W%02d", year, week)
}

// weekLabel formats week number as "W43"
func weekLabel(week int) string {
	return fmt.Sprintf("W%02d", week)
}

// getLastNWeeks returns the last N weeks, filling gaps with zero counts
func getLastNWeeks(weeks []WeekStats, n int) []WeekStats {
	if len(weeks) == 0 {
		// No data, return current week with zero counts
		now := time.Now()
		year, week := GetISOWeek(now)
		return []WeekStats{{
			Year:       year,
			Week:       week,
			WeekLabel:  weekLabel(week),
			EntryCount: 0,
			TodoCount:  0,
		}}
	}

	// Get the most recent week
	latestWeek := weeks[len(weeks)-1]
	latestYear := latestWeek.Year
	latestWeekNum := latestWeek.Week

	// Build last N weeks including gaps
	result := make([]WeekStats, 0, n)
	weekMap := make(map[string]WeekStats)

	for _, w := range weeks {
		weekMap[weekKey(w.Year, w.Week)] = w
	}

	// Walk backwards from latest week
	currentYear := latestYear
	currentWeek := latestWeekNum

	for i := 0; i < n; i++ {
		key := weekKey(currentYear, currentWeek)

		if stats, exists := weekMap[key]; exists {
			result = append([]WeekStats{stats}, result...)
		} else {
			// Fill gap with zero counts
			result = append([]WeekStats{{
				Year:       currentYear,
				Week:       currentWeek,
				WeekLabel:  weekLabel(currentWeek),
				EntryCount: 0,
				TodoCount:  0,
			}}, result...)
		}

		// Move to previous week
		currentWeek--
		if currentWeek < 1 {
			currentYear--
			// ISO 8601: most years have 52 weeks, some have 53
			dec31 := time.Date(currentYear, 12, 31, 0, 0, 0, 0, time.UTC)
			_, weeksInYear := GetISOWeek(dec31)
			currentWeek = weeksInYear
		}
	}

	return result
}
