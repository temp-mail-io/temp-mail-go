package tempmail

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_ListEmailMessages(t *testing.T) {
	mDoer := newMockDoer(t)
	mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(200, readFile(t, "testdata/list_email_messages.json")), nil)

	c := newClient()
	c.doer = mDoer
	r, resp, err := c.ListEmailMessages(context.Background(), "user@example.com")
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.Len(t, r, 2)
}
