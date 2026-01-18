package users

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service implements the business logic for the users module.
type Service struct {
	// config holds the configuration for the service
	config Config

	// users provides access to data operations
	users *Repository

	// metrics is used for collecting and reporting metrics
	metrics *Metrics

	// logger is used for structured logging
	logger *zap.Logger
}

// New creates and initializes a new Service instance.
func New(config Config, users *Repository, metrics *Metrics, logger *zap.Logger) *Service {
	return &Service{
		config:  config,
		users:   users,
		metrics: metrics,
		logger:  logger,
	}
}

// CreateUser creates a new user with the given credentials.
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.ObserveCreateUserDuration(time.Since(start).Seconds())
	}()

	// Validate input
	if req.Email == "" {
		return nil, ErrInvalidEmail
	}
	if req.Username == "" {
		return nil, ErrInvalidUsername
	}
	if len(req.Username) < s.config.UsernameMinLength || len(req.Username) > s.config.UsernameMaxLength {
		return nil, fmt.Errorf("username must be between %d and %d characters: %w",
			s.config.UsernameMinLength, s.config.UsernameMaxLength, ErrInvalidUsername)
	}

	// Check if username is already taken
	exists, err := s.users.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("failed to check username existence", zap.Error(err))
		return nil, ErrDatabaseError
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// Check if email is already registered
	exists, err = s.users.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("failed to check email existence", zap.Error(err))
		return nil, ErrDatabaseError
	}
	if exists {
		return nil, ErrEmailTaken
	}

	// Create user
	id := uuid.New()
	user := NewUser(id, req.Email, req.Username, req.Password)

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	s.metrics.IncTotalUsers()
	s.logger.Info("user created", zap.String("id", id.String()), zap.String("email", req.Email))

	return user.toResponse(), nil
}

// GetUser retrieves a user by their ID.
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.ObserveGetUserDuration(time.Since(start).Seconds())
	}()

	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user.toResponse(), nil
}

// GetUserByEmail retrieves a user by their email address.
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.ObserveGetUserDuration(time.Since(start).Seconds())
	}()

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user.toResponse(), nil
}

// GetUserByUsername retrieves a user by their username.
func (s *Service) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.ObserveGetUserDuration(time.Since(start).Seconds())
	}()

	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user.toResponse(), nil
}

// UpdateUser updates an existing user with the given data.
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*UserResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.ObserveUpdateUserDuration(time.Since(start).Seconds())
	}()

	// Get existing user
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Username != nil {
		// Check if new username is already taken by another user
		if *req.Username != user.Username {
			exists, err := s.users.ExistsByUsername(ctx, *req.Username)
			if err != nil {
				s.logger.Error("failed to check username existence", zap.Error(err))
				return nil, ErrDatabaseError
			}
			if exists {
				return nil, ErrUsernameTaken
			}
		}
		user.Username = *req.Username
	}

	if req.Status != nil {
		if !IsValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		user.Status = *req.Status
	}

	if req.Role != nil {
		if !IsValidRole(*req.Role) {
			return nil, ErrInvalidRole
		}
		user.Role = *req.Role
	}

	// Save changes
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}

	s.logger.Info("user updated", zap.String("id", id.String()))

	return user.toResponse(), nil
}

// DeleteUser removes a user from the system.
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		s.metrics.ObserveDeleteUserDuration(time.Since(start).Seconds())
	}()

	// Check if user exists
	_, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.users.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info("user deleted", zap.String("id", id.String()))

	return nil
}

// UpdateLastLogin updates the last login timestamp for a user.
func (s *Service) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if err := s.users.UpdateLastLogin(ctx, id); err != nil {
		return err
	}

	s.logger.Debug("user last login updated", zap.String("id", id.String()))

	return nil
}

// CheckUsernameAvailability checks if a username is available.
func (s *Service) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	exists, err := s.users.ExistsByUsername(ctx, username)
	if err != nil {
		return false, ErrDatabaseError
	}
	return !exists, nil
}

// CheckEmailAvailability checks if an email is available for registration.
func (s *Service) CheckEmailAvailability(ctx context.Context, email string) (bool, error) {
	exists, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return false, ErrDatabaseError
	}
	return !exists, nil
}

// ActivateUser sets a user's status to active.
func (s *Service) ActivateUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Status = UserStatusActive
	if err := s.users.Update(ctx, user); err != nil {
		return err
	}

	s.metrics.IncActiveUsers()

	return nil
}

// SuspendUser sets a user's status to suspended.
func (s *Service) SuspendUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Status = UserStatusSuspended
	if err := s.users.Update(ctx, user); err != nil {
		return err
	}

	s.metrics.DecActiveUsers()

	return nil
}
