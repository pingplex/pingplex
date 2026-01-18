package users

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// userModel represents a data model used for database operations.
type userModel struct {
	// ID is the unique identifier for the user.
	ID uuid.UUID `db:"id"`

	// Email is the user's email address.
	Email string `db:"email"`

	// EmailNormalized is the lowercase version of the email for case-insensitive lookups.
	EmailNormalized string `db:"email_normalized"`

	// Username is the unique username for the user.
	Username string `db:"username"`

	// UsernameNormalized is the lowercase version of the username for case-insensitive lookups.
	UsernameNormalized string `db:"username_normalized"`

	// PasswordHash is the hashed password.
	PasswordHash string `db:"password_hash"`

	// Status is the current status of the user account.
	Status string `db:"status"`

	// Role is the user's role in the system.
	Role string `db:"role"`

	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt is the timestamp when the user was last updated.
	UpdatedAt time.Time `db:"updated_at"`

	// LastLogin is the timestamp of the last successful login.
	LastLogin time.Time `db:"last_login"`

	// MFAEnabled indicates whether multi-factor authentication is enabled.
	MFAEnabled bool `db:"mfa_enabled"`
}

// toDomain converts a userModel to a User domain entity.
func (m *userModel) toDomain() *User {
	return &User{
		ID:           m.ID,
		Email:        m.Email,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Status:       UserStatus(m.Status),
		Role:         UserRole(m.Role),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		LastLogin:    m.LastLogin,
		MFAEnabled:   m.MFAEnabled,
	}
}

// fromDomain converts a User domain entity to a userModel.
func fromDomain(u *User) *userModel {
	return &userModel{
		ID:                 u.ID,
		Email:              u.Email,
		EmailNormalized:    normalizeEmail(u.Email),
		Username:           u.Username,
		UsernameNormalized: normalizeUsername(u.Username),
		PasswordHash:       u.PasswordHash,
		Status:             string(u.Status),
		Role:               string(u.Role),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		LastLogin:          u.LastLogin,
		MFAEnabled:         u.MFAEnabled,
	}
}

// CreateUserRequest represents a request to create a new user.
type CreateUserRequest struct {
	// Email is the user's email address.
	Email string `json:"email"`

	// Username is the unique username for the user.
	Username string `json:"username"`

	// Password is the plaintext password (will be hashed).
	Password string `json:"password"`
}

// UpdateUserRequest represents a request to update an existing user.
type UpdateUserRequest struct {
	// Username is the new username (optional).
	Username *string `json:"username,omitempty"`

	// Status is the new status (optional).
	Status *UserStatus `json:"status,omitempty"`

	// Role is the new role (optional).
	Role *UserRole `json:"role,omitempty"`
}

// UserResponse represents the response returned for user operations.
type UserResponse struct {
	// ID is the unique identifier for the user.
	ID string

	// Email is the user's email address.
	Email string

	// Username is the unique username for the user.
	Username string

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

// toResponse converts a User domain entity to a UserResponse.
func (u *User) toResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID.String(),
		Email:      u.Email,
		Username:   u.Username,
		Status:     u.Status,
		Role:       u.Role,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		LastLogin:  u.LastLogin,
		MFAEnabled: u.MFAEnabled,
	}
}

// normalizeEmail converts an email to lowercase for case-insensitive lookups.
func normalizeEmail(email string) string {
	// In a real implementation, this would properly handle email normalization
	// including handling of internationalized email addresses.
	return lowercaseAndTrim(email)
}

// normalizeUsername converts a username to lowercase for case-insensitive lookups.
func normalizeUsername(username string) string {
	return lowercaseAndTrim(username)
}

// lowercaseAndTrim converts a string to lowercase and trims whitespace.
func lowercaseAndTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
