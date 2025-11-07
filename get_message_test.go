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

func TestClient_GetMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/get_message.json")), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.GetMessage(context.Background(), "01JE97FT950QRPDYGDXJ4R43QR")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		assert.Equal(t, GetMessageResponse{
			ID:        "01JE97FT950QRPDYGDXJ4R43QR",
			From:      "sender@example.com",
			To:        "recipient@example.com",
			CC:        []string{"cc@example.com"},
			Subject:   "Test Message",
			BodyText:  "This is a test message.",
			BodyHTML:  "<p>This is a test message.</p>",
			CreatedAt: time.Date(2022, 1, 31, 22, 0, 0, 0, time.UTC),
			Attachments: []GetMessageAttachmentResponse{
				{
					ID:   "01JE97K1PBYVGKY0PVE3KXSBF9",
					Name: "document.pdf",
					Size: 2048,
				},
			},
		}, result)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.GetMessage(nil, "01JE97FT950QRPDYGDXJ4R43QR")
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.GetMessage(context.Background(), "nonexistent_id")
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
