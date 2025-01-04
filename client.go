package temp_mail_go

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a client for the Temp Mail API.
type Client struct {
	// doer is an HTTP client.
	doer doer
	// apiKey is an API key for the Temp Mail API.
	apiKey string
}

const (
	baseURL = "https://api.temp-mail.io"

	headerAPIKey        = "X-API-Key"
	headerRateLimit     = "X-Ratelimit-Limit"
	headerRateRemaining = "X-Ratelimit-Remaining"
	headerRateUsed      = "X-Ratelimit-Used"
	headerRateReset     = "X-Ratelimit-Reset"

	userAgent = "temp-mail-go/v1.0.0"
)

// NewClient creates ready to use Client.
func NewClient(apiKey string, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		doer:   client,
		apiKey: apiKey,
	}
}

// newRequest creates a new HTTP request.
func (c *Client) newRequest(ctx context.Context, method, path string, data interface{}) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerAPIKey, c.apiKey)
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

// do sends an HTTP request and decodes the response.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	r, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		httpErr := HTTPError{
			Response: r,
		}
		if err := json.NewDecoder(r.Body).Decode(&httpErr); err != nil {
			return nil, err
		}
		return nil, &httpErr
	}

	if v != nil {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return newResponse(r), nil
}
