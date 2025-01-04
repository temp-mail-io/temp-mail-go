package temp_mail_go

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ListEmailMessagesResponse represents the response to list email messages.
type ListEmailMessagesResponse struct {
	Messages []ListEmailMessagesMessageResponse `json:"messages"`
}

// ListEmailMessagesMessageResponse represents an email message.
type ListEmailMessagesMessageResponse struct {
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
	Attachments []ListEmailMessagesAttachmentResponse `json:"attachments"`
}

// ListEmailMessagesAttachmentResponse represents an attachment of an email message.
type ListEmailMessagesAttachmentResponse struct {
	// ID is the unique identifier of the attachment.
	ID string `json:"id"`
	// Name is the name of the attachment.
	// For example, "image.png".
	Name string `json:"name"`
	// Size is the size of the attachment in bytes.
	Size int `json:"size"`
}

// ListEmailMessages returns all messages for the email address.
func (c *Client) ListEmailMessages(ctx context.Context, email string) ([]ListEmailMessagesResponse, *Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/emails/%s/messages", email), nil)
	if err != nil {
		return nil, nil, err
	}

	var messages []ListEmailMessagesResponse
	resp, err := c.do(req, &messages)
	if err != nil {
		return nil, nil, err
	}

	return messages, resp, nil
}
