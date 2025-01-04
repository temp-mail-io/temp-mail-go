package temp_mail_go

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// DownloadAttachment downloads an attachment by its ID and returns the raw bytes.
func (c *Client) DownloadAttachment(ctx context.Context, attachmentID string) ([]byte, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/attachments/%s", attachmentID), nil)
	if err != nil {
		return nil, nil, err
	}

	r, err := c.rawDo(req)
	if err != nil {
		return nil, nil, err
	}
	defer r.Body.Close()

	if err := c.checkResponse(r); err != nil {
		return nil, nil, err
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}

	return b, r, nil
}
