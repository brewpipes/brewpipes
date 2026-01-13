package entity

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Identifiers is a set of identifiers for an entity.
type Identifiers struct {
	ID   int64
	UUID uuid.UUID
}

// Timestamps is a set of timestamps for an entity.
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
