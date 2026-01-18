package users

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when trying to create a user that already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUsernameTaken is returned when the requested username is already in use.
	ErrUsernameTaken = errors.New("username is already taken")

	// ErrEmailTaken is returned when the requested email is already in use.
	ErrEmailTaken = errors.New("email is already taken")

	// ErrInvalidUserID is returned when the provided user ID is not valid.
	ErrInvalidUserID = errors.New("invalid user ID")

	// ErrInvalidUsername is returned when the provided username does not meet requirements.
	ErrInvalidUsername = errors.New("invalid username")

	// ErrInvalidEmail is returned when the provided email is not valid.
	ErrInvalidEmail = errors.New("invalid email")

	// ErrInvalidStatus is returned when the provided status is not valid.
	ErrInvalidStatus = errors.New("invalid user status")

	// ErrInvalidRole is returned when the provided role is not valid.
	ErrInvalidRole = errors.New("invalid user role")

	// ErrDatabaseError is returned when a database operation fails.
	ErrDatabaseError = errors.New("database error")

	// ErrCacheError is returned when a cache operation fails.
	ErrCacheError = errors.New("cache error")
)
