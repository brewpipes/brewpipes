package uuidutil

import "github.com/gofrs/uuid/v5"

// ToStringPointer converts a *uuid.UUID to a *string, returning nil if the input is nil.
func ToStringPointer(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}

	str := value.String()
	return &str
}
