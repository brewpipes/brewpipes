package dto

import "github.com/gofrs/uuid/v5"

func uuidToStringPointer(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}

	str := value.String()
	return &str
}
