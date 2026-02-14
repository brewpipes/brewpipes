package database

import (
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDParam converts a *uuid.UUID to a value suitable for use as a SQL parameter.
// Returns nil if the pointer is nil, otherwise returns the dereferenced UUID.
func UUIDParam(value *uuid.UUID) any {
	if value == nil {
		return nil
	}

	return *value
}

// AssignUUIDPointer scans a pgtype.UUID into a **uuid.UUID destination.
// Sets the destination to nil if the pgtype value is not valid.
func AssignUUIDPointer(destination **uuid.UUID, value pgtype.UUID) {
	if value.Valid {
		uuidValue := uuid.UUID(value.Bytes)
		*destination = &uuidValue
		return
	}

	*destination = nil
}
