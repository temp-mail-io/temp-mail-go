package tempmail

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
