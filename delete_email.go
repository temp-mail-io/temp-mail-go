package temp_mail_go

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteEmail deletes an email address.
func (c *Client) DeleteEmail(ctx context.Context, email string) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/emails/%s", email), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}
