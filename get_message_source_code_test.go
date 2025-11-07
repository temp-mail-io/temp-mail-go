package tempmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMessageSourceCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := `{
			"data": "Return-Path: <sender@example.com>\nReceived: from mail.example.com\nSubject: Test Message\n\nThis is the message body."
		}`

		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, []byte(responseData)), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.GetMessageSourceCode(context.Background(), "01JE97FT950QRPDYGDXJ4R43QR")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, GetMessageSourceCodeResponse{
			Data: "Return-Path: <sender@example.com>\nReceived: from mail.example.com\nSubject: Test Message\n\nThis is the message body.",
		}, result)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.GetMessageSourceCode(nil, "01JE97FT950QRPDYGDXJ4R43QR")
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.GetMessageSourceCode(context.Background(), "nonexistent_id")
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
