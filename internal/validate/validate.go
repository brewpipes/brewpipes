package validate

import (
	"fmt"
	"strings"
)

// Required returns an error if value is empty or whitespace-only.
func Required(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}
