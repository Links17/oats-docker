package errors

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotUpdate = errors.New("record not updated")
)

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsNotUpdate(err error) bool {
	return errors.Is(err, ErrRecordNotUpdate)
}
