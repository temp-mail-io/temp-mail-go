package temp_mail_go

import "net/http"

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
