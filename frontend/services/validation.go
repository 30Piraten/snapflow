package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/30Piraten/snapflow/utils"
)

// ValidateOrder checks if all required fields are present and valid
func ValidateOrder(order *utils.PhotoOrder) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	var missingFields []string

	if strings.TrimSpace(order.FullName) == "" {
		missingFields = append(missingFields, "full name")
	}
	if strings.TrimSpace(order.Email) == "" {
		missingFields = append(missingFields, "email")
	}
	if strings.TrimSpace(order.Location) == "" {
		missingFields = append(missingFields, "location")
	}
	if strings.TrimSpace(order.Size) == "" {
		missingFields = append(missingFields, "size")
	}
	if strings.TrimSpace(order.PaperType) == "" {
		missingFields = append(missingFields, "paper type")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Basic email validation
	if !strings.Contains(order.Email, "@") || !strings.Contains(order.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}
