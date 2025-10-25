package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/google/uuid"
)

var (
	entries = flag.Int("entries", 100, "Number of entries to generate")
	todos   = flag.Int("todos", 50, "Number of todos to generate")
	path    = flag.String("path", "", "Custom output path (default: ~/.amos)")
)

var (
	sampleTitles = []string{
		"Meeting notes",
		"Project planning",
		"Weekly review",
		"Sprint retrospective",
		"Client feedback",
		"Design brainstorm",
		"Code review notes",
		"Bug investigation",
		"Feature research",
		"Team sync",
	}

	sampleTags = []string{
		"work", "personal", "urgent", "research", "meeting",
		"client", "design", "dev", "planning", "review",
	}

	sampleTodoTitles = []string{
		"Follow up on email",
		"Review pull request",
		"Update documentation",
		"Fix bug in authentication",
		"Deploy to staging",
		"Write test cases",
		"Schedule meeting",
		"Send project update",
		"Research new library",
		"Refactor component",
	}

	sampleBodyFragments = []string{
		"Discussed the upcoming features and roadmap.",
		"Need to prioritize the high-impact items first.",
		"The team agreed on the new architecture approach.",
		"Several blockers were identified and assigned.",
		"Action items documented for next sprint.",
		"Reviewed the user feedback and analytics.",
		"Key decisions made regarding tech stack.",
		"Timeline adjusted based on new requirements.",
		"Security concerns addressed in depth.",
		"Performance benchmarks look promising.",
	}
)

func main() {
	flag.Parse()

	// No need to seed - Go 1.20+ auto-seeds

	// Determine output path
	outputPath := *path
	if outputPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		outputPath = filepath.Join(home, ".amos")
	}

	// Ensure directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generating test data...\n")
	fmt.Printf("  Entries: %d\n", *entries)
	fmt.Printf("  Todos: %d\n", *todos)
	fmt.Printf("  Output: %s\n\n", outputPath)

	startTime := time.Now()

	// Generate entries
	generatedEntries := generateEntries(*entries)
	entriesPath := filepath.Join(outputPath, "entries.json")
	if err := saveJSON(entriesPath, generatedEntries); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving entries: %v\n", err)
		os.Exit(1)
	}

	// Generate todos (some linked to entries)
	generatedTodos := generateTodos(*todos, generatedEntries)
	todosPath := filepath.Join(outputPath, "todos.json")
	if err := saveJSON(todosPath, generatedTodos); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving todos: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)

	// Get file sizes
	entriesSize := getFileSize(entriesPath)
	todosSize := getFileSize(todosPath)

	fmt.Printf("✓ Generation complete in %v\n\n", elapsed)
	fmt.Printf("Results:\n")
	fmt.Printf("  entries.json: %d entries, %s\n", len(generatedEntries), formatBytes(entriesSize))
	fmt.Printf("  todos.json:   %d todos, %s\n", len(generatedTodos), formatBytes(todosSize))
	fmt.Printf("  Total size:   %s\n", formatBytes(entriesSize+todosSize))
	fmt.Printf("\nRun 'make run' to test performance with this data.\n")
}

func generateEntries(count int) []models.Entry {
	entries := make([]models.Entry, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		// Generate timestamp (spread over past year)
		daysAgo := rand.Intn(365)
		timestamp := now.AddDate(0, 0, -daysAgo)

		// Random title
		title := fmt.Sprintf("%s %03d", randomChoice(sampleTitles), i+1)

		// Generate body with random tags
		tags := randomTags(2, 5)
		bodyParts := make([]string, rand.Intn(3)+2) // 2-4 fragments
		for j := range bodyParts {
			bodyParts[j] = randomChoice(sampleBodyFragments)
		}
		body := ""
		for _, part := range bodyParts {
			body += part + " "
		}
		// Add tags to body
		for _, tag := range tags {
			body += "@" + tag + " "
		}

		entries[i] = models.Entry{
			ID:        uuid.New().String(),
			Title:     title,
			Body:      body,
			Tags:      tags,
			Timestamp: timestamp,
			TodoIDs:   []string{}, // Will be populated when we link todos
		}
	}

	return entries
}

func generateTodos(count int, entries []models.Entry) []models.Todo {
	todos := make([]models.Todo, count)
	now := time.Now()

	// Decide how many are linked vs standalone (60% linked, 40% standalone)
	linkedCount := int(float64(count) * 0.6)

	for i := 0; i < count; i++ {
		daysAgo := rand.Intn(365)
		createdAt := now.AddDate(0, 0, -daysAgo)

		title := fmt.Sprintf("%s %03d", randomChoice(sampleTodoTitles), i+1)
		tags := randomTags(1, 3)

		// Random status
		statuses := []string{"open", "open", "open", "next", "done"} // Weight towards open
		status := randomChoice(statuses)

		// Link to entry if in linked range
		var entryID *string
		if i < linkedCount && len(entries) > 0 {
			// Pick a random entry to link to
			linkedEntry := &entries[rand.Intn(len(entries))]
			entryID = &linkedEntry.ID
			// Add this todo ID to the entry
			linkedEntry.TodoIDs = append(linkedEntry.TodoIDs, uuid.New().String())
		}

		todos[i] = models.Todo{
			ID:        uuid.New().String(),
			Title:     title,
			Status:    status,
			Tags:      tags,
			CreatedAt: createdAt,
			EntryID:   entryID,
		}
	}

	return todos
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func randomTags(min, max int) []string {
	count := min + rand.Intn(max-min+1)
	tags := make([]string, count)
	used := make(map[string]bool)

	for i := 0; i < count; i++ {
		// Ensure unique tags
		for {
			tag := randomChoice(sampleTags)
			if !used[tag] {
				tags[i] = tag
				used[tag] = true
				break
			}
		}
	}

	return tags
}

func saveJSON(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0644)
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
