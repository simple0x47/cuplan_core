package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/simpleg-eu/cuplan-core/pkg/core/authorization"
	"github.com/simpleg-eu/cuplan-core/pkg/core/config"
	"github.com/simpleg-eu/cuplan-core/pkg/core/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

const protectedApi = "/protected"
const protectedMessage = "This is a protected resource."

type AuthorizationTestSuite struct {
	suite.Suite
	router         chi.Router
	configProvider *config.FileProvider
}

func TestAuthorizationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorizationTestSuite))
}

func (a *AuthorizationTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	a.configProvider = config.GetConfigForTest(testFile).Unwrap()
}

func (a *AuthorizationTestSuite) TestAuthorize_NoAuthorizationHeader_Unauthorized() {
	a.initializeRouter("", "")
	req := httptest.NewRequest("GET", protectedApi, nil)
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	var errorResponse core.Error
	unmarshalError := json.Unmarshal(rec.Body.Bytes(), &errorResponse)

	if unmarshalError != nil {
		assert.Fail(a.T(), fmt.Sprintf("Failed to read response as an error: %v", unmarshalError))
	}
	assert.Equal(a.T(), http.StatusUnauthorized, rec.Code)
	assert.Equal(a.T(), core.InvalidToken, errorResponse.ErrorKind)
}

func (a *AuthorizationTestSuite) TestAuthorize_InvalidToken_Unauthorized() {
	a.initializeRouter("", "")
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", "Bearer invalidToken")
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	var errorResponse core.Error
	unmarshalError := json.Unmarshal(rec.Body.Bytes(), &errorResponse)

	if unmarshalError != nil {
		assert.Fail(a.T(), fmt.Sprintf("Failed to read response as an error: %v", unmarshalError))
	}
	assert.Equal(a.T(), http.StatusUnauthorized, rec.Code)
	assert.Equal(a.T(), core.InvalidToken, errorResponse.ErrorKind)
}

func (a *AuthorizationTestSuite) TestAuthorize_IssuerMismatchToken_Unauthorized() {
	a.initializeRouter("lmao", a.configProvider.Get("config.yaml", "Audience").Unwrap().(string))
	req := httptest.NewRequest("GET", protectedApi, nil)
	token := secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken()).Get(a.configProvider.Get("config.yaml", "ValidTokenSecret").Unwrap().(string)).Unwrap()
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	var errorResponse core.Error
	unmarshalError := json.Unmarshal(rec.Body.Bytes(), &errorResponse)

	if unmarshalError != nil {
		assert.Fail(a.T(), fmt.Sprintf("Failed to read response as an error: %v", unmarshalError))
	}
	assert.Equal(a.T(), http.StatusUnauthorized, rec.Code)
	assert.Equal(a.T(), core.InvalidToken, errorResponse.ErrorKind)
	assert.Contains(a.T(), errorResponse.Message, "failed to validate token: \"iss\" not satisfied: values do not match")
}

func (a *AuthorizationTestSuite) TestAuthorize_AudienceMismatchToken_Unauthorized() {
	a.initializeRouter(a.configProvider.Get("config.yaml", "Issuer").Unwrap().(string), "audience?")
	req := httptest.NewRequest("GET", protectedApi, nil)
	token := secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken()).Get(a.configProvider.Get("config.yaml", "ValidTokenSecret").Unwrap().(string)).Unwrap()
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	var errorResponse core.Error
	unmarshalError := json.Unmarshal(rec.Body.Bytes(), &errorResponse)

	if unmarshalError != nil {
		assert.Fail(a.T(), fmt.Sprintf("Failed to read response as an error: %v", unmarshalError))
	}
	assert.Equal(a.T(), http.StatusUnauthorized, rec.Code)
	assert.Equal(a.T(), core.InvalidToken, errorResponse.ErrorKind)
	assert.Contains(a.T(), errorResponse.Message, "failed to validate token: aud not satisfied")
}

func (a *AuthorizationTestSuite) TestAuthorize_ValidToken_Authorized() {
	a.initializeRouter(a.configProvider.Get("config.yaml", "Issuer").Unwrap().(string), a.configProvider.Get("config.yaml", "Audience").Unwrap().(string))
	req := httptest.NewRequest("GET", protectedApi, nil)
	token := secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken()).Get(a.configProvider.Get("config.yaml", "ValidTokenSecret").Unwrap().(string)).Unwrap()
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	response := string(rec.Body.Bytes())

	assert.Equal(a.T(), http.StatusOK, rec.Code)
	assert.Equal(a.T(), protectedMessage, response)
}

func (a *AuthorizationTestSuite) initializeRouter(issuer string, audience string) {
	logger, _ := zap.NewDevelopment()
	auth := NewAuthorization(logger, authorization.GetJwks(a.configProvider.Get("config.yaml", "JwksUri").Unwrap().(string)).Unwrap(), audience, issuer)

	a.router = chi.NewRouter()
	a.router.Use(auth.Authorize)
	a.router.Get(protectedApi, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(protectedMessage))
	})
}
