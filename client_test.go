package tempmail

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
}
