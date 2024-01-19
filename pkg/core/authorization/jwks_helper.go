package authorization

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/simpleg-eu/cuplan_core/pkg/core"
	"io"
	"net/http"
)

func GetJwks(url string) core.Result[*jwk.Set, core.Error] {
	response, err := http.Get(url)

	if err != nil {
		return core.Err[*jwk.Set, core.Error](*core.NewError(core.IOFailure, fmt.Sprintf("failed to get JWKS from url '%s': %v", url, err)))
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return core.Err[*jwk.Set, core.Error](*core.NewError(core.IOFailure, fmt.Sprintf("failed to read body of response: %v", err)))
	}

	jwks, err := jwk.Parse(body)

	if err != nil {
		return core.Err[*jwk.Set, core.Error](*core.NewError(core.SerializationFailure, fmt.Sprintf("failed to parse JWKS response: %v", err)))
	}

	return core.Ok[*jwk.Set, core.Error](&jwks)
}
