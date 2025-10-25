package helpers

import (
	"testing"
	"time"

	"github.com/apodacaa/amos/internal/models"
)

func TestGetISOWeek(t *testing.T) {
	tests := []struct {
		date         time.Time
		expectedYear int
		expectedWeek int
	}{
		{time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC), 2025, 43},
		{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), 2025, 1},
		{time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), 2025, 1}, // Last day of 2024 is in week 1 of 2025
	}

	for _, test := range tests {
		year, week := GetISOWeek(test.date)
		if year != test.expectedYear || week != test.expectedWeek {
			t.Errorf("GetISOWeek(%v) = %d, %d; want %d, %d",
				test.date, year, week, test.expectedYear, test.expectedWeek)
		}
	}
}

func TestWeekLabel(t *testing.T) {
	tests := []struct {
		week     int
		expected string
	}{
		{1, "W01"},
		{9, "W09"},
		{10, "W10"},
		{43, "W43"},
		{52, "W52"},
	}

	for _, test := range tests {
		result := weekLabel(test.week)
		if result != test.expected {
			t.Errorf("weekLabel(%d) = %s; want %s", test.week, result, test.expected)
		}
	}
}

func TestAggregateByWeek(t *testing.T) {
	now := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC) // Week 43
	lastWeek := now.AddDate(0, 0, -7)                     // Week 42
	twoWeeksAgo := now.AddDate(0, 0, -14)                 // Week 41

	entries := []models.Entry{
		{ID: "1", Timestamp: now},
		{ID: "2", Timestamp: now},
		{ID: "3", Timestamp: lastWeek},
		{ID: "4", Timestamp: twoWeeksAgo},
	}

	todos := []models.Todo{
		{ID: "1", CreatedAt: now},
		{ID: "2", CreatedAt: lastWeek},
		{ID: "3", CreatedAt: lastWeek},
	}

	stats := AggregateByWeek(entries, todos, 8)

	// Should have 8 weeks
	if len(stats) != 8 {
		t.Errorf("Expected 8 weeks, got %d", len(stats))
	}

	// Check most recent week (should be last in array)
	currentWeek := stats[len(stats)-1]
	if currentWeek.Week != 43 {
		t.Errorf("Expected current week to be 43, got %d", currentWeek.Week)
	}
	if currentWeek.EntryCount != 2 {
		t.Errorf("Expected 2 entries in current week, got %d", currentWeek.EntryCount)
	}
	if currentWeek.TodoCount != 1 {
		t.Errorf("Expected 1 todo in current week, got %d", currentWeek.TodoCount)
	}

	// Check week 42
	week42 := stats[len(stats)-2]
	if week42.Week != 42 {
		t.Errorf("Expected week 42, got %d", week42.Week)
	}
	if week42.EntryCount != 1 {
		t.Errorf("Expected 1 entry in week 42, got %d", week42.EntryCount)
	}
	if week42.TodoCount != 2 {
		t.Errorf("Expected 2 todos in week 42, got %d", week42.TodoCount)
	}

	// Check week 41
	week41 := stats[len(stats)-3]
	if week41.Week != 41 {
		t.Errorf("Expected week 41, got %d", week41.Week)
	}
	if week41.EntryCount != 1 {
		t.Errorf("Expected 1 entry in week 41, got %d", week41.EntryCount)
	}
	if week41.TodoCount != 0 {
		t.Errorf("Expected 0 todos in week 41, got %d", week41.TodoCount)
	}

	// Check that older weeks have zero counts
	week36 := stats[0]
	if week36.EntryCount != 0 || week36.TodoCount != 0 {
		t.Errorf("Expected week %d to have zero counts, got %d entries and %d todos",
			week36.Week, week36.EntryCount, week36.TodoCount)
	}
}

func TestAggregateByWeek_EmptyData(t *testing.T) {
	entries := []models.Entry{}
	todos := []models.Todo{}

	stats := AggregateByWeek(entries, todos, 8)

	// Should have 1 week (current week with zero counts)
	if len(stats) != 1 {
		t.Errorf("Expected 1 week for empty data, got %d", len(stats))
	}

	if stats[0].EntryCount != 0 || stats[0].TodoCount != 0 {
		t.Errorf("Expected zero counts for empty data, got %d entries and %d todos",
			stats[0].EntryCount, stats[0].TodoCount)
	}
}
