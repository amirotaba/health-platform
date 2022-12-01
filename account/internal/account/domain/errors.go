package domain

type (
	// IHTTPError HTTPError custom http error to handle custom errors
	IHTTPError interface {
		GetStatusCode() int
	}
	HTTPError struct {
		IHTTPError `json:"-"`

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

func (c *HTTPError) GetStatusCode() int {
	return c.StatusCode
}
