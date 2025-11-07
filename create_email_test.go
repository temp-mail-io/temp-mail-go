package tempmail

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	client *Client
	suite.Suite
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	apiKey := os.Getenv("TEST_API_KEY")
	if apiKey == "" {
		s.T().Skip("TEST_API_KEY is not set")
	}
	s.client = NewClient(apiKey, nil)
}

func (s *IntegrationSuite) TestCreateEmail() {
	email, resp, err := s.client.CreateEmail(context.Background(), CreateEmailOptions{})
	s.Require().NoError(err)
	s.Assert().NotEmpty(email.Email)
	s.Assert().NotZero(email.TTL)
	s.Require().NotNil(resp)
	testRateLimit(s.T(), resp)
}

func (s *IntegrationSuite) TestDeleteEmail() {
	email, _, err := s.client.CreateEmail(context.Background(), CreateEmailOptions{})
	s.Require().NoError(err)
	resp, err := s.client.DeleteEmail(context.Background(), email.Email)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	testRateLimit(s.T(), resp)
}

func (s *IntegrationSuite) TestListDomains() {
	domains, resp, err := s.client.ListDomains(context.Background())
	s.Require().NoError(err)
	s.Require().NotEmpty(domains.Domains)
	s.Require().NotNil(resp)
	testRateLimit(s.T(), resp)
}

func testRateLimit(t *testing.T, resp *Response) {
	t.Helper()
	assert.Positive(t, resp.Rate.Limit)
	assert.NotZero(t, resp.Rate.Remaining)
	assert.NotZero(t, resp.Rate.Used)
	assert.True(t, resp.Rate.Reset.After(time.Now()), "reset time should be in the future")
}

func TestClient_CreateEmail(t *testing.T) {
	t.Run("success with empty options", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/create_email.json")), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.CreateEmail(context.Background(), CreateEmailOptions{})
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		expected := CreateEmailResponse{
			Email: "test@example.com",
			TTL:   3600 * time.Second,
		}
		assert.Equal(t, expected, result)
	})

	t.Run("success with custom options", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusOK, readFile(t, "testdata/create_email.json")), nil)

		c := newClient()
		c.doer = mDoer
		result, resp, err := c.CreateEmail(context.Background(), CreateEmailOptions{
			Email:      "custom@example.com",
			DomainType: DomainTypePublic,
			Domain:     "example.com",
		})
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, result.Email)
		assert.NotZero(t, result.TTL)
	})

	t.Run("error from newRequest", func(t *testing.T) {
		c := newClient()
		_, _, err := c.CreateEmail(nil, CreateEmailOptions{})
		assert.EqualError(t, err, "net/http: nil Context")
	})

	t.Run("error from do", func(t *testing.T) {
		mDoer := newMockDoer(t)
		mDoer.EXPECT().Do(mock.Anything).Return(newTestResponse(http.StatusBadRequest, readFile(t, "testdata/error_response.json")), nil)

		c := newClient()
		c.doer = mDoer
		_, _, err := c.CreateEmail(context.Background(), CreateEmailOptions{Email: "invalid"})
		require.Error(t, err)
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, "request_error", httpErr.ErrorDetails.Type)
		assert.Equal(t, "not_found", httpErr.ErrorDetails.Code)
	})
}
