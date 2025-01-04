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

	resp, err := c.rawDo(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return b, resp, nil
}
