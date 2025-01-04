package tempmail

import (
	"context"
	"net/http"
)

// ListDomainsResponse represents the response for the listDomainsHandler.
type ListDomainsResponse struct {
	Domains []ListDomainsDomainResponse `json:"domains"`
}

// ListDomainsDomainResponse represents one domain in the listDomainsHandler response.
type ListDomainsDomainResponse struct {
	// Name of the domain.
	Name string `json:"name"`
	// Type of the domain.
	// Possible values: "public", "premium", "custom"
	Type string `json:"type"`
}

// ListDomains returns a list of domains available for use.
func (c *Client) ListDomains(ctx context.Context) (ListDomainsResponse, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/domains", nil)
	if err != nil {
		return ListDomainsResponse{}, nil, err
	}

	var resp ListDomainsResponse
	r, err := c.do(req, &resp)
	if err != nil {
		return ListDomainsResponse{}, nil, err
	}

	return resp, r, nil
}
