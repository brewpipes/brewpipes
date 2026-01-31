package storage

import (
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func uuidParam(value *uuid.UUID) any {
	if value == nil {
		return nil
	}

	return *value
}

func assignUUIDPointer(destination **uuid.UUID, value pgtype.UUID) {
	if value.Valid {
		uuidValue := uuid.UUID(value.Bytes)
		*destination = &uuidValue
		return
	}

	*destination = nil
}
