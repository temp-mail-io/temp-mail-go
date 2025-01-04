package temp_mail_go

import (
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

const baseURL = "https://api.temp-mail.io"

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
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	return req, nil
}

// do sends an HTTP request and decodes the response.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var httpErr HTTPError
		if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
			return err
		}
		return &httpErr
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
