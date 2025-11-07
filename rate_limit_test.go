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

func TestClient_RateLimit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/rate_limit.json")), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.RateLimit(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		expected := Rate{
			Limit:     1000,
			Used:      100,
			Remaining: 900,
			Reset:     time.Unix(1640995200, 0),
		}
		assert.Equal(t, expected, result)
		assert.Equal(t, expected, resp.Rate)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.RateLimit(nil)
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.RateLimit(context.Background())
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
