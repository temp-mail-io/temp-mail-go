package temp_mail_go

import (
	"fmt"
	"io"
)

// HTTPError is the error response that will be returned by the API.
type HTTPError struct {
	ErrorDetails HTTPErrorError `json:"error"`
	Meta         HTTPErrorMeta  `json:"meta"`
}

type HTTPErrorError struct {
	// Type is the type of the error.
	// Possible values: api_error, request_error.
	Type string `json:"type"`
	// Code is the error code.
	Code string `json:"code"`
	// Detail is the error message.
	Detail string `json:"detail"`
}

type HTTPErrorMeta struct {
	RequestID string `json:"request_id"`
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("error: %s, code: %s, detail: %s", h.ErrorDetails.Type, h.ErrorDetails.Code, h.ErrorDetails.Detail)
}

func (h *HTTPError) fullError() string {
	return fmt.Sprintf("error: %s, code: %s, detail: %s, request_id: %s", h.ErrorDetails.Type, h.ErrorDetails.Code, h.ErrorDetails.Detail, h.Meta.RequestID)
}

// Format implements fmt.Formatter interface.
// It adds request ID to the error message when called with %+v verb.
func (h *HTTPError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, h.fullError())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, h.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", h.Error())
	}
}
