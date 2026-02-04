package users

import "time"

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type User struct {
	ID     string
	Status Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Identity struct {
	Provider     string
	ProviderID   string
	ProviderData string
}
