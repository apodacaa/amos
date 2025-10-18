package models

import "time"

// Entry represents a journal entry
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
	TodoRefs  []string  `json:"todo_refs,omitempty"` // References to related todos (future)
}
