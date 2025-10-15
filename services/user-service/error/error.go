package error

type AppError struct {
	Code    int    `json:"code"`    // Custom error code
	Message string `json:"message"` // Error message
	Err     error  `json:"-"`       // Internal error (not exposed)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAppError creates a new custom error
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithErr creates a new custom error with underlying error
func NewAppErrorWithErr(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
