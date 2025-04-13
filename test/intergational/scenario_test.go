//go:build integrations

package intergational

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	ctx            context.Context
	cancel         context.CancelFunc
	client         *resty.Client
	moderatorToken string
	employeeToken  string
}

func (s *IntegrationSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)

	s.client = resty.New().SetBaseURL("http://localhost:8080")
}

func (s *IntegrationSuite) TearDownTest() {
	s.cancel()
}

func (s *IntegrationSuite) dummyLoginHelper(role string) string {
	body := map[string]interface{}{
		"role": role,
	}
	var tokenResp struct {
		Token string `json:"Token"`
	}
	r, err := s.client.R().
		SetBody(body).
		Post("/dummyLogin")
	s.T().Logf("DummyLogin (%s) response: %s", role, r.Body())
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, r.StatusCode())

	err = json.Unmarshal(r.Body(), &tokenResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), tokenResp.Token)

	return tokenResp.Token
}

// TestFullPVZPipeline realizes scenario:
//  1. moderator gets token creates new PVZ
//  2. employee gets token and craetes new reception
//  3. employee adds 50 products to reception
//  4. employee closes receptions
func (s *IntegrationSuite) TestPVZPipeline() {
	// 1. Get Moderator token
	s.moderatorToken = s.dummyLoginHelper("moderator")

	// 2. Create PVZ with POST /pvz
	pvzID := uuid.New()
	registrationDate := time.Now().Format(time.RFC3339)
	pvzBody := map[string]interface{}{
		"id":                pvzID,
		"registration_date": registrationDate,
		"city":              "Москва",
	}

	var pvzResp struct {
		ID string `json:"id"`
	}
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.moderatorToken)).
		SetBody(pvzBody).
		Post("/pvz")
	s.T().Logf("Create PVZ response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), r.StatusCode(), http.StatusCreated)

	err = json.Unmarshal(r.Body(), &pvzResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), pvzResp.ID)

	// 3. Get Employee token
	s.employeeToken = s.dummyLoginHelper("employee")

	// 4. Create Reception with POST /receptions
	receptionBody := map[string]interface{}{
		"pvz_id": pvzID,
	}

	var receptionResp struct {
		ID string `json:"id"`
	}
	r, err = s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.employeeToken)).
		SetBody(receptionBody).
		Post("/receptions")
	s.T().Logf("Craete Reception response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), r.StatusCode(), http.StatusCreated)

	err = json.Unmarshal(r.Body(), &receptionResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), receptionResp.ID)

	// 5. Add 50 Products to Receptions with POST /reception
	for i := 1; i <= 50; i++ {
		productBody := map[string]interface{}{
			"type":   "одежда",
			"pvz_id": pvzID.String(),
		}

		r, err = s.client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.employeeToken)).
			SetBody(productBody).
			Post("/products")
		s.T().Logf("Add Product %d response: %s", i, r.Body())
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.StatusCreated, r.StatusCode())
	}

	// 6. Close Reception with /POST /pvz/{pvzId}/close_last_reception
	closeReceptionURL := fmt.Sprintf("/pvz/%s/close_last_reception", pvzID.String())
	r, err = s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.employeeToken)).
		Post(closeReceptionURL)
	s.T().Logf("Close Reception response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, r.StatusCode())
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
