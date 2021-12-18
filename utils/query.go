package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func QueryErrorHandler(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("can't find any record")
	}
	return err
}

func UpdateResultHandler(result *gorm.DB, returnVal interface{}) (interface{}, error) {
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		err = QueryErrorHandler(err)
		return nil, err
	}
	if ra > 0 {
		return returnVal, nil
	} else {
		return nil, fmt.Errorf("no record affected")
	}
}

func DeleteResultHandler(result *gorm.DB) (bool, error) {
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		err = QueryErrorHandler(err)
		return false, err
	}
	if ra > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("no record affected")
	}
}
