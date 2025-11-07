package tempmail

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponse(t *testing.T) {
	httpResp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			headerRateLimit:     []string{"1000"},
			headerRateRemaining: []string{"900"},
			headerRateUsed:      []string{"100"},
			headerRateReset:     []string{"1640995200"},
		},
	}

	resp := newResponse(httpResp)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1000, resp.Rate.Limit)
	assert.Equal(t, 900, resp.Rate.Remaining)
	assert.Equal(t, 100, resp.Rate.Used)
	assert.Equal(t, time.Unix(1640995200, 0), resp.Rate.Reset)
}

func TestParseRate(t *testing.T) {
	t.Run("all headers present", func(t *testing.T) {
		httpResp := &http.Response{
			Header: http.Header{
				headerRateLimit:     []string{"1000"},
				headerRateRemaining: []string{"900"},
				headerRateUsed:      []string{"100"},
				headerRateReset:     []string{"1640995200"},
			},
		}

		rate := parseRate(httpResp)
		assert.Equal(t, 1000, rate.Limit)
		assert.Equal(t, 900, rate.Remaining)
		assert.Equal(t, 100, rate.Used)
		assert.Equal(t, time.Date(2022, 1, 1, 4, 0, 0, 0, time.Local), rate.Reset)
	})

	t.Run("no headers present", func(t *testing.T) {
		httpResp := &http.Response{
			Header: http.Header{},
		}

		rate := parseRate(httpResp)
		assert.Equal(t, 0, rate.Limit)
		assert.Equal(t, 0, rate.Remaining)
		assert.Equal(t, 0, rate.Used)
		assert.Equal(t, time.Time{}, rate.Reset)
	})

	t.Run("invalid header values", func(t *testing.T) {
		httpResp := &http.Response{
			Header: http.Header{
				headerRateLimit:     []string{"invalid"},
				headerRateRemaining: []string{"also_invalid"},
				headerRateUsed:      []string{"not_a_number"},
				headerRateReset:     []string{"not_a_timestamp"},
			},
		}

		rate := parseRate(httpResp)
		assert.Equal(t, 0, rate.Limit)
		assert.Equal(t, 0, rate.Remaining)
		assert.Equal(t, 0, rate.Used)
		assert.Equal(t, time.Time{}, rate.Reset)
	})

	t.Run("zero reset timestamp", func(t *testing.T) {
		httpResp := &http.Response{
			Header: http.Header{
				headerRateReset: []string{"0"},
			},
		}

		rate := parseRate(httpResp)
		assert.Equal(t, time.Time{}, rate.Reset)
	})

	t.Run("partial headers present", func(t *testing.T) {
		httpResp := &http.Response{
			Header: http.Header{
				headerRateLimit: []string{"500"},
				headerRateUsed:  []string{"50"},
			},
		}

		rate := parseRate(httpResp)
		assert.Equal(t, 500, rate.Limit)
		assert.Equal(t, 0, rate.Remaining) // not present
		assert.Equal(t, 50, rate.Used)
		assert.Equal(t, time.Time{}, rate.Reset) // not present
	})
}
