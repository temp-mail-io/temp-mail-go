package tempmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_DeleteMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, []byte("{}")), nil)

		c := newClient()
		c.doer = mDoer
		resp, err := c.DeleteMessage(context.Background(), "01JE97FT950QRPDYGDXJ4R43QR")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, err := c.DeleteMessage(nil, "01JE97FT950QRPDYGDXJ4R43QR")
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, err := c.DeleteMessage(context.Background(), "nonexistent_id")
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
