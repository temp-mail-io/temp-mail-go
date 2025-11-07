package tempmail

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewClient("API_KEY", nil)
		require.NotNil(t, c)
		assert.Equal(t, "API_KEY", c.apiKey)
		assert.Same(t, http.DefaultClient, c.doer)
	})

	t.Run("custom client", func(t *testing.T) {
		httpClient := &http.Client{}
		c := NewClient("API_KEY", httpClient)
		require.NotNil(t, c)
		assert.Same(t, httpClient, c.doer)
	})
}

func newClient() *Client {
	return NewClient("API_KEY", nil)
}

func TestNewRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := newClient()
		req, err := c.newRequest(context.Background(), http.MethodGet, "/v1/emails", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "https://api.temp-mail.io/v1/emails", req.URL.String())
		assert.Equal(t, "API_KEY", req.Header.Get(headerAPIKey))
		assert.Equal(t, "temp-mail-go/v1.0.0", req.UserAgent())
	})

	t.Run("custom body", func(t *testing.T) {
		c := newClient()
		req, err := c.newRequest(context.Background(), http.MethodPost, "/v1/emails", createEmailRequest{
			Domain: "example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, req)
		assert.Equal(t, http.MethodPost, req.Method)
		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"domain":"example.com"}`, string(b))
	})

	t.Run("JSON marshal error", func(t *testing.T) {
		c := newClient()
		// Create a value that cannot be marshaled to JSON
		invalidData := make(chan int)
		req, err := c.newRequest(context.Background(), http.MethodPost, "/v1/emails", invalidData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})
}

// readFile reads the file from the given path.
func readFile(t *testing.T, path string) []byte {
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return b
}

// newTestResponse creates a new test response with the given status code and body.
func newTestResponse(statusCode int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

func TestClient_do(t *testing.T) {
	t.Run("success with nil v", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, []byte{}), nil)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		resp, err := c.do(req, nil)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("error from rawDo", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(nil, assert.AnError)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		_, err = c.do(req, nil)
		assert.EqualError(t, err, assert.AnError.Error())
	})

	t.Run("error from checkResponse", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		_, err = c.do(req, nil)
		require.Error(t, err)
		var httpErr *HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.NotNil(t, httpErr.Response)
		assert.Equal(t, HTTPErrorError{
			Type:   "request_error",
			Code:   "not_found",
			Detail: "Attachment not found",
		}, httpErr.ErrorDetails)
		assert.Equal(t, HTTPErrorMeta{
			RequestID: "req_123456789",
		}, httpErr.Meta)
	})

	t.Run("JSON decode error", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(200, []byte("invalid json")), nil)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		var result map[string]interface{}
		_, err = c.do(req, &result)
		require.Error(t, err)
		assert.EqualError(t, err, "invalid character 'i' looking for beginning of value")
	})
}

func TestClient_checkResponse(t *testing.T) {
	t.Run("success response", func(t *testing.T) {
		c := newClient()
		resp := newResponse(newTestResponse(http.StatusOK, []byte{}))
		err := c.checkResponse(resp)
		require.NoError(t, err)
	})

	t.Run("error response", func(t *testing.T) {
		c := newClient()
		resp := newResponse(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")))
		err := c.checkResponse(resp)
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 400, httpErr.Response.StatusCode)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})

	t.Run("error response with invalid JSON", func(t *testing.T) {
		c := newClient()
		resp := newResponse(newTestResponse(http.StatusBadGateway, []byte("invalid json")))
		err := c.checkResponse(resp)
		require.Error(t, err)
		assert.Errorf(t, err, "invalid character 'i' looking for beginning of value")
	})
}

func TestClient_rawDo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		expectedResp := newTestResponse(http.StatusOK, []byte("test"))
		mDoer.EXPECT().Do(mock.Anything).Return(expectedResp, nil)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		resp, err := c.rawDo(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("error from doer", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(nil, assert.AnError)

		c := newClient()
		c.doer = mDoer
		req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
		require.NoError(t, err)
		_, err = c.rawDo(req)
		assert.EqualError(t, err, assert.AnError.Error())
	})
}
