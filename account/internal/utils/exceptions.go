package utils

type (
	// HTTPError custom http error to handle custom errors
	HTTPError struct {
		error      `json:"-"`
		StatusCode int    `json:"-"`
		Message    string `json:"message"`
	}
	// UnauthorizedError to handle unauthorized errors
	UnauthorizedError struct {
		HTTPError
	}
	// ForbiddenError to handle forbidden errors
	ForbiddenError struct {
		HTTPError
	}
)

// GetNotFoundError returns error associated with HTTP NotFound error
func GetNotFoundError(message string) error {
	if message == "" {
		message = "آیتمی با اطلاعات ارسالی یافت نشد"
	}
	err := HTTPError{}
	err.StatusCode = 404
	err.Message = message
	return err
}

// GetUnAuthorizedError returns error associated with HTTP Unauthorized error
func GetUnAuthorizedError(message string) error {
	if message == "" {
		message = "شما به این قسمت دسترسی ندارید"
	}
	err := HTTPError{}
	err.StatusCode = 401
	err.Message = message
	return err
}

// GetForbiddenError returns error associated with HTTP Forbidden error
func GetForbiddenError(message string) error {
	if message == "" {
		message = "شما به این قسمت دسترسی ندارید"
	}
	err := HTTPError{}
	err.StatusCode = 403
	err.Message = message
	return err
}

// GetValidationError returns error associated with HTTP Vlidation error
func GetValidationError(message string) error {
	if message == "" {
		message = "درخواست معتبر نمی‌باشد"
	}
	err := HTTPError{}
	err.StatusCode = 400
	err.Message = message
	return err
}

// ReturnSuccess returns json associated with HTTP
func ReturnSuccess(message string) error {
	if message == "" {
		message = ""
	}
	err := HTTPError{}
	err.StatusCode = 200
	err.Message = message
	return err
}
