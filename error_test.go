package tempmail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError_Error(t *testing.T) {
	httpErr := &HTTPError{
		Response: &http.Response{StatusCode: http.StatusBadRequest},
		ErrorDetails: HTTPErrorError{
			Type:   "request_error",
			Code:   "not_found",
			Detail: "The requested resource was not found",
		},
		Meta: HTTPErrorMeta{
			RequestID: "req_123456789",
		},
	}

	expected := "status 400, error type: request_error, code: not_found, detail: The requested resource was not found"
	assert.Equal(t, expected, httpErr.Error())
}

func TestHTTPError_Format(t *testing.T) {
	httpErr := &HTTPError{
		Response: &http.Response{StatusCode: http.StatusBadRequest},
		ErrorDetails: HTTPErrorError{
			Type:   "request_error",
			Code:   "not_found",
			Detail: "The requested resource was not found",
		},
		Meta: HTTPErrorMeta{
			RequestID: "req_123456789",
		},
	}

	t.Run("format with %s", func(t *testing.T) {
		result := fmt.Sprintf("%s", httpErr)
		expected := "status 400, error type: request_error, code: not_found, detail: The requested resource was not found"
		assert.Equal(t, expected, result)
	})

	t.Run("format with %v", func(t *testing.T) {
		result := fmt.Sprintf("%v", httpErr)
		expected := "status 400, error type: request_error, code: not_found, detail: The requested resource was not found"
		assert.Equal(t, expected, result)
	})

	t.Run("format with %+v", func(t *testing.T) {
		result := fmt.Sprintf("%+v", httpErr)
		expected := "status 400, error type: request_error, code: not_found, detail: The requested resource was not found, request_id: req_123456789"
		assert.Equal(t, expected, result)
	})

	t.Run("format with %q", func(t *testing.T) {
		result := fmt.Sprintf("%q", httpErr)
		expected := "\"status 400, error type: request_error, code: not_found, detail: The requested resource was not found\""
		assert.Equal(t, expected, result)
	})
}

func TestHTTPError_fullError(t *testing.T) {
	httpErr := &HTTPError{
		Response: &http.Response{StatusCode: http.StatusBadGateway},
		ErrorDetails: HTTPErrorError{
			Type:   "api_error",
			Code:   "internal_error",
			Detail: "Internal server error",
		},
		Meta: HTTPErrorMeta{
			RequestID: "req_987654321",
		},
	}

	expected := "status 502, error type: api_error, code: internal_error, detail: Internal server error, request_id: req_987654321"
	assert.Equal(t, expected, httpErr.fullError())
}
