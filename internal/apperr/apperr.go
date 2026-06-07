package apperr

import "time"

type UserError struct{ Msg string }
type RetryableError struct {
	Msg        string
	RetryAfter time.Duration
}
type InternalError struct{ Msg string }

func (e *UserError) Error() string      { return e.Msg }
func (e *RetryableError) Error() string { return e.Msg }
func (e *InternalError) Error() string  { return e.Msg }
