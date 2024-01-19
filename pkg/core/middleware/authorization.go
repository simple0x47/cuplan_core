package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/simpleg-eu/cuplan_core/pkg/core"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Authorization struct {
	logger   *zap.Logger
	jwks     *jwk.Set
	audience string
	issuer   string
}

func NewAuthorization(logger *zap.Logger, jwks *jwk.Set, audience string, issuer string) *Authorization {
	a := new(Authorization)
	a.logger = logger
	a.jwks = jwks
	a.audience = audience
	a.issuer = issuer

	return a
}

func HasTokenPermissionTo(r *http.Request, permission string) core.Result[bool, core.Error] {
	result := extractToken(r)

	if result.IsErr() {
		return core.Err[bool, core.Error](result.UnwrapErr())
	}

	token := result.Unwrap()

	objects, exists := token.Get("permissions")

	if !exists {
		return core.Ok[bool, core.Error](false)
	}

	objs, ok := objects.([]any)

	if !ok {
		return core.Err[bool, core.Error](*core.NewError(core.InvalidToken, "failed to read 'permissions' claim as a slice"))
	}

	for _, object := range objs {
		tokenPermission, ok := object.(string)

		if !ok {
			return core.Err[bool, core.Error](*core.NewError(core.InvalidToken, fmt.Sprintf("expected 'permissions' to be a slice of strings but it contains: %v", object)))
		}

		if tokenPermission == permission {
			return core.Ok[bool, core.Error](true)
		}
	}

	return core.Ok[bool, core.Error](false)
}

func extractToken(r *http.Request) core.Result[jwt.Token, core.Error] {
	str := extractBearerToken(r)

	if str.IsNone() {
		return core.Err[jwt.Token, core.Error](*core.NewError(core.NotFound, "could not find token within request"))
	}

	token, err := jwt.ParseString(str.Unwrap())

	if err != nil {
		return core.Err[jwt.Token, core.Error](*core.NewError(core.InvalidToken, fmt.Sprintf("failed to parse token: %v", err)))
	}

	return core.Ok[jwt.Token, core.Error](token)
}

func extractBearerToken(r *http.Request) core.Option[string] {
	header := r.Header.Get("Authorization")

	if len(header) == 0 {
		return core.None[string]()
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return core.None[string]()
	}

	return core.Some[string](strings.TrimPrefix(header, "Bearer "))
}

func (a *Authorization) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		optionalToken := extractBearerToken(r)

		if optionalToken.IsNone() {
			a.writeInvalidTokenError("bearer token is missing", w)
			return
		}

		stringToken := optionalToken.Unwrap()

		token, parseError := jwt.ParseString(stringToken, jwt.WithKeySet(*a.jwks))

		if parseError != nil {
			a.writeInvalidTokenError(fmt.Sprintf("invalid token: %v", parseError), w)
			return
		}

		validationResult := a.validateToken(&token)

		if validationResult.IsErr() {
			a.writeInvalidTokenError(validationResult.UnwrapErr().Message, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authorization) writeInvalidTokenError(message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	err := core.NewError(core.InvalidToken, message)
	data, marshalError := json.Marshal(*err)

	if marshalError != nil {
		a.logger.Warn(fmt.Sprintf("failed to json marshal error: %v", marshalError))
	}

	_, writeError := w.Write(data)

	if writeError != nil {
		a.logger.Warn(fmt.Sprintf("failed to write error bytes as response: %v", writeError))
	}
}

func (a *Authorization) validateToken(token *jwt.Token) core.Result[core.Empty, core.Error] {
	err := jwt.Validate(*token, jwt.WithIssuer(a.issuer), jwt.WithAudience(a.audience))

	if err != nil {
		return core.Err[core.Empty, core.Error](*core.NewError(core.InvalidToken, fmt.Sprintf("failed to validate token: %v", err)))
	}

	return core.Ok[core.Empty, core.Error](core.Empty{})
}
