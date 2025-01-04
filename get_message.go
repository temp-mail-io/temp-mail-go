package temp_mail_go

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// GetMessageResponse is a response to get a message.
type GetMessageResponse struct {
	// ID is the unique identifier of the email message.
	ID string `json:"id"`
	// From is the email address of the sender.
	From string `json:"from"`
	// To is the email address of the recipient.
	To string `json:"to"`
	// CC is the email addresses of the CC recipients.
	CC []string `json:"cc"`
	// Subject is the subject of the email message.
	Subject string `json:"subject"`
	// BodyText is the plain text body of the email message.
	BodyText string `json:"body_text"`
	// BodyHTML is the HTML body of the email message.
	BodyHTML string `json:"body_html"`
	// CreatedAt is the time when the email message was created.
	CreatedAt time.Time `json:"created_at"`
	// Attachments is the list of attachments of the email message.
	Attachments []GetMessageAttachmentResponse `json:"attachments"`
}

// GetMessageAttachmentResponse represents an attachment of an email message.
type GetMessageAttachmentResponse struct {
	// ID is the unique identifier of the attachment.
	ID string `json:"id"`
	// Name is the name of the attachment.
	// For example, "image.png".
	Name string `json:"name"`
	// Size is the size of the attachment in bytes.
	Size int `json:"size"`
}

// GetMessage gets a message by its ID.
func (c *Client) GetMessage(ctx context.Context, messageID string) (GetMessageResponse, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/messages/%s", messageID), nil)
	if err != nil {
		return GetMessageResponse{}, nil, err
	}

	var resp GetMessageResponse
	r, err := c.do(req, &resp)
	if err != nil {
		return GetMessageResponse{}, nil, err
	}

	return resp, r, nil
}
