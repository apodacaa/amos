package models

import "time"

// Todo represents a task item
type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"` // "open", "next", or "done"
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	EntryID   *string   `json:"entry_id,omitempty"` // Pointer - nil if standalone
}

// GetTags returns the tags for the todo (implements Taggable interface)
func (t Todo) GetTags() []string {
	return t.Tags
}

// GetTimestamp returns the creation time for the todo (implements Timestamped interface)
func (t Todo) GetTimestamp() time.Time {
	return t.CreatedAt
}

// GetUpdatedTimestamp returns the update time for the todo (implements Timestamped interface)
// Date filters use OR logic: match by CreatedAt OR UpdatedAt
func (t Todo) GetUpdatedTimestamp() time.Time {
	return t.UpdatedAt
}
