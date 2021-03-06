// Package utils provides utility functions using in entire app
package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Error handler for ORM query
func HandleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("can't find any record")
	}
	return err
}

// Define return values for update, delete query
func ReturnBoolStateFromResult(result *gorm.DB) (bool, error) {
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		err = HandleQueryError(err)
		return false, err
	}
	if ra > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("no record affected")
	}
}
