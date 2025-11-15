package models

// Config holds user preferences for the app
type Config struct {
	Theme string `json:"theme"` // Theme name (e.g., "cyberpunk", "brutalist")
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Theme: "cyberpunk", // Default to cyberpunk theme
	}
}
