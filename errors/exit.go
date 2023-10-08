package errors

type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	return e.Message
}

func (e *ExitError) ExitCode() int {
	return e.Code
}
