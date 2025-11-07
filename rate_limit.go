package tempmail

import (
	"context"
	"net/http"
	"time"
)

// Rate represents the rate limit for the current client.
type Rate struct {
	// Limit is the number of requests per hour the client is currently limited to.
	Limit int
	// Used is the number of requests the client has made within the current rate limit window.
	Used int
	// Remaining is the number of requests remaining in the current rate limit window.
	Remaining int
	// Reset is the time at which the current rate limit will reset.
	Reset time.Time
}

type rateLimitResponse struct {
	// Limit is the number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`
	// Used is the number of requests the client has made within the current rate limit window.
	Used int `json:"used"`
	// Remaining is the number of requests remaining in the current rate limit window.
	Remaining int `json:"remaining"`
	// Reset is the time at which the current rate limit window resets, in UTC epoch seconds.
	Reset int64 `json:"reset"`
}

// RateLimit returns the current rate limit for the client.
func (c *Client) RateLimit(ctx context.Context) (Rate, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/rate_limit", nil)
	if err != nil {
		return Rate{}, nil, err
	}

	var resp rateLimitResponse
	r, err := c.do(req, &resp)
	if err != nil {
		return Rate{}, nil, err
	}

	result := Rate{
		Limit:     resp.Limit,
		Used:      resp.Used,
		Remaining: resp.Remaining,
		Reset:     time.Unix(resp.Reset, 0),
	}
	// Set the Rate field in the Response since API doesn't return rate limit headers for this endpoint.
	r.Rate = result

	return result, r, nil
}
