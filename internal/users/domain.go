package users

import "time"

type Status string
type Provider string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"

	ProviderTelegram Provider = "telegram"
)

type User struct {
	ID     string
	Status Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Identity struct {
	Provider     Provider
	ProviderID   string
	ProviderData string
}
