package tempmail

import (
	"net/http"
	"strconv"
	"time"
)

// Response is a GitHub API response. This wraps the standard http.Response
// returned from GitHub and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// Explicitly specify the Rate type so Rate's String() receiver doesn't
	// propagate to Response.
	Rate Rate
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.Rate = parseRate(r)
	return response
}

// parseRate parses the rate related headers.
func parseRate(r *http.Response) Rate {
	var rate Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if used := r.Header.Get(headerRateUsed); used != "" {
		rate.Used, _ = strconv.Atoi(used)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = time.Unix(v, 0)
		}
	}
	return rate
}
