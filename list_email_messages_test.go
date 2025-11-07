package tempmail

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_ListEmailMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/list_email_messages.json")), nil)

		c := newClient()
		c.doer = mDoer
		r, resp, err := c.ListEmailMessages(context.Background(), "user@example.com")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		require.Len(t, r.Messages, 1)
		assert.Equal(t, ListEmailMessagesMessageResponse{
			ID:        "01JE97FT950QRPDYGDXJ4R43QR",
			From:      "admin@example.com",
			To:        "user@example.com",
			CC:        []string{"another_user@example.com"},
			Subject:   "Your account has been created",
			BodyText:  "Welcome to our service! Your account has been created successfully.",
			BodyHTML:  "<p>Welcome to our service! Your account has been created successfully.</p>",
			CreatedAt: time.Date(2022, 1, 31, 22, 0, 0, 0, time.UTC),
			Attachments: []ListEmailMessagesAttachmentResponse{
				{
					ID:   "01JE97K1PBYVGKY0PVE3KXSBF9",
					Name: "invoice.pdf",
					Size: 5120,
				},
			},
		}, r.Messages[0])
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.ListEmailMessages(nil, "user@example.com")
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.ListEmailMessages(context.Background(), "nonexistent@example.com")
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
