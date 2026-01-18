package users

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v3"
)

// Repository handles data access operations for the users module.
type Repository struct {
	// scylla provides access to the ScyllaDB database.
	scylla gocqlx.Session
}

// NewRepository creates and initializes a new Repository instance.
func NewRepository(scylla gocqlx.Session) *Repository {
	return &Repository{
		scylla: scylla,
	}
}

// Create inserts a new user into the database.
func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (
		id, email, email_normalized, username, username_normalized,
		password_hash, status, role, created_at, updated_at,
		mfa_enabled
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := r.scylla.Query(query, nil).BindStruct(fromDomain(user)).Exec(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID retrieves a user by their ID from the database.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, email, email_normalized, username, username_normalized,
		password_hash, status, role, created_at, updated_at,
		last_login, mfa_enabled
		FROM users WHERE id = ? LIMIT 1`

	var user userModel
	if err := r.scylla.Query(query, nil).Bind(id).Get(&user); err != nil {
		return nil, ErrUserNotFound
	}

	return user.toDomain(), nil
}

// GetByEmail retrieves a user by their email address.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	normalizedEmail := normalizeEmail(email)

	query := `SELECT id FROM users WHERE email_normalized = ? LIMIT 1`

	var userID uuid.UUID
	if err := r.scylla.Query(query, nil).Bind(normalizedEmail).Get(&userID); err != nil {
		return nil, ErrUserNotFound
	}

	// Use GetByID
	return r.GetByID(ctx, userID)
}

// GetByUsername retrieves a user by their username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	normalizedUsername := normalizeUsername(username)

	query := `SELECT id FROM users WHERE username_normalized = ? LIMIT 1`

	var userID uuid.UUID
	if err := r.scylla.Query(query, nil).Bind(normalizedUsername).Get(&userID); err != nil {
		return nil, ErrUserNotFound
	}

	// Use GetByID
	return r.GetByID(ctx, userID)
}

// Update updates an existing user in the database.
func (r *Repository) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE users SET
		username = ?, username_normalized = ?, status = ?, role = ?,
		updated_at = ?, last_login = ?, mfa_enabled = ?
		WHERE id = ?`

	if err := r.scylla.Query(query, nil).BindStruct(fromDomain(user)).Exec(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// Delete removes a user from the database by ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`

	if err := r.scylla.Query(query, nil).Bind(id).Exec(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// ExistsByUsername checks if a username is already taken.
func (r *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	normalizedUsername := normalizeUsername(username)

	query := `SELECT id FROM users WHERE username_normalized = ? LIMIT 1`

	var userID uuid.UUID
	err := r.scylla.Query(query, nil).Bind(normalizedUsername).Get(&userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ExistsByEmail checks if an email is already registered.
func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	normalizedEmail := normalizeEmail(email)

	query := `SELECT id FROM users WHERE email_normalized = ? LIMIT 1`

	var userID uuid.UUID
	err := r.scylla.Query(query, nil).Bind(normalizedEmail).Get(&userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// UpdateLastLogin updates the last login timestamp for a user.
func (r *Repository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login = ?, updated_at = ? WHERE id = ?`

	now := time.Now()
	if err := r.scylla.Query(query, nil).Bind(now, now, id).Exec(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}
