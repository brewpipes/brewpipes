package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/identity/storage"
)

// CreateUserRequest is the request body [POST /users].
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the request.
func (r CreateUserRequest) Validate() error {
	if err := validateRequired(r.Username, "username"); err != nil {
		return err
	}

	return validateRequired(r.Password, "password")
}

// UpdateUserRequest is the request body [PUT /users/{id}].
type UpdateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// Validate validates the request.
func (r UpdateUserRequest) Validate() error {
	if r.Username == nil && r.Password == nil {
		return validateRequired("", "username or password")
	}
	if r.Username != nil {
		if err := validateRequired(*r.Username, "username"); err != nil {
			return err
		}
	}
	if r.Password != nil {
		if err := validateRequired(*r.Password, "password"); err != nil {
			return err
		}
	}

	return nil
}

// UserResponse is the DTO for queried user objects.
type UserResponse struct {
	ID       int64  `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Role     string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

// NewUserResponse creates a new UserResponse from a user.
func NewUserResponse(entity storage.User) UserResponse {
	return newUser(entity)
}

// NewUsersResponse creates a new UsersResponse from a slice of users.
func NewUsersResponse(entities []storage.User) UsersResponse {
	collection := newUserCollection(entities)
	return UsersResponse{
		Users: collection,
	}
}

func newUser(u storage.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		UUID:      u.UUID.String(),
		Username:  u.Username,
		Role:      "user",
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func newUserCollection(entities []storage.User) []UserResponse {
	users := []UserResponse{}
	for _, l := range entities {
		users = append(users, newUser(l))
	}

	return users
}
