package tempmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_DownloadAttachment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedContent := []byte("test attachment content")

		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, expectedContent), nil)

		c := newClient()
		c.doer = mDoer
		content, resp, err := c.DownloadAttachment(context.Background(), "01JE97K1PBYVGKY0PVE3KXSBF9")
		require.NoError(t, err)
		assert.Equal(t, expectedContent, content)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.DownloadAttachment(nil, "01JE97K1PBYVGKY0PVE3KXSBF9")
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from rawDo", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(nil, assert.AnError)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.DownloadAttachment(context.Background(), "01JE97K1PBYVGKY0PVE3KXSBF9")
		assert.EqualError(t, err, assert.AnError.Error())
	})

	t.Run("error response", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.DownloadAttachment(context.Background(), "nonexistent")
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})

	t.Run("io.ReadAll error", func(t *testing.T) {
		mReadCloser := newMockReadCloser(t)
		mReadCloser.EXPECT().Read(mock.Anything).Return(0, assert.AnError)
		mReadCloser.EXPECT().Close().Return(nil)

		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       mReadCloser,
		}, nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.DownloadAttachment(context.Background(), "01JE97K1PBYVGKY0PVE3KXSBF9")
		assert.EqualError(t, err, assert.AnError.Error())
	})
}
