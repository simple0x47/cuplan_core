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
