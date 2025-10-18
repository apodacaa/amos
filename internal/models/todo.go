package models

import "time"

// Todo represents a task item
type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"` // "open" or "done"
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	EntryID   *string   `json:"entry_id,omitempty"` // Pointer - nil if standalone
}
