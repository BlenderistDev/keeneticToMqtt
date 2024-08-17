package errs

import "errors"

var (
	// ErrUnauthorized ошибка авторизации от keenetic.
	ErrUnauthorized = errors.New("unauthorized")
)
