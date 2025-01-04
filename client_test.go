package tempmail

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewClient("API_KEY", nil)
		require.NotNil(t, c)
		assert.Equal(t, "API_KEY", c.apiKey)
		assert.Same(t, http.DefaultClient, c.doer)
	})

	t.Run("custom client", func(t *testing.T) {
		httpClient := &http.Client{}
		c := NewClient("API_KEY", httpClient)
		require.NotNil(t, c)
		assert.Same(t, httpClient, c.doer)
	})
}
