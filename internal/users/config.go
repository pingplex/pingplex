package users

import "time"

// Config holds the configuration for the users module.
type Config struct {
	// CacheTTL defines the time-to-live for cached user data.
	// This controls how long user data is kept in Redis cache before
	// being refreshed from the database.
	CacheTTL time.Duration

	// UsernameMaxLength defines the maximum allowed length for usernames.
	UsernameMaxLength int

	// UsernameMinLength defines the minimum allowed length for usernames.
	UsernameMinLength int
}
