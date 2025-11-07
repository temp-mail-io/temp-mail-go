package tempmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_ListDomains(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/list_domains.json")), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.ListDomains(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		assert.Equal(t, ListDomainsResponse{
			Domains: []ListDomainsDomainResponse{
				{
					Name: "example.com",
					Type: "public",
				},
				{
					Name: "premium.com",
					Type: "premium",
				},
				{
					Name: "custom.com",
					Type: "custom",
				},
			},
		}, result)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.ListDomains(nil)
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.ListDomains(context.Background())
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
