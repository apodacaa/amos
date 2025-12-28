package models

import (
	"encoding/json"
	"time"
)

// Entry represents a journal entry
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TodoIDs   []string  `json:"todo_ids,omitempty"` // IDs of todos created in this entry

	// DEPRECATED: For backward compatibility during JSON migration
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// GetTags returns the tags for the entry (implements Taggable interface)
func (e Entry) GetTags() []string {
	return e.Tags
}

// GetTimestamp returns CreatedAt (implements Timestamped interface for filtering)
// This ensures date filters use creation date, not edit date
func (e Entry) GetTimestamp() time.Time {
	return e.CreatedAt
}

// UnmarshalJSON implements custom JSON unmarshaling for backward compatibility
// Handles migration from old "timestamp" field to new "created_at"/"updated_at" fields
func (e *Entry) UnmarshalJSON(data []byte) error {
	// Create alias type to avoid recursion
	type Alias Entry
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Migration: if old format (timestamp exists but CreatedAt is zero)
	if !e.Timestamp.IsZero() && e.CreatedAt.IsZero() {
		e.CreatedAt = e.Timestamp
		e.UpdatedAt = e.Timestamp
		e.Timestamp = time.Time{} // Clear deprecated field (omitempty will exclude from JSON)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling to exclude deprecated Timestamp field
func (e Entry) MarshalJSON() ([]byte, error) {
	// Create a type without the Timestamp field for marshaling
	type Alias struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		Tags      []string  `json:"tags"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		TodoIDs   []string  `json:"todo_ids,omitempty"`
	}

	return json.Marshal(&Alias{
		ID:        e.ID,
		Title:     e.Title,
		Body:      e.Body,
		Tags:      e.Tags,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		TodoIDs:   e.TodoIDs,
	})
}
