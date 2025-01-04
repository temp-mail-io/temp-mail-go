package temp_mail_go

import "context"

type RateLimitResponse struct {
	// Limit is the maximum number of requests that you can make per period.
	Limit int64 `json:"limit"`
	// Used is the number of requests you have made in the current rate limit window.
	Used int64 `json:"used"`
	// Remaining is the number of requests remaining in the current rate limit window.
	Remaining int64 `json:"remaining"`
	// Reset is the time at which the current rate limit window resets, in UTC epoch seconds.
	Reset int64 `json:"reset"`
}

// RateLimit returns the current rate limit for the client.
func (c *Client) RateLimit(ctx context.Context) (RateLimitResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/v1/rate_limit", nil)
	if err != nil {
		return RateLimitResponse{}, err
	}

	var resp RateLimitResponse
	if err := c.do(req, &resp); err != nil {
		return RateLimitResponse{}, err
	}

	return resp, nil
}
