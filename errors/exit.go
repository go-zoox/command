package errors

// ExitError is an error that indicates an exit.
type ExitError struct {
	Code    int
	Message string
}

// Error returns the error message.
func (e *ExitError) Error() string {
	return e.Message
}

// ExitCode returns the exit code.
func (e *ExitError) ExitCode() int {
	return e.Code
}
