package manifest

func NewRetryableError(msg string) *RetryableError {
	return &RetryableError{
		message: msg,
	}
}

type RetryableError struct {
	message string
}

func (e *RetryableError) Error() string {
	return e.message
}
