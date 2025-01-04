package tempmail

import (
	"context"
	"fmt"
	"net/http"
)

// GetMessageSourceCodeResponse is a response to get the source code of a message.
type GetMessageSourceCodeResponse struct {
	// Data is the source code of the message.
	Data string `json:"data"`
}

func (c *Client) GetMessageSourceCode(ctx context.Context, messageID string) (GetMessageSourceCodeResponse, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/messages/%s/source", messageID), nil)
	if err != nil {
		return GetMessageSourceCodeResponse{}, nil, err
	}

	var resp GetMessageSourceCodeResponse
	r, err := c.do(req, &resp)
	if err != nil {
		return GetMessageSourceCodeResponse{}, nil, err
	}

	return resp, r, nil
}
