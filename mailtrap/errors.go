package mailtrap

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response

	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Errors)
}
