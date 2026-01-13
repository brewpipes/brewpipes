package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/identity/storage"
	"github.com/gofrs/uuid/v5"
)

// User is the DTO for queried user objects.
type User struct {
	ID uuid.UUID

	Username string `json:"username"`
	Role     string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

// NewUsersResponse creates a new UsersResponse from a slice of users.
func NewUsersResponse(entities []storage.User) UsersResponse {
	collection := newUserCollection(entities)
	return UsersResponse{
		Users: collection,
	}
}

func newUser(u storage.User) User {
	return User{
		ID:        u.UUID,
		Username:  u.Username,
		Role:      "user",
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func newUserCollection(entities []storage.User) []User {
	users := []User{}
	for _, l := range entities {
		users = append(users, newUser(l))
	}

	return users
}
