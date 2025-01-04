package tempmail

import (
	"context"
	"net/http"
	"time"
)

const (
	// DomainTypePublic is a public domain.
	DomainTypePublic = "public"
	// DomainTypeCustom is user-provided domain.
	DomainTypeCustom = "custom"
	// DomainTypePremium is a premium domain.
	DomainTypePremium = "premium"
)

// CreateEmailOptions represents the options to create an email.
type CreateEmailOptions struct {
	// Email is the email address to create. If not provided, a random email address will be generated
	Email string
	// DomainType is the type of domain to use for the email address.
	// Possible values are: "public", "custom", "premium"
	DomainType string
	// Domain is the domain to use for the email address.
	Domain string
}

// createEmailRequest represents the request to create an email
type createEmailRequest struct {
	// Email is the email address to create. If not provided, a random email address will be generated
	Email string `json:"email"`
	// DomainType is the type of domain to use for the email address.
	// Possible values are: "public", "custom", "premium"
	DomainType string `json:"domain_type"`
	// Domain is the domain to use for the email address.
	Domain string `json:"domain"`
}

// createEmailResponse represents the response to create an email.
type createEmailResponse struct {
	// Email is the email address that was created.
	Email string `json:"email"`
	// TTL is the time to live of the email address in seconds.
	TTL int `json:"ttl"`
}

type CreateEmailResponse struct {
	// Email is the email address that was created.
	Email string
	// TTL is the time to live of the email address in seconds.
	TTL time.Duration
}

// CreateEmail creates an email address.
// It returns the email address and the time to live of the email address.
// You should use this method before getting messages for the email address.
func (c *Client) CreateEmail(ctx context.Context, options CreateEmailOptions) (CreateEmailResponse, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/emails", createEmailRequest{
		Email:      options.Email,
		DomainType: options.DomainType,
		Domain:     options.Domain,
	})
	if err != nil {
		return CreateEmailResponse{}, nil, err
	}

	var resp createEmailResponse
	r, err := c.do(req, &resp)
	if err != nil {
		return CreateEmailResponse{}, nil, err
	}

	return CreateEmailResponse{
		Email: resp.Email,
		TTL:   time.Duration(resp.TTL) * time.Second,
	}, r, nil
}
