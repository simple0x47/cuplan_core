package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/simpleg-eu/cuplan_core/pkg/core"
	"github.com/simpleg-eu/cuplan_core/pkg/core/authorization"
	"github.com/simpleg-eu/cuplan_core/pkg/core/config"
	"github.com/simpleg-eu/cuplan_core/pkg/core/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
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
	token := a.getExpiredToken()
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
	token := a.getExpiredToken()
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
	token := a.getValidToken()
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	a.router.ServeHTTP(rec, req)
	response := string(rec.Body.Bytes())

	assert.Equal(a.T(), http.StatusOK, rec.Code)
	assert.Equal(a.T(), protectedMessage, response)
}

func (a *AuthorizationTestSuite) TestExtractToken_NoToken_Error() {
	req := httptest.NewRequest("GET", protectedApi, nil)

	result := extractToken(req)

	assert.True(a.T(), result.IsErr())
	assert.Equal(a.T(), core.NotFound, result.UnwrapErr().ErrorKind)
	assert.Contains(a.T(), result.UnwrapErr().Message, "could not find token within request")
}

func (a *AuthorizationTestSuite) TestExtractToken_MalformedToken_Error() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", "Bearer malformedToken")

	result := extractToken(req)

	assert.True(a.T(), result.IsErr())
	assert.Equal(a.T(), core.InvalidToken, result.UnwrapErr().ErrorKind)
	assert.Contains(a.T(), result.UnwrapErr().Message, "failed to parse token: ")
}

func (a *AuthorizationTestSuite) TestExtractToken_ValidToken_Ok() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	token := a.getValidToken()
	req.Header.Set("Authorization", token)

	result := extractToken(req)

	assert.True(a.T(), result.IsOk())
}

func (a *AuthorizationTestSuite) TestHasTokenPermissionTo_InvalidToken_Error() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", "Bearer invalidToken")

	result := hasTokenPermissionTo(req, "some:thing")

	assert.True(a.T(), result.IsErr())
	assert.Equal(a.T(), core.InvalidToken, result.UnwrapErr().ErrorKind)
}

func (a *AuthorizationTestSuite) TestHasTokenPermissionTo_NoPermissionsToken_False() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", a.getValidToken())

	result := hasTokenPermissionTo(req, "some:thing")

	assert.True(a.T(), result.IsOk())
	assert.False(a.T(), result.Unwrap())
}

func (a *AuthorizationTestSuite) TestHasTokenPermissionTo_StringPermission_Error() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.configProvider.Get("config.yaml", "StringPermissionToken").Unwrap()))

	result := hasTokenPermissionTo(req, "some:thing")

	assert.True(a.T(), result.IsErr())
	assert.Equal(a.T(), core.InvalidToken, result.UnwrapErr().ErrorKind)
}

func (a *AuthorizationTestSuite) TestHasTokenPermissionTo_ExistingPermission_True() {
	req := httptest.NewRequest("GET", protectedApi, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.configProvider.Get("config.yaml", "ArrayPermissionsToken").Unwrap()))

	result := hasTokenPermissionTo(req, "some:thing")

	assert.True(a.T(), result.IsOk())
	assert.True(a.T(), result.Unwrap())
}

func (a *AuthorizationTestSuite) initializeRouter(issuer string, audience string) {
	logger, _ := zap.NewDevelopment()
	auth := NewAuthorization(logger, authorization.GetJwks(a.configProvider.Get("config.yaml", "JwksUri").Unwrap().(string)).Unwrap(), audience, issuer)

	a.router = chi.NewRouter()
	a.router.Use(auth.Authorize)
	a.router.Get(protectedApi, func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(protectedMessage))

		log.Printf("failed to write protected message: %v", err)
	})
}

func (a *AuthorizationTestSuite) getExpiredToken() string {
	return secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken()).Get(a.configProvider.Get("config.yaml", "ValidTokenSecret").Unwrap().(string)).Unwrap()
}

func (a *AuthorizationTestSuite) getValidToken() string {
	return secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken()).Get(a.configProvider.Get("config.yaml", "ValidTokenSecret").Unwrap().(string)).Unwrap()
}
