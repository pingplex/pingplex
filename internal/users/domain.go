package users

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	// UserStatusPending indicates the user account is pending verification.
	UserStatusPending UserStatus = "pending"
	// UserStatusActive indicates the user account is active and can log in.
	UserStatusActive UserStatus = "active"
	// UserStatusSuspended indicates the user account is suspended.
	UserStatusSuspended UserStatus = "suspended"
)

// UserRole represents the role of a user in the system.
type UserRole string

const (
	// UserRoleUser is the default role for regular users.
	UserRoleUser UserRole = "user"
	// UserRoleAdmin is the role for administrators with elevated privileges.
	UserRoleAdmin UserRole = "admin"
)

// User represents a user domain entity.
type User struct {
	// ID is the unique identifier for the user.
	ID uuid.UUID

	// Email is the user's email address.
	Email string

	// Username is the unique username for the user.
	Username string

	// PasswordHash is the hashed password.
	PasswordHash string

	// Status is the current status of the user account.
	Status UserStatus

	// Role is the user's role in the system.
	Role UserRole

	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time

	// UpdatedAt is the timestamp when the user was last updated.
	UpdatedAt time.Time

	// LastLogin is the timestamp of the last successful login.
	LastLogin time.Time

	// MFAEnabled indicates whether multi-factor authentication is enabled.
	MFAEnabled bool
}

// NewUser creates a new user with default values.
func NewUser(id uuid.UUID, email, username, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Status:       UserStatusPending,
		Role:         UserRoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
		MFAEnabled:   false,
	}
}

// IsActive returns true if the user account is active.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsValidStatus checks if the given status is valid.
func IsValidStatus(status UserStatus) bool {
	switch status {
	case UserStatusPending, UserStatusActive, UserStatusSuspended:
		return true
	default:
		return false
	}
}

// IsValidRole checks if the given role is valid.
func IsValidRole(role UserRole) bool {
	switch role {
	case UserRoleUser, UserRoleAdmin:
		return true
	default:
		return false
	}
}
