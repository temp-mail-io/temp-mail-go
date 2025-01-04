package tempmail

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteMessage deletes a message by its ID.
func (c *Client) DeleteMessage(ctx context.Context, messageID string) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/messages/%s", messageID), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}
