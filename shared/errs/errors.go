package errs

import (
	"errors"
	"fmt"
)

var (
	errInternal       = errors.New("internal error")
	errEntityNotFound = errors.New("entity not found")
)

func ErrEntityNotFound(entity string) error {
	if entity == "" {
		return errEntityNotFound
	}
	return fmt.Errorf("%w: %s", errEntityNotFound, entity)
}

func IsEntityNotFoundErr(err error) bool {
	return errors.Is(err, errEntityNotFound)
}

func ErrInternal(inner error) error {
	if inner == nil {
		return errInternal
	}
	return fmt.Errorf("%w: %w", errInternal, inner)
}

func IsInternalErr(err error) bool {
	return errors.Is(err, errInternal)
}
