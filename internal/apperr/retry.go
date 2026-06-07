package apperr

import "errors"

func AsRetryable(err error, target **RetryableError) bool {
	return errors.As(err, target)
}
